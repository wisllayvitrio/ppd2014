package main

import (
	"os"
	"fmt"
	"net"
	"net/rpc"
	"github.com/wisllayvitrio/ppd2014/space"
)

func main() {
	addr := ":8666"
	space := new(space.TupleSpace)
	
	// Register in the RPC default server
	rpc.Register(space)
	
	// Create a TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("ERROR -", err)
		os.Exit(1)
	}
	
	// Accept and handle connections
	fmt.Println("DEBUG: Listening in", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Just ignore this connection
			continue
		}
		
		// Handle the RPC request
		go rpc.ServeConn(conn)
	}
}