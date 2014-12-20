package main

import (
	"fmt"
	"log"
	"net"
	"flag"
	"net/rpc"
	"runtime"
	"github.com/wisllayvitrio/ppd2014/space"
)

// Command-line flags
var addr string
// Default values
const defaultAddr string = ":8666"
// Descriptions
const usageAddr string = "IP:PORT to receive connections (empty IP means all interfaces)"
// Set the flag names (long and short for each flag var
func init() {
	flag.StringVar(&addr, "address", defaultAddr, usageAddr)
	flag.StringVar(&addr, "a", defaultAddr, usageAddr)
}

func main() {
	flag.Parse()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	
	// Create the TupleSpace and register in the RPC default server
	space := space.NewTupleSpace()
	rpc.Register(space)
	
	// Create a TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("ERROR creating TCP listener:", err)
	}
	
	// Accept and handle connections
	fmt.Println("DEBUG - Listening in:", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("DEBUG - Ignoring failed connection. Err:", err)
			// Just ignore this connection
			continue
		}
		
		// Handle the RPC request
		go rpc.ServeConn(conn)
	}
}