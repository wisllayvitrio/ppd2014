package space 

import (
	"fmt"
	"errors"
	"github.com/wisllayvitrio/ppd2014/encode"
)

type Tuple struct {
	Args [][]byte
}

func NewTuple(args ...interface{}) (*Tuple, error) {
	tuple := new(Tuple)
	for _, arg := range args {
		data, err := encode.EncodeBytes(arg)
		
		if err != nil {
			return nil, errors.New("Error while encoding arguments")
		}
		
		tuple.Args = append(tuple.Args, data)
	}
	
	return tuple, nil
}

func (t* Tuple) Get(i int, ptr interface{}) error {
	err := encode.DecodeBytes(t.Args[i], ptr)
	if err != nil {
		return errors.New("Error while decoding tuple argument")
	}

	return nil
}

func (t* Tuple) Size() int {
	return len(t.Args)
}