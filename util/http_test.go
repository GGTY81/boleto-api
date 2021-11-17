package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderToMap(t *testing.T) {
	h := getHeader()

	result := HeaderToMap(h)

	assert.Equal(t, len(h), len(result))
	assert.Equal(t, result["Authorization"], "***")
}

func BenchmarkHeaderToMap(b *testing.B) {
	h := getHeader()
	for n := 0; n < b.N; n++ {
		HeaderToMap(h)
	}
}

func getHeader() http.Header {
	h := make(http.Header)
	h.Add("Authorization", "Basic NWVjNWI5zBkOmFiNjE3YzM4LTliNzAtNGE1OS1hMzhmLTMzMTU0ZmFiMDEwYw==")
	h.Add("PostmanToken", "7dc1c2cd-4a22-49ed-85f8-f081e3cb25dc")
	h.Add("AcceptEncoding", "gzip, deflate, br")
	h.Add("ContentLength", "1933")
	h.Add("ContentType", "application/json")
	h.Add("Connection", "no-cache")
	h.Add("CacheControl", "keep-alive")
	h.Add("Accept", "*/*")
	h.Add("UserAgent", "PostmanRuntime/7.28.4")
	return h
}
