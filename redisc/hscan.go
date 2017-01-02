package redisc

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func hscan(key string) error {

	var (
		count int
		cursor int64
		items  []string
	)

	//results := make([]string, 0)

	c := getRedisConn()
	defer c.Close()

	for {
		values, err := redis.Values(c.Do("HSCAN", key, cursor))
		if err != nil {
			fmt.Println("hscan error on redis.Values")
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			fmt.Println("hscan error on redis.Scan")
		}

		fmt.Println(items)
		fmt.Println(count)
		count = count + 1

		if cursor == 0 {
			break
		}
	}
	return nil
}
