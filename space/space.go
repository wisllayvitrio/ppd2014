package space

import (
	"fmt"
)

type TupleSpace struct{
	entryList []Tuple
	// TODO: The tuple space has more stuff
}

type Request struct{
	Data Tuple
	Timeout int64
}

func (ts *TupleSpace) Write(tuple Request, dumb *interface{}) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)

	*ok = true
	return nil
}

func (ts *TupleSpace) Read(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Read Called!")
	fmt.Println("Template provided", tuple)
}

func (ts *TupleSpace) Take(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Take Called!")
	fmt.Println("Template provided", tuple)
}

