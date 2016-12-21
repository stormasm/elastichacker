package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/attic-labs/noms/go/types"
	"github.com/stormasm/firego"
	"github.com/stormasm/hackerelastic1/elastic"
)

type datum struct {
	index float64
	value types.Struct
}

func main() {
	hv := bigSync()
	fmt.Println(hv.Hash().String())
	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	fmt.Println(input)
}

func bigSync() types.Value {
	newIndex := make(chan float64, 1000)
	newDatum := make(chan datum, 100)
	streamData := make(chan types.Value, 100)
	newMap := types.NewStreamingMap(types.NewTestValueStore(), streamData)

	go func() {
		for i := 8432709.0; i < 8432712.0; i++ {
			newIndex <- i
		}
		close(newIndex)
	}()

	workerPool(500, func() {
		churn(newIndex, newDatum)
	}, func() {
		close(newDatum)
	})

	for datum := range newDatum {
		streamData <- types.Number(datum.index)
		streamData <- datum.value
	}
	close(streamData)

	fmt.Println("generating map...")
	mm := <-newMap
	return types.NewStruct("HackerNoms", types.StructData{
		"items": mm,
		"top":   types.NewList(types.Number(0)),
	})
}

func workerPool(count int, work func(), done func()) {
	workerDone := make(chan bool, 1)
	for i := 0; i < count; i += 1 {
		go func() {
			work()
			workerDone <- true
		}()
	}

	go func() {
		for i := 0; i < count; i += 1 {
			_ = <-workerDone
		}
		close(workerDone)
		done()
	}()
}

func makeClient() *http.Client {
	var tr *http.Transport
	tr = &http.Transport{
		Dial: func(network, address string) (net.Conn, error) {
			return net.DialTimeout(network, address, 30*time.Second)
		},
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	return client
}

func churn(newIndex <-chan float64, newData chan<- datum) {
	client := makeClient()

	for index := range newIndex {
		id := int(index)
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d", id)
		fb := firego.New(url, client)
		bytes, err := fb.DoRequest("GET", nil)
		if err != nil {
			fmt.Printf("DoRequest failed for %d %s\n", id, err)
		}
		n := len(bytes)
		var dst []byte
		dst = make([]byte, n, n)
		numofbytes := copy(dst,bytes)
		fmt.Println("Number of bytes copied = ",numofbytes)
		name, data := processBytes(id, dst)
		fmt.Println(id, name)
		// elastic.Process_json_bytes("hackernews", name, string(id), bytes)
		// send data to redis
		sendDatum(newData, name, index, data)
	}
}

func processBytes(id int, bytes []byte) (name string, data map[string]types.Value) {
	var val map[string]interface{}
	err := json.Unmarshal(bytes, &val)
	if err != nil {
		fmt.Printf("json Unmarshal failed for %d %s\n", id, err)
	}

	data = make(map[string]types.Value)
	for k, v := range val {
		switch vv := v.(type) {
		case string:
			data[k] = types.String(vv)
		case float64:
			data[k] = types.Number(vv)
		case bool:
			data[k] = types.Bool(vv)
		case []interface{}:
			ll := types.NewList()
			for _, elem := range vv {
				ll = ll.Append(types.Number(elem.(float64)))
			}
			data[k] = ll
		default:
			panic(reflect.TypeOf(v))
		}
	}

	name, ok := val["type"].(string)
	if !ok {
		fmt.Printf("no type for id %d; trying again\n", id)
	}
	return name, data
}

func sendDatum(newData chan<- datum, name string, id float64, data map[string]types.Value) {
	st := types.NewStruct(name, data)
	d := datum{
		index: id,
		value: st,
	}

	newData <- d
}
