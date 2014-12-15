package client

import (
	"github.com/wisllayvitrio/ppd2014/middleware"
	"code.google.com/p/go-uuid/uuid"
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
	id := uuid.NewRandom().String()
	args := []interface{}{interface{}(a), interface{}(b)}
	req := middleware.Request{"testServ", "Sum", id, args}
	
	// Set leasing
	s.mid.SetWriteLeasing("10s")
	// Send to tuple space
	err := s.mid.SendRequest(req)
	if err != nil {
		return -1, err
	}
	
	// Set timeout
	s.mid.SetReadTimeout("10s")
	// Wait for response
	res, err := s.mid.ReceiveResponse(id)
	if err != nil {
		return -1, err
	}
	
	// This function understands middleware.Response
	var values []interface{}
	values = res.Args
	result := int(values[0].(int))
	
	// Return response
	return result, err
}