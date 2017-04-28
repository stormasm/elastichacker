package main

import (
	"fmt"
	"github.com/stormasm/elastichacker/redisc"
	"time"
)

func printer(c chan redisc.Datum) {
	for {
		msg := <-c
		fmt.Println(msg)
		fmt.Println()
		time.Sleep(time.Second * 1)
	}
}

func main() {

	var newStory chan redisc.Datum = make(chan redisc.Datum, 100)
	var newComment chan redisc.Datum = make(chan redisc.Datum, 100)

	go redisc.Hscan("story", newStory)
	go redisc.Hscan("comment", newComment)
	go printer(newStory)
	go printer(newComment)

	var input string
	fmt.Scanln(&input)
}
