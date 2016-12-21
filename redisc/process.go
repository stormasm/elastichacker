package redisc

import (
//	"github.com/garyburd/redigo/redis"
)

func Process_json_bytes(index, itype, id string, byteArray []byte) {
}

func Process_json_test(index, itype string, id int) error {
    c := getRedisConn()
    defer c.Close()

    _, err := c.Do("HSET", index, id, itype)
    return err
}
