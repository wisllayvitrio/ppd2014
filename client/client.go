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
	ptr, err := middleware.NewMiddlewareDefault(spaceAddr)
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
	fmt.Println("DEBUG: Middleware Request:", req)
	// Set leasing
	s.mid.SetWriteLeasing("10s")
	// Send to tuple space
	err := s.mid.SendRequest(req)
	if err != nil {
		return -1, err
	}
	
	// Set timeout
	s.mid.SetReadTimeout("1s")
	// Wait for response
	res, err := s.mid.ReceiveResponse("666")
	if err != nil {
		return -1, err
	}
	fmt.Println("DEBUG: Middleware Response:", res)
	// This function understands middleware.Response
	var values []interface{}
	values = res.Args
	result := int(values[0].(int))
	
	fmt.Println("DEBUG: a+b:", result)
	
	// Return response
	return result, err
}