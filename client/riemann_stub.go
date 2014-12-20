package client

import (
	"time"
	//"code.google.com/p/go-uuid/uuid"
	"github.com/wisllayvitrio/ppd2014/middleware"
	"github.com/wisllayvitrio/ppd2014/logger"
)

type execRes struct {
	res []interface{}
	err error
}

type RiemannStub struct {
	name string
	m middleware.Middleware
	L *logger.Logger
}

func NewRiemannStub(spaceAddr, timeout, leasing string) (*RiemannStub, error) {
	r := new(RiemannStub)
	ptr, err := middleware.NewMiddleware(spaceAddr, timeout, leasing)
	if err != nil {
		return nil, err
	}
	
	r.name = "Riemann"
	r.m = *ptr
	r.L = logger.NewLogger()
	return r, nil
}
/**/
func (r *RiemannStub) Integral(a, b, dx float64, coefs []float64, numParts int) (float64, int, error) {
	funcName := "Integral"
	done := make(chan execRes, numParts)
	partDelta := (b-a)/float64(numParts)
	
	// Call the execution of each part
	for i := 0; i < numParts; i++ {
		// Calculate and set this part args
		args := make([]interface{}, 4)
		args[0] = interface{}(a + float64(i) * partDelta)
		args[1] = interface{}(a + float64(i) * partDelta + partDelta)
		args[2] = interface{}(dx)
		args[3] = interface{}(coefs)
		
		go r.execute(r.name, funcName, args, done)
	}
	
	// Wait for each part to execute
	sum := 0.0
	errCount := 0
	for i := 0; i < numParts; i++ {
		exeRes := <-done
		if exeRes.err != nil {
			errCount++
			continue
		}
		v, ok := exeRes.res[0].(float64)
		if ok {
			sum += v
		} else {
			errCount++
		}
	}
	
	return sum, errCount, nil
}

func (r *RiemannStub) execute(name, funct string, args []interface{}, done chan execRes) {
	aux := time.Now()
	res, err := r.m.Execute(name, funct, args)
	r.L.AddTime("execute", time.Since(aux))
	
	if err != nil {
		done<- execRes{nil, err}
		return
	}
	
	done<- execRes{res, nil}
}
/**/
/*
func (r *RiemannStub) Integral(a, b, dx float64, coefs []float64, numParts int) (float64, int, error) {
	// Create the header arguments (same for every part)
	funcName := "Integral"
	id := uuid.NewRandom().String()
	partDelta := (b-a)/float64(numParts)
	
	// Prepare the request and send for each part
	for i := 0; i < numParts; i++ {
		req := middleware.Request{}
		req.ServiceName = r.name
		req.FuncName = funcName
		req.ResponseID = id
		req.Args = make([]interface{}, 4)
		// Calculate and set this part args
		req.Args[0] = interface{}(a + float64(i) * partDelta)
		req.Args[1] = interface{}(a + float64(i) * partDelta + partDelta)
		req.Args[2] = interface{}(dx)
		req.Args[3] = interface{}(coefs)
		
		// Send the request of this part
		err := r.m.SendRequest(req)
		if err != nil {
			return 0.0, 0, err
		}
	}

	// Wait, get each response, and calculate the final value
	sum := 0.0
	errCount := 0
	for i := 0; i < numParts; i++ {
		res, err := r.m.ReceiveResponse(id)
		if err != nil {
			errCount++
			continue
		}
		// Guarantee that Args[0] is a float64
		v, ok := res.Args[0].(float64)
		if ok {
			sum += v
		} else {
			errCount++
		}
	}
	
	return sum, errCount, nil
}
/**/