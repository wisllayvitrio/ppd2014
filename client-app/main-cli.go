package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"time"
	"runtime"
	"strings"
	"strconv"
	"go/build"
	"github.com/wisllayvitrio/ppd2014/client"
)

// Polynomial defined as list of coeficients (flag parse use)
type polynomial []float64
// Implement the Value interface (String() and Set() - flag parse use)
func (p *polynomial) String() string {
	var l []string
	if len(*p) == 0 {
		return "[]"
	}
	for _, v := range *p {
		l = append(l, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return ("[" + strings.Join(l, "#") + "]")
}
func (p *polynomial) Set(value string) error {
	// Parse the '#' separated float64 list
	for _, v := range strings.Split(value, "#") {
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*p = append(*p, num)
	}
	return nil
}

// Command-line flags
var addr string
var timeout string
var leasing string
var start float64
var fin float64
var dx float64
var numPart int
var numExec int
var coefs polynomial
// Default values
const defaultAddr string = "localhost:8666"
const defaultTimeout string = "10s"
const defaultLeasing string = "5s"
const defaultStart float64 = 0.0
const defaultFin float64 = 1000000.0
const defaultDx float64 = 0.1
const defaultNumPart int = 200
const defaultnumExec int = 1
// Descriptions
const usageAddr string = "IP:PORT of the Tuple Space"
const usageTimeout string = "Time to wait for messages from the Tuple Space"
const usageLeasing string = "Leasing time of messages sent to the Tuple Space"
const usageStart string = "Approximate the integral value starting from here"
const usageFin string = "Approximate the integral value until here"
const usageDx string = "DeltaX used when calculating the areas in the Riemann sum"
const usageNumPart string = "Ammout of different requests to send to the Tuple Space"
const usagenumExec string = "Number of clients running concurrently (goroutines)"
const usageCoefs string = "Number of sequential executions of the client"
// Set the flag names (long and short for each flag var
func init() {
	flag.StringVar(&addr, "address", defaultAddr, usageAddr)
	flag.StringVar(&addr, "a", defaultAddr, usageAddr)
	
	flag.StringVar(&timeout, "timeout", defaultTimeout, usageTimeout)
	flag.StringVar(&timeout, "t", defaultTimeout, usageTimeout)
	
	flag.StringVar(&leasing, "leasing", defaultLeasing, usageLeasing)
	flag.StringVar(&leasing, "l", defaultLeasing, usageLeasing)
	
	flag.Float64Var(&start, "start", defaultStart, usageStart)
	flag.Float64Var(&start, "s", defaultStart, usageStart)
	
	flag.Float64Var(&fin, "finish", defaultFin, usageFin)
	flag.Float64Var(&fin, "f", defaultFin, usageFin)
	
	flag.Float64Var(&dx, "delta", defaultDx, usageDx)
	flag.Float64Var(&dx, "d", defaultDx, usageDx)
	
	flag.IntVar(&numPart, "partitions", defaultNumPart, usageNumPart)
	flag.IntVar(&numPart, "p", defaultNumPart, usageNumPart)
	
	flag.IntVar(&numExec, "executions", defaultnumExec, usagenumExec)
	flag.IntVar(&numExec, "e", defaultnumExec, usagenumExec)
	
	flag.Var(&coefs, "coefficients", usageCoefs)
	flag.Var(&coefs, "c", usageCoefs)
}

func execute(file *os.File) {
	aux := time.Now()
	r, err := client.NewRiemannStub(addr, timeout, leasing)
	if err != nil {
		log.Fatalln("ERROR creating RiemannStub:", err)
	}
	createTime := time.Since(aux)
	
	aux = time.Now()
	res, errCount, err := r.Integral(start, fin, dx, coefs, numPart)
	if err != nil {
		log.Fatalln("ERROR calculating the Integral:", err)
	}
	executeTime := time.Since(aux)
	
	// Actually write the results (separated by a line)
	str := fmt.Sprintln(createTime.Nanoseconds(), executeTime.Nanoseconds(), errCount)
	c, err := file.WriteString(str)
	if err != nil {
		log.Fatal(fmt.Sprintln("ERROR writing", c, "characters on log file", file.Name(), ":", err))
	}
	
	// Print results and times (on screen)
	fmt.Println("Done! After", errCount, "errors, the final sum is:", res)
	fmt.Println("Creating the Stub (and middleware) took:", createTime)
	fmt.Println("Executing the function took:", executeTime)
	
	done<- true
}

func main() {
	flag.Parse()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	
	// Open a file to store the results
	path := build.Default.GOPATH + "/src/github.com/wisllayvitrio/ppd2014/logs/ppd2014_client_log.txt"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			if err != nil {
				log.Fatal(fmt.Sprintln("ERROR creating file:", err))
			}
		} else {
			log.Fatal(fmt.Sprintln("ERROR opening file:", err))
		}
	}
	defer file.Close()
	
	// Print linte to differentiate between executions
	c, err := file.WriteString("####################################################\n")
	if err != nil {
		log.Fatal(fmt.Sprintln("ERROR writing", c, "characters on log file", file.Name(), ":", err))
	}

	// Call the client executions
	for i := 0; i < numExec; i++ {
		execute(file)
	}
}