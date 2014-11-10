package scgi

import (
	"fmt"
	"io"
	"net/rpc"
	"testing"
)

func newClient(cli io.ReadWriteCloser) *rpc.Client {
	scc := scgiClientCodec{
		conn: cli,
		seq:  make(chan requestData),
	}
	return rpc.NewClientWithCodec(&scc)
}

func TestClient(t *testing.T) {
	client := newClient(&stubConnection{})
	args := &Args{1, 2}
	reply := new(int)
	err := client.Call("Arith.Add", args, reply)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if 3 != *reply {
		t.Fatalf("Wrong response, expected:%d, actual:%d.", 3, reply)
	}
}

type Args struct {
	A, B int
}

type stubConnection struct {
	io.ReadWriteCloser
	resp []byte
}

func (sc *stubConnection) Read(p []byte) (n int, err error) {
	if len(p) < len(sc.resp) {
		return 0, nil
	}

	copy(p, sc.resp)

	return len(sc.resp), io.EOF
}

func (sc *stubConnection) Write(p []byte) (n int, err error) {
	sc.resp = createResponse(3)

	return len(p), nil
}

func (sc *stubConnection) Close() error {
	return nil
}

func createResponse(n int) []byte {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><methodResponse><params><param><value><i4>%d</i4></value></param></params></methodResponse>`, n)
	return []byte(xml)
}
