package main

// considered harmful
import "server"

import "log"
import "fmt"
import "net/rpc"

func main() {

	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &server.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith Multiply: %d*%d=%d\n", args.A, args.B, reply)
	//fmt.Printf("Arith Divide: %d", replyCall)
}
