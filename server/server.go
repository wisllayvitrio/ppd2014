package server

import (
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type Service struct {
	name string
	mid middleware.Middleware
}

// Work once to a random stranger (for free)
func (s *Service) WorkDefault() error {
	err = mid.Serve(s, s.name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Work(waitTimeout string) error {
	err := mid.SetReadTimeout(waitTimeout)
	if err != nil {
		return err
	}
	
	err = mid.Serve(s, s.name)
	if err != nil {
		return err
	}
	return nil
}

// Actual service function
func (s *Service) Sum(a, b int) int {
	return a + b
}