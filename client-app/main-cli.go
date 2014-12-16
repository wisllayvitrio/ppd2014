package main

import (
	"log"
	"fmt"
	"github.com/wisllayvitrio/ppd2014/client"
)

func main() {
	/*stub, err := client.NewStub("localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	for i := 0; ; i++ {
		<-time.After(10 * time.Millisecond)
	
		res, err := stub.Sum(i, i + 2)
		if err != nil {
			fmt.Println("ERROR -", err)
			debug.PrintStack()
		}
		fmt.Println("Done! sum:", res, "Execution number:", i)
	}*/
	
	r, err := client.NewRiemannStub("localhost:8666", "10s", "5s")
	if err != nil {
		log.Fatal("ERROR creating RiemannStub: %s\n", err)
	}
	
	coefs := []float64{4, 2}
	res, errCount, err := r.Integral(0, 1000000, 0.1, coefs, 200)
	if err != nil {
		log.Fatal("ERROR calculating the Integral: %s\n", err)
	}
	fmt.Println("Done! After", errCount, "errors, the final sum is:", res)
}