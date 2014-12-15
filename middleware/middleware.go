package middleware

import (
	"fmt"
	"time"
	"errors"
	"reflect"
	"net/rpc"
	"github.com/wisllayvitrio/ppd2014/space"
)

const DefaultReadTimeout string = "10s"
const DefaultWriteLeasing string = "30s"

type Service interface {
	Name() string
	WorkDefault() error
	Work(timeout string, replyLeasing string) error
}

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
	writeLeasing time.Duration
}

// Constructors
func NewMiddlewareDefault(address string) (*Middleware, error) {
	return NewMiddleware(address, DefaultReadTimeout, DefaultWriteLeasing)
}

func NewMiddleware(address string, timeout string, leasing string) (*Middleware, error) {
	m := new(Middleware)
	// This is not checking the address (net.Dial calls will check that)
	m.spaceAddr = address
	
	// Set the read timeout configuration (used when reading from TupleSpace)
	m.SetReadTimeout(timeout)
	// Set the write leasing configuration (used when writing to the TupleSpace)
	m.SetWriteLeasing(leasing)
	
	return m, nil
}

// Function to change the read timeout configuration
func (m *Middleware) SetReadTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}
	m.readTimeout = dur
	return nil
}

// Function to change the write leasing configuration
func (m *Middleware) SetWriteLeasing(leasing string) error {
	dur, err := time.ParseDuration(leasing)
	if err != nil {
		return err
	}
	m.writeLeasing = dur
	return nil
}

// Send to TupleSpace
func (m *Middleware) SendRequest(req Request) error {
	// Create a tuple
	tuple, err := space.NewTuple(req.ServerName, req.FuncName, req.ResponseID, req.Args)
	if err != nil {
		return err
	}
	
	dummyTuple, err := space.NewTuple()
	if err != nil {
		return err
	}
	
	// Talk to the TupleSpace
	err = m.communicate("TupleSpace.Write", m.writeLeasing, *tuple, dummyTuple)
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
	
	dummyTuple, err := space.NewTuple()
	if err != nil {
		return err
	}
	
	err = m.communicate("TupleSpace.Write", m.writeLeasing, *tuple, dummyTuple)
	if err != nil {
		return err
	}
	
	return nil
}

// Read from TupleSpace
func (m *Middleware) ReceiveResponse(id string) (*Response, error) {
	tuple, err := space.NewTuple(id, space.NilValue())
	if err != nil {
		return nil, err
	}
	
	// Create response tuple to put the result on it
	resTuple, err := space.NewTuple()
	if err != nil {
		return nil, err
	}

	// Send tuple to TupleSpace
	err = m.communicate("TupleSpace.Take", m.readTimeout, *tuple, resTuple)
	if err != nil {
		return nil, err
	}

	// Create response struct
	res := new(Response)
	err = resTuple.Get(0, &res.ID)
	if err != nil {
		return nil, err
	}
	err = resTuple.Get(1, &res.Args)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Middleware) ReceiveRequest(serverName string) (*Request, error) {
	tuple, err := space.NewTuple(serverName, space.NilValue(), space.NilValue(), space.NilValue())
	if err != nil {
		return nil, err
	}
	
	reqTuple, err := space.NewTuple()
	if err != nil {
		return nil, err
	}
	
	err = m.communicate("TupleSpace.Take", m.readTimeout, *tuple, reqTuple)
	if err != nil {
		return nil, err
	}
	
	req := new(Request)
	err = reqTuple.Get(0, &req.ServerName)
	if err != nil {
		return nil, err
	}
	err = reqTuple.Get(1, &req.FuncName)
	if err != nil {
		return nil, err
	}
	err = reqTuple.Get(2, &req.ResponseID)
	if err != nil {
		return nil, err
	}
	err = reqTuple.Get(3, &req.Args)
	if err != nil {
		return nil, err
	}
	
	return req, nil
}

func (m *Middleware) communicate(call string, t time.Duration, req space.Tuple, res *space.Tuple) error {
	// Create request
	message := space.Request{req, t}

	// Dial the RPC server
	rpcClient, err := rpc.Dial("tcp", m.spaceAddr)
	if err != nil {
		return err
	}
	fmt.Println("DEBUG: Communicate", call, "- time:", t, "- tuple:", req)
	// Call the write function of the TupleSpace
	err = rpcClient.Call(call, message, res)
	if err != nil {
		return err
	}

	return nil
}

// Get a request and call the worker function to execute it
func (m *Middleware) Serve(obj Service) error {
	req, err := m.ReceiveRequest(obj.Name())
	if err != nil {
		return err
	}
	
	results, err := invoke(obj, req.FuncName, req.Args)
	if err != nil {
		return err
	}
	
	res := Response{req.ResponseID, results}
	err = m.SendResponse(res)
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