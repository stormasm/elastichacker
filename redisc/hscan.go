package redisc

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Datum struct {
	hnid string    // hackernews id
	json string
}

func Hscan(key string, newData chan<- Datum) error {

	var (
		hackernewsid string
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
				hackernewsid = item
				_, err = c.Do("SADD", setkeyname, item)
				if err != nil {
					fmt.Println("error on SADD")
				}
			}
			if evenodd == 1 {
				//fmt.Println(hackernewsid)
				//fmt.Println(item)
				// Build the struct here and put it on a channel
				c := Datum{
					hnid:hackernewsid,
					json:item,
				}
				newData <- c
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
