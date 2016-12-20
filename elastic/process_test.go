package elastic

import (
	"testing"
//	"github.com/stretchr/testify/assert"
)

func TestProcessJsonString(t *testing.T) {
	tweet2 := `{"user" : "olivere", "message" : "It's not a Raggy Waltz"}`
    Process_json_string("berk",tweet2)
}

func TestProcessJsonByteArray(t *testing.T) {
	tweet2 := []byte(`{"user" : "olivere", "message" : "It's not a Raggy Waltz"}`)
    Process_json_bytes("angerman",tweet2)
}
