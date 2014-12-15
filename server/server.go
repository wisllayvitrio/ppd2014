package server

import (
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type Service struct {
	name string
	mid middleware.Middleware
}

func NewService(name string, spaceAddr string) (*Service, error) {
	s := new(Service)
	ptr, err := middleware.NewMiddlewareDefault(spaceAddr)
	if err != nil {
		return nil, err
	}
	
	s.name = name
	s.mid = *ptr
	return s, nil
}

// Work once to a random stranger (for free)
func (s *Service) WorkDefault() error {
	err := s.mid.Serve(s, s.name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Work(waitTimeout string) error {
	err := s.mid.SetReadTimeout(waitTimeout)
	if err != nil {
		return err
	}
	
	err = s.mid.Serve(s, s.name)
	if err != nil {
		return err
	}
	return nil
}

// Actual service function
func (s *Service) Sum(a, b int) int {
	return a + b
}