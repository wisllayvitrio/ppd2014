package main

import (
	"fmt"
	"github.com/wisllayvitrio/ppd2014/client"
)

func main() {
	stub := new(client.Stub)
	res, err := stub.Sum(1, 1)
	
	if err != nil {
		fmt.Println("ERROR -", err)
	}
	fmt.Println("Not done yet! sum:", res)
}