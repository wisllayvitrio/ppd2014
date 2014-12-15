package main

import (
	"fmt"
	"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/server"
)

func main() {
	srv, err := server.NewService("testServ", "localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	//err = srv.WorkDefault()
	err = srv.Work("1m")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
}