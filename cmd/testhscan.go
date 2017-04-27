package main

import (
	"fmt"
	"github.com/stormasm/elastichacker/redisc"
)

func main() {

	var item redisc.Datum

	newStory := make(chan redisc.Datum, 100)
	newComment := make(chan redisc.Datum, 100)

	go func() {
		item = <-newComment
		fmt.Println(item)
	}()

	redisc.Hscan("story", newStory)
	redisc.Hscan("comment", newComment)

}
