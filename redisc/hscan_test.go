package redisc

import (
	"testing"
)

func TestHscan(t *testing.T) {
	hscan("story")
	hscan("comment")
}
