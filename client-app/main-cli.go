package main

import (
	"fmt"
	"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/client"
)

func main() {
	stub, err := client.NewStub("localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	res, err := stub.Sum(1, 1)
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	fmt.Println("Not done yet! sum:", res)
}