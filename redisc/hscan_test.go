package redisc

import (
	"testing"
)

func TestHscan(t *testing.T) {
	Hscan("story")
	Hscan("comment")
}
