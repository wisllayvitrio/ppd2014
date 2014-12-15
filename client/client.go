package client

import (
	"fmt"
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type Stub struct {
	mid middleware.Middleware
}

func NewStub(spaceAddr string) (*Stub, error) {
	s := new(Stub)
	ptr, err := middleware.NewMiddleware(spaceAddr)
	if err != nil {
		return nil, err
	}
	
	s.mid = *ptr
	return s, nil
}

// Test function sum
func (s *Stub) Sum(a,b int) (int, error) {
	// Create Request
	args := []interface{}{interface{}(a), interface{}(b)}
	req := middleware.Request{"testServ", "Sum", "666", args}
	
	// Set lease
	mid.SetWriteLease("10s")
	// Send to tuple space
	err := mid.SendRequest(req)
	if err != nil {
		return -1, err
	}
	
	// Set timeout
	mid.SetReadTimeout("1s")
	// Wait for response
	res, err := mid.ReceiveResponse("666")
	if err != nil {
		return -1, err
	}
	
	// This function understands middleware.Response
	var values []interface{}
	values = res.Args
	result := int(values[0])
	
	// Return response
	return result, err
}