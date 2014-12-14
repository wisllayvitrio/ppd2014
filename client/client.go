package client

import (
	"fmt"
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type Stub struct {
	mid middleware.Middleware
}

// Test function sum
func (s *Stub) Sum(a,b int) (int, error) {
	// Create Request
	args := []interface{}{interface{}(a), interface{}(b)}
	req := Request{"testServ", "Sum", "666", args}
	
	// Send to tuple space
	err := middleware.Send(req)
	
	// Wait for response
	// TODO
	
	// Return response
	return 0, err
}