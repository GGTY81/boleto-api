package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderToMap(t *testing.T) {
	header := getHeader()

	result := HeaderToMap(header)

	assert.Equal(t, len(header), len(result))
	assert.Equal(t, result["Authorization"], "***")
}

func BenchmarkHeaderToMap(b *testing.B) {
	header := getHeader()
	for n := 0; n < b.N; n++ {
		HeaderToMap(header)
	}
}

func getHeader() http.Header {
	header := make(http.Header)
	header.Add("Authorization", "Basic NWVjNWI5zBkOmFiNjE3YzM4LTliNzAtNGE1OS1hMzhmLTMzMTU0ZmFiMDEwYw==")
	header.Add("PostmanToken", "7dc1c2cd-4a22-49ed-85f8-f081e3cb25dc")
	header.Add("AcceptEncoding", "gzip, deflate, br")
	header.Add("ContentLength", "1933")
	header.Add("ContentType", "application/json")
	header.Add("Connection", "no-cache")
	header.Add("CacheControl", "keep-alive")
	header.Add("Accept", "*/*")
	header.Add("UserAgent", "PostmanRuntime/7.28.4")
	return header
}
