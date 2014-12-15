package space

import (
	"fmt"
	"time"
	"github.com/wisllayvitrio/ppd2014/index"
)

const maxTupleSize int = 10
const tupleIndexDuration time.Duration = 1 * time.Minute
const waitIndexDuration time.Duration = 1 * time.Minute

type TupleSpace struct{
	tupleIndex index.Index //armazena a tupla que esta relacionada a cada indice
	searchTable [maxTupleSize]searchIndex
}

type searchIndex struct {
	numAttributes int
	tupleIndex []*index.Index
	waitIndex []*index.Index
}

func makeSearchIndex(numAttributes int) searchIndex {
	search := searchIndex{
		numAttributes,
		make([]*index.Index, numAttributes),
		make([]*index.Index, numAttributes),
	}

	for i := 0; i < numAttributes; i++ {
		search.tupleIndex[i] = index.NewIndex(tupleIndexDuration)
		search.waitIndex[i] = index.NewIndex(waitIndexDuration)
	}

	return search
}

type Request struct{
	Data Tuple
	Timeout int64
}

func (ts *TupleSpace) Write(tuple Request, dumb *interface{}) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)

	return nil
}

func (ts *TupleSpace) Read(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Read Called!")
	fmt.Println("Template provided", tuple)

	return nil
}

func (ts *TupleSpace) Take(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Take Called!")
	fmt.Println("Template provided", tuple)

	return nil
}
