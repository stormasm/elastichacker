package redisc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"github.com/attic-labs/noms/go/hash"
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

	nbytearray := encode_struct_tobytes(itype, id, byteArray)
	hashString := hash.Of(nbytearray).String()
	_, err := c.Do("HSET", index, id, nbytearray)

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

func Read_hash_of_struct(index string, id int) (myhash string) {
	c := getRedisConn()
	defer c.Close()

    strary := []string{index, "hash"}
	indexhash := strings.Join(strary, "")

	myinterface, err := c.Do("HGET", indexhash, id)

    if err != nil {
        fmt.Println("Read_hash_of_struct hget error")
    }

    // do a type assertion to convert the interface to a byte array
    byteary := myinterface.([]byte)
    n := len(byteary)
    myhash = string(byteary[:n])
    return myhash
}
