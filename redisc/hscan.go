package redisc

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
)

func hscan(key string) error {

	var (
		setkeyname string
		total  int
		count  int
		cursor int64
		items  []string
	)

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

		// fmt.Println("items length = ", len(items))

		strary := []string{"set", key}
		setkeyname = strings.Join(strary, "")

		for num, item := range items {
			evenodd := num % 2
			// Grab the ID
			if evenodd == 0 {
				_, err = c.Do("SADD", setkeyname, item)
				if err != nil {
					fmt.Println("error on SADD")
				}
			}
		}

		// fmt.Println("count = ", count)
		total = total + len(items)
		count = count + 1
		if cursor == 0 {
			break
		}
	}
	fmt.Println(setkeyname, " total = ", total/2)
	return nil
}
