package middleware

import (
	"os"
	"fmt"
	"net"
	"time"
	"github.com/wisllayvitrio/ppd2014/space"
)

type Request struct {
	// Demux header
	ServerName string
	FuncName string
	ResponseID string
	// Data for the actual execution
	Args []interface{}
}

type Response struct {
	ID string
	Args []interface{}
}

type Middleware struct {
	spaceAddr string
	timeout time.Duration
}

// Constructors
func NewMiddleware(address string) (*Middleware, error) {
	return NewMiddleware(address, "300ms")
}

func NewMiddleware(address string, timeout string) (*Middleware, error) {
	m := new(Middleware)
	
	// Set the TupleSpace addr (check before)
	addrs, err := net.LookupHost(address)
	if err != nil {
		return nil, error
	}
	fmt.Println("DEBUG: TupleSpace has these addresses:", addrs)
	m.spaceAddr = addrs[0]
	
	// Set the timeout configuration (used when reading from TupleSpace)
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: Timeout configured:", dur)
	m.timeout = dur
	
	return m, nil
}

// Function to change the timeout configuration
func (m *Middleware) SetTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	m.timeout = dur
	return nil
}

// Send to TupleSpace
func (m *Middleware) Send(req Request) error{
	// Create a tuple
	tuple, err := space.NewTuple(req.ServerName, req.FuncName, req.ResponseID, req.Args)
	if err != nil {
		return err
	}
	
	// Talk to the TupleSpace
	err = m.communicate("TupleSpace.Write", tuple, nil)
	if err != nil {
		return err
	}
	
	// No error
	return nil
}

// Read from TupleSpace
func (m *Middleware) Read(id string) (*Response, error) {
	req, err := space.NewTuple(id, nil)
	if err != nil {
		return err
	}
	
	// Create response struct
	res := new(Response)
	
	err = m.communicate("TupleSpace.Take", tuple, res)
	if err != nil {
		return nil, err
	}
	
	return res, nil
}

func (m *Middleware) communicate(call string, req space.Tuple, res *space.Tuple) error {
	// Create request
	message := space.Request{req, m.timeout}
	
	// Dial the RPC server
	rpcClient, err := rpc.Dial("tcp", m.spaceAddr)
	if err != nil {
		return err
	}
	
	// Call the write function of the TupleSpace
	err = rpcClient.Call(call, message, res)
	if err != nil {
		return err
	}
	
	return nil
}