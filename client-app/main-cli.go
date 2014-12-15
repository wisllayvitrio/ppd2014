package main

import (
	"fmt"
	"time"
	"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/client"
)

func main() {
	stub, err := client.NewStub("localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	for i := 0; ; i++ {
		<-time.After(100 * time.Millisecond)
	
		res, err := stub.Sum(1, 1)
		if err != nil {
			fmt.Println("ERROR -", err)
			debug.PrintStack()
		}
		fmt.Println("Done! sum:", res, "Execution number:", i)
	}
}