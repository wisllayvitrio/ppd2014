package middleware

import (
	"os"
	"fmt"
	"net"
	"time"
	"github.com/wisllayvitrio/ppd2014/space"
)

const DefaultReadTimeout string = "300ms"
const DefaultWriteLease string = "1s"

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
	readTimeout time.Duration
	writeLease time.Duration
}

// Constructors
func NewMiddleware(address string) (*Middleware, error) {
	return NewMiddleware(address, DefaultReadTimeout, DefaultWriteLease)
}

func NewMiddleware(address string, timeout string, lease string) (*Middleware, error) {
	m := new(Middleware)
	
	// Set the TupleSpace addr (check before)
	addrs, err := net.LookupHost(address)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: TupleSpace has these addresses:", addrs)
	m.spaceAddr = addrs[0]
	
	// Set the read timeout configuration (used when reading from TupleSpace)
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: ReadTimeout configured:", dur)
	m.readTimeout = dur
	
	// Set the write lease configuration (used when writing to the TupleSpace)
	dur, err = time.ParseDuration(lease)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: WriteLease configured:", dur)
	m.writeLease = dur
	
	return m, nil
}

// Function to change the read timeout configuration
func (m *Middleware) SetReadTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	m.readTimeout = dur
	return nil
}

// Function to change the write lease configuration
func (m *Middleware) SetWriteLease(lease string) error {
	dur, err := time.ParseDuration(lease)
	if err != nil {
		return nil, err
	}
	m.writeLease = dur
	return nil
}

// Send to TupleSpace
func (m *Middleware) SendRequest(req Request) error {
	// Create a tuple
	tuple, err := space.NewTuple(req.ServerName, req.FuncName, req.ResponseID, req.Args)
	if err != nil {
		return err
	}
	
	// Talk to the TupleSpace
	err = communicate("TupleSpace.Write", m.writeLease, tuple, nil)
	if err != nil {
		return err
	}
	
	// No error
	return nil
}

func (m *Middleware) SendResponse(res Response) error {
	tuple, err := space.NewTuple(res.ID, res.Args)
	if err != nil {
		return err
	}
	
	err = communicate("TupleSpace.Write", m.writeLease, tuple, nil)
	if err != nil {
		return err
	}
	
	return nil
}

// Read from TupleSpace
func (m *Middleware) ReceiveResponse(id string) (*Response, error) {
	tuple, err := space.NewTuple(id, nil)
	if err != nil {
		return nil, err
	}
	
	// Create response struct
	res := new(Response)
	
	err = communicate("TupleSpace.Take", m.readTimeout, tuple, res)
	if err != nil {
		return nil, err
	}
	
	return res, nil
}

func (m *Middleware) ReceiveRequest(serverName string) (*Request, error) {
	tuple, err := space.NewTuple(serverName, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	
	req := new(Request)
	
	err = communicate("TupleSpace.Take", m.readTimeout, tuple, req)
	if err != nil {
		return nil, err
	}
	
	return req, nil
}

// Get a request and call the worker function to execute it
func (m *Middleware) Serve(obj interface{}, serviceName string) error {
	req, err := ReceiveRequest(serviceName)
	if err != nil {
		return err
	}
	
	results, err := invoke(obj, req.FuncName, req.Args)
	if err != nil {
		return err
	}
	
	res := Response{req.ResponseID, results}
	err = SendResponse(res)
	if err != nil {
		return err
	}
	
	return nil
}

func communicate(call string, t time.Duration, req space.Tuple, res *space.Tuple) error {
	// Create request
	message := space.Request{req, t}
	
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

func invoke(obj interface{}, methodName string, args []interface{}) ([]interface{}, error) {
	inputs := make([]reflect.Value, len(args))
	for i, v := range args {
		inputs[i] = reflect.ValueOf(v)
	}
	
	// Get the method (of type []reflect.Value)
	objValue := reflect.ValueOf(obj)
	method := objValue.MethodByName(methodName)
	
	// Test if method was found
	if !method.IsValid() {
		return nil, errors.New(fmt.Sprintf("%s was not found on object of type %s", methodName, objValue.Type()))
	}
	
	// Call the method (outputs is of type []reflect.Value)
	outputs := method.Call(inputs)
	// Get the []interface{} of the outputs
	result := make([]interface{}, len(outputs))
	for i, v := range outputs {
		result[i] = v.Interface()
	}
	
	return result, nil
}