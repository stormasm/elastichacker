package redisc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/attic-labs/noms/go/hash"
	"strings"
	"github.com/garyburd/redigo/redis"
)

type P struct {
	Itype string
	Id    int
	Json  []byte
}

func Process_json_test(index, itype string, id int) error {
	c := getRedisConn()
	defer c.Close()

	_, err := c.Do("HSET", index, id, itype)
	return err
}

func Write_json_bytes(index, itype string, id int, byteArray []byte) error {
	c := getRedisConn()
	defer c.Close()

	// Write the byteArray using the elastic type as the key
	// and the hackernews id as the field

	_, err := c.Do("HSET", itype, id, byteArray)

	// Encode and write the struct as GOB data to redis as well...

	nbytearray := encode_struct_tobytes(itype, id, byteArray)
	hashString := hash.Of(nbytearray).String()
	_, err = c.Do("HSET", index, id, nbytearray)

	strary := []string{index, "hash"}
	indexhash := strings.Join(strary, "")
	strary = []string{index, "set"}
	indexset := strings.Join(strary, "")

	_, err = c.Do("HSET", indexhash, id, hashString)
	_, err = c.Do("SADD", indexset, id)

	return err
}

func encode_struct_tobytes(itype string, id int, byteArray []byte) []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(P{itype, id, byteArray})
	if err != nil {
		fmt.Println("process_bytes error in Encoder")
	}
	return buf.Bytes()
}

func Read_json_bytes(index string, id int) error {
	c := getRedisConn()
	defer c.Close()

	myinterface, err := c.Do("HGET", index, id)
	if err != nil {
		fmt.Println("Read_json_bytes redis hget error")
	}

	byteary := myinterface.([]byte)
	decode_bytes_to_struct(byteary)
	return nil
}

func decode_bytes_to_struct(byteArray []byte) {
	bytebuf := bytes.NewBuffer(byteArray)
	dec := gob.NewDecoder(bytebuf)

	var p P
	err := dec.Decode(&p)
	if err != nil {
		fmt.Println("decode error:", err)
	}

	fmt.Printf("Id = %d\n", p.Id)
	fmt.Printf("Id = %s\n", p.Itype)

	n := len(p.Json)
	json := string(p.Json[:n])
	fmt.Println(json)
}

func Read_hash_of_struct(index string, id int) (myhash string) {
	c := getRedisConn()
	defer c.Close()

	strary := []string{index, "hash"}
	indexhash := strings.Join(strary, "")

	s, err := redis.String(c.Do("HGET", indexhash, id))
	myinterface, err := c.Do("HGET", indexhash, id)

	if err != nil {
		fmt.Println("Read_hash_of_struct hget error")
	}

	// do a type assertion to convert the interface to a byte array
	byteary := myinterface.([]byte)
	n := len(byteary)
	myhash = string(byteary[:n])

	mycompare := strings.Compare(s,myhash)
	if mycompare == 0 {
		fmt.Println("Strings are equal ", s, myhash)
	}

	return myhash
}
