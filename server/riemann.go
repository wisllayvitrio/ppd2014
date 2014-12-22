package server

import (
	"fmt"
	"time"
	"math"
	"github.com/wisllayvitrio/ppd2014/logger"
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type RiemannCalculator struct {
	name string
	maxExec int
	m middleware.Middleware
	l *logger.Logger
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
	
	r.l = logger.NewLogger("ppd2014_service_log.txt", time.Second)
	go r.l.LogStart()
	
	return r, nil
}

// Name() and Work() implements the middleware.Service interface
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
	aux := time.Now()
	err := r.m.Serve(r)
	r.l.AddTime(true, time.Since(aux))
	if err != nil {
		done<- err
		return
	}
	
	done<- nil
}

// Set the middleware timeouts (to allow users of this object to set)
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
	aux := time.Now()
	for i := x; i < y; i += dx {
		sum += dx * poly(i, coefs)
	}
	r.l.AddTime(false, time.Since(aux))
	return sum
}

func poly(x float64, coefs []float64) float64 {
	y := 0.0
	for i, coef := range coefs {
		y += coef * math.Pow(x, float64(i))
	}
	return y
}