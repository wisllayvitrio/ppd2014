package server

import (
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type Calculator struct {
	name string
	mid middleware.Middleware
}

func NewCalculator(name string, spaceAddr string) (*Calculator, error) {
	c := new(Calculator)
	ptr, err := middleware.NewMiddlewareDefault(spaceAddr)
	if err != nil {
		return nil, err
	}
	
	c.name = name
	c.mid = *ptr
	return c, nil
}

func (c *Calculator) Name() string {
	return c.name
}

// Work once to a random stranger (for free)
func (c *Calculator) WorkDefault() error {
	err := c.mid.Serve(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Calculator) Work(waitTimeout string, replyLeasing string) error {
	err := c.mid.SetReadTimeout(waitTimeout)
	if err != nil {
		return err
	}
	
	err = c.mid.SetWriteLeasing(replyLeasing)
	if err != nil {
		return err
	}
	
	err = c.mid.Serve(c)
	if err != nil {
		return err
	}
	return nil
}

// Actual service function
func (c *Calculator) Sum(a, b int) int {
	return a + b
}