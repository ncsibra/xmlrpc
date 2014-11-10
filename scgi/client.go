package scgi

import (
	"fmt"
	"net"
	"net/rpc"

	"bytes"

	"io"

	"github.com/kolo/xmlrpc"
)

type requestData struct {
	seq           uint64
	serviceMethod string
}

type scgiClientCodec struct {
	conn io.ReadWriteCloser
	seq  chan requestData
}

func (codec *scgiClientCodec) WriteRequest(request *rpc.Request, args interface{}) error {
	params, ok := args.([]interface{})
	if !ok {
		if args != nil {
			params = []interface{}{args}
		}
	}

	body, err := xmlrpc.EncodeMethodCall(request.ServiceMethod, params...)
	if err != nil {
		return err
	}

	req, err := NewRequest(body, nil)
	if err != nil {
		return nil
	}

	n, err := codec.conn.Write(append(req.Netstring, req.Body...))
	if err != nil {
		return err
	}
	fmt.Println("Request:", request.Seq, request.ServiceMethod, n, string(append(req.Netstring, req.Body...)))
	codec.seq <- requestData{request.Seq, request.ServiceMethod}

	return nil
}

func (codec *scgiClientCodec) ReadResponseHeader(response *rpc.Response) error {
	fmt.Println("ReadResponseHeader")
	d, ok := <-codec.seq
	if !ok {
		return nil
	}

	response.Seq = d.seq
	response.ServiceMethod = d.serviceMethod
	fmt.Println("Response:", response.Seq, response.ServiceMethod)

	return nil
}

func (codec *scgiClientCodec) ReadResponseBody(body interface{}) error {
	fmt.Println("ReadResponseBody nil")
	if body == nil {
		return nil
	}

	fmt.Println("ReadResponseBody non nil")

	var r bytes.Buffer

	//	codec.conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, err := r.ReadFrom(codec.conn)
	fmt.Println("Read finished")
	if err != nil {
		fmt.Println("Conn failed")
		return err
	}

	fmt.Println("Readed response:", string(r.Bytes()))
	res := xmlrpc.NewResponse(r.Bytes())
	if res.Failed() {
		fmt.Println("Resp failed")
		return res.Err()
	}

	fmt.Println("Unmarshal")
	return res.Unmarshal(body)
}

func (codec *scgiClientCodec) Close() error {
	fmt.Println("Close")
	close(codec.seq)
	return codec.conn.Close()
}

func NewScgiClient(url string) (*rpc.Client, error) {
	c, err := net.Dial("tcp", url)

	if err != nil {
		return nil, err
	}

	scc := scgiClientCodec{
		conn: c,
		seq:  make(chan requestData),
	}

	return rpc.NewClientWithCodec(&scc), nil
}
