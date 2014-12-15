package space

import (
	"github.com/wisllayvitrio/ppd2014/encode"
)

type Nil struct {
	Dummy bool
}

func NilValue() interface{} {
	return Nil{false}
}

func NilHash() string {
	data, _ := encode.EncodeBytes(NilValue())
	return encode.Hash(data)
}