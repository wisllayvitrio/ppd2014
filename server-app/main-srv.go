package main

import (
	"fmt"
	"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/server"
)

func main() {
	calc, err := server.NewCalculator("testServ", "localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	//err = srv.WorkDefault()
	err = calc.Work("1m")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
}