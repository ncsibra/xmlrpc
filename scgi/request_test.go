package scgi

import (
	"bytes"
	"strconv"
	"testing"
)

type testData struct {
	body     string
	headers  map[string]string
	expected Request
}

var requests = []testData{
	createRequest("", nil),
	createRequest("system_clientVersion", nil),
	createRequest("", map[string]string{"a": "b"}),
	createRequest("system_clientVersion", map[string]string{"a": "b"}),
}

func TestValidRequests(t *testing.T) {
	for i, d := range requests {
		req, err := NewRequest([]byte(d.body), d.headers)

		if err != nil {
			t.Fatalf("Error not expected here, index: %d, method: %s, headers: %v, error: %v", i, d.body, d.headers, err)
		}

		if !bytes.Equal(req.Netstring, d.expected.Netstring) {
			t.Fatalf("The actual netstring is not equal with the expected.\n index: %d\n method: %s\n headers: %v\n actstr: %s\n expstr: %s\n act: %v\n exp: %v", i, d.body, d.headers, string(req.Netstring), string(d.expected.Netstring), req.Netstring, d.expected.Netstring)
		}
	}
}

func createRequest(body string, headers map[string]string) testData {
	var ns []byte

	h := []byte("CONTENT_LENGTH")
	h = append(h, 0)
	h = append(h, []byte(strconv.Itoa(len(body)))...)
	h = append(h, 0)
	h = append(h, []byte("SCGI")...)
	h = append(h, 0)
	h = append(h, []byte("1")...)
	h = append(h, 0)

	for k, v := range headers {
		h = append(h, []byte(k)...)
		h = append(h, 0)
		h = append(h, []byte(v)...)
		h = append(h, 0)
	}

	ns = append(ns, []byte(strconv.Itoa(len(h)))...)
	ns = append(ns, []byte(":")...)
	ns = append(ns, h...)
	ns = append(ns, []byte(",")...)

	return testData{body, headers, Request{ns, []byte(body)}}
}
