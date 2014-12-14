package client

import (
	"fmt"
	"net/rpc"
	"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/space"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
}

// TODO: Where should this be?
type Request struct {
	// Demux header
	serverName string
	funcName string
	responseID string
	// Data for the actual execution
	args []interface{}
}

// Send to tuplespace
func spaceSend(req Request) {
	// TODO: put this in the correct location (command arguments)
	addr := "localhost:8666"

	// Create an tuple
	tuple, err := space.NewTuple(req.serverName, req.funcName, req.responseID, req.args)
	checkErr(err)
	
	// Dial the RPC server
	rpcClient, err := rpc.Dial("tcp", addr)
	checkErr(err)
	
	// Call the write function of the TupleSpace
	var ok bool
	err = rpcClient.Call("TupleSpace.Write", tuple, &ok)
	checkErr(err)
	// Call done
	fmt.Println("Everything ok?", ok)
}

type Stub struct {
}

// Test function sum
func (s *Stub) Sum(a,b int) (int, error) {
	// Create Request
	args := []interface{}{interface{}(a), interface{}(b)}
	req := Request{"testServ", "Sum", "666", args}
	
	// Send to tuple space
	spaceSend(req)
	
	// Wait for response
	// TODO
	
	// Return response
	return 0, nil
}