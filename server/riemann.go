package server

import (
	"fmt"
	"math"
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type RiemannCalculator struct {
	name string
	maxExec int
	m middleware.Middleware
}

func NewRiemannCalculator(spaceAddr, timeout, leasing string, maxExec int) (*RiemannCalculator, error) {
	r := new(RiemannCalculator)
	ptr, err := middleware.NewMiddleware(spaceAddr, timeout, leasing)
	if err != nil {
		return nil, err
	}
	
	r.name = "Riemann"
	r.maxExec = maxExec
	r.m = *ptr
	return r, nil
}

func (r *RiemannCalculator) Name() string {
	return r.name
}

func (r *RiemannCalculator) Work() error {
	// Create maxExec goroutines to execute
	done := make(chan error, r.maxExec)
	for i := 0; i < r.maxExec; i++ {
		go r.execute(done)
	}
	// Continue executing forever
	for {
		if err := <-done; err != nil {
			fmt.Println("WORK ERROR:", err)
		}
		go r.execute(done)
	}
	
	return nil
}

func (r *RiemannCalculator) execute(done chan error) {
	err := r.m.Serve(r)
	if err != nil {
		done<- err
		return
	}
	
	done<- nil
	return
}

func (r *RiemannCalculator) SetReadTimeout(timeout string) error {
	err := r.m.SetReadTimeout(timeout)
	if err != nil {
		return err
	}
	return nil
}
func (r *RiemannCalculator) SetWriteLeasing(leasing string) error {
	err := r.m.SetWriteLeasing(leasing)
	if err != nil {
		return err
	}
	return nil
}

// Riemann Integral calculation
func (r *RiemannCalculator) Integral(x, y, dx float64, coefs []float64) float64 {
	sum := 0.0
	for i := x; i < y; i += dx {
		sum += dx * poly(i, coefs)
	}
	return sum
}

func poly(x float64, coefs []float64) float64 {
	y := 0.0
	for i, coef := range coefs {
		y += coef * math.Pow(x, float64(i))
	}
	return y
}