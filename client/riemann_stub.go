package client

import (
	"time"
	"code.google.com/p/go-uuid/uuid"
	"github.com/wisllayvitrio/ppd2014/logger"
	"github.com/wisllayvitrio/ppd2014/middleware"
)

type execRes struct {
	res []interface{}
	err error
}

type RiemannStub struct {
	name string
	m middleware.Middleware
	l *logger.Logger
}

func NewRiemannStub(spaceAddr, timeout, leasing string) (*RiemannStub, error) {
	r := new(RiemannStub)
	ptr, err := middleware.NewMiddleware(spaceAddr, timeout, leasing)
	if err != nil {
		return nil, err
	}
	
	r.name = "Riemann"
	r.m = *ptr
	
	r.l = logger.NewLogger("ppd2014_stub_log.txt", time.Second)
	go r.l.LogStart()
	
	return r, nil
}

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
		aux := time.Now()
		go r.m.SendRequest(req)
		r.l.AddTime(false, time.Since(aux))}
	}

	// Wait, get each response, and calculate the final value
	sum := 0.0
	errCount := 0
	for i := 0; i < numParts; i++ {
		aux := time.Now()
		res, err := r.m.ReceiveResponse(id)
		r.l.AddTime(true, time.Since(aux))
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