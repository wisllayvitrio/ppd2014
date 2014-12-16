package main

import (
	"log"
	//"fmt"
	//"runtime/debug"
	"github.com/wisllayvitrio/ppd2014/server"
)

func main() {
	/*calc, err := server.NewCalculator("testServ", "localhost:8666")
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}
	
	//err = srv.WorkDefault()
	err = calc.Work()
	if err != nil {
		fmt.Println("ERROR -", err)
		debug.PrintStack()
	}*/
	
	r, err := server.NewRiemannCalculator("riemann", "localhost:8666", "10s", "5s", 10)
	if err != nil {
		log.Fatal("ERROR creating RiemannCalculator: %s\n", err)
	}
	
	// Keep working forever
	err = r.Work()
	if err != nil {
		log.Fatal("ERROR working: %s\n", err)
	}
}