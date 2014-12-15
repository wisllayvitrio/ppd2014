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
			return nil, errors.New(fmt.Sprintln("Error while encoding arguments -", err))
		}
		
		tuple.Args = append(tuple.Args, data)
	}
	
	return tuple, nil
}

func (t* Tuple) Get(i int, ptr interface{}) error {
	err := encode.DecodeBytes(t.Args[i], ptr)
	if err != nil {
		return errors.New(fmt.Sprintln("Error while decoding tuple argument -", err))
	}
	
	return nil
}

func (t* Tuple) Size() int {
	return len(t.Args)
}

func (t *Tuple) Match(template Tuple) bool {
	if t.Size() != template.Size() {
		return false
	}

	for i, arg1 := range t.Args {
		arg2 := template.Args[i]
		hash1, hash2 := encode.Hash(arg1), encode.Hash(arg2)

		if hash2 != NilHash() && hash1 != hash2 {
			return false
		}
	}

	return true
}