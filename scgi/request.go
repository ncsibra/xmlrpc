package scgi

import (
	"bytes"
	"strconv"
)

const (
	comma = byte(',')
	colon = byte(':')
)

var (
	contentLength = []byte("CONTENT_LENGTH")
	scgi          = []byte("SCGI")
	one           = []byte("1")
	nul           = []byte{0}
)

type Request struct {
	Netstring []byte
	Body      []byte
}

func NewRequest(body []byte, extraHeaders map[string]string) (*Request, error) {
	req := new(Request)
	req.Body = body
	h := req.createRequiredHeaders()
	h = req.addExtraHeaders(h, extraHeaders)
	req.Netstring = req.createNetstring(h)

	return req, nil
}

func (req Request) createRequiredHeaders() [][]byte {
	cl := []byte(strconv.Itoa(len(req.Body)))
	return [][]byte{contentLength, cl, scgi, one}
}

func (req Request) addExtraHeaders(headers [][]byte, extraHeaders map[string]string) [][]byte {
	for k, v := range extraHeaders {
		headers = append(headers, []byte(k))
		headers = append(headers, []byte(v))
	}

	return headers
}

func (req *Request) createNetstring(headers [][]byte) []byte {
	var ns []byte
	h := append(bytes.Join(headers, nul), nul...)

	ns = append(ns, []byte(strconv.Itoa(len(h)))...)
	ns = append(ns, colon)
	ns = append(ns, h...)
	ns = append(ns, comma)

	return ns
}
