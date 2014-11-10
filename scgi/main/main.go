package main

import (
	"fmt"

	"github.com/kolo/xmlrpc/scgi"
)

func main() {
	client, _ := scgi.NewScgiClient("localhost:5000")
	var result string

	err := client.Call("system.client_version", nil, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %v\n", result)

	//	c := make(chan *rpc.Call, 3)
	//	go client.Go("system.client_version", nil, &result, c)
	//	go client.Go("system.client_version", nil, &result, c)
	//	go client.Go("system.client_version", nil, &result, c)
	//
	//	for call := range c {
	//		if call.Error != nil {
	//			fmt.Println(call.Error)
	//			client.Close()
	//			break
	//		}
	//		fmt.Printf("Result: %v\n", result)
	//	}

}
