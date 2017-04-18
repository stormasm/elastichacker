package main

import (
	"github.com/stormasm/elastichacker/redisc"
)

func main() {
	redisc.Hscan("story")
	redisc.Hscan("comment")

}
