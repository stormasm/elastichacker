package elastic

import (
	"testing"
	//	"github.com/stretchr/testify/assert"
)

func TestProcessJsonString(t *testing.T) {
	tweet2 := `{"user" : "olivere", "message" : "poor man"}`
	Process_json_string("lascruces", "redis", "102", tweet2)
}

func TestProcessJsonByteArray(t *testing.T) {
	tweet2 := []byte(`{"user" : "olivere", "message" : "sunshine clouds"}`)
	Process_json_bytes("oregon", "raton", "103", tweet2)
}
