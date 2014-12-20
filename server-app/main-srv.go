package main

import (
	"log"
	//"fmt"
	"flag"
	"runtime"
	"github.com/wisllayvitrio/ppd2014/server"
)

// Command-line flags
var addr string
var timeout string
var leasing string
var maxExec int
// Default values
const defaultAddr string = "localhost:8666"
const defaultTimeout string = "10s"
const defaultLeasing string = "5s"
const defaultMaxExec int = 10
// Descriptions
const usageAddr string = "IP:PORT of the Tuple Space"
const usageTimeout string = "Time to wait for messages from the Tuple Space"
const usageLeasing string = "Leasing time of messages sent to the Tuple Space"
const usageMaxExec string = "Max number of goroutines to use at the same time"
// Set the flag names (long and short for each flag var
func init() {
	flag.StringVar(&addr, "address", defaultAddr, usageAddr)
	flag.StringVar(&addr, "a", defaultAddr, usageAddr)
	
	flag.StringVar(&timeout, "timeout", defaultTimeout, usageTimeout)
	flag.StringVar(&timeout, "t", defaultTimeout, usageTimeout)
	
	flag.StringVar(&leasing, "leasing", defaultLeasing, usageLeasing)
	flag.StringVar(&leasing, "l", defaultLeasing, usageLeasing)
	
	flag.IntVar(&maxExec, "goroutines", defaultMaxExec, usageMaxExec)
	flag.IntVar(&maxExec, "g", defaultMaxExec, usageMaxExec)
}

func main() {
	flag.Parse()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	
	r, err := server.NewRiemannCalculator(addr, timeout, leasing, maxExec)
	if err != nil {
		log.Fatalln("ERROR creating RiemannCalculator:", err)
	}
	
	// Keep working forever
	err = r.Work()
	if err != nil {
		log.Fatalln("ERROR working:", err)
	}
}