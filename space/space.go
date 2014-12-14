package space

import (
	"fmt"
)

type TupleSpace struct{
	entryList []Tuple
	// TODO: The tuple space has more stuff
}

func (ts *TupleSpace) Write(tuple Tuple, ok *bool) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)

	*ok = true
	return nil
}

