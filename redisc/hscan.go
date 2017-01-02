package redisc

import (
	"fmt"
	"reflect"
	"github.com/garyburd/redigo/redis"
)

func hscan(key string) error {

	var (
		total int
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

		fmt.Println("items length = ",len(items))
		for num, item := range items {
	  		fmt.Println(num)
			fmt.Println(reflect.TypeOf(item))
			fmt.Println(item)

			_, err = c.Do("SADD", "storyset", item[0])
			if err != nil {
				fmt.Println("error on SADD")
			}
  		}

		//fmt.Println(reflect.TypeOf(items))

/*
		fmt.Println(items[0])
		fmt.Println(items[1])
		_, err = c.Do("SADD", "storyset", items[0])
		if err != nil {
			fmt.Println("error on SADD")
		}
*/
		fmt.Println("count = ",count)
		fmt.Println("\n\n")
		total = total + len(items)
		count = count + 1
		if cursor == 0 {
			break
		}
	}
	fmt.Println("total = ", total/2, " ", key)
	return nil
}
