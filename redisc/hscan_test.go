package redisc

import (
	"testing"
)

func TestHscan(t *testing.T) {

	newStory := make(chan Datum, 100)
	newComment := make(chan Datum, 100)

	Hscan("story", newStory)
	Hscan("comment", newComment)
}
