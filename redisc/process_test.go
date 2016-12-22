package redisc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadJsonBytes(t *testing.T) {
	var index string = "hackernews"
	var id int = 8432709
	Read_json_bytes(index, id)
}

func TestReadJsonStruct(t *testing.T) {
	var index string = "hackernews"
	var id int = 8432709

	// Read struct out of redis
	// Read hash out of redis
	// Test hash of struct

	myhash := Read_hash_of_struct(index, id)
	assert := assert.New(t)
	s1 := "67to51ntpmub261mlapf31jvdos04gk5"
	assert.Equal(myhash, s1)
}
