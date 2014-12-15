package space

import (
	"fmt"
	"time"
	"github.com/wisllayvitrio/ppd2014/index"
)

const maxTupleSize int = 10
const tupleIndexDuration time.Duration = 1 * time.Minute
const waitIndexDuration time.Duration = 1 * time.Minute

type Nil struct {
	Dummy bool
}

func NilValue() interface{} {
	return Nil{false}
}

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
	Timeout time.Duration
}

func (ts *TupleSpace) Write(tuple Request, dummy *Tuple) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)
	fmt.Println()
	
	// dummy return
	dummy, err := NewTuple(4, 2)
	if err != nil {
		fmt.Println("dummy err:", err)
	}

	return nil
}

func (ts *TupleSpace) Read(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Read Called!")
	fmt.Println("Template provided:", template)
	fmt.Println("Return tuple:", tuple)
	fmt.Println()

	return nil
}

func (ts *TupleSpace) Take(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Take Called!")
	fmt.Println("Template provided:", template)
	fmt.Println("Return tuple:", tuple)
	//fmt.Println()
	
	// Test return
	res := new(Tuple)
	
	var name string
	var sum int
	args := make([]interface{}, 0)
	 
	template.Data.Get(0, &name)
	template.Data.Get(1, &sum)
	args = append(args, sum)
	
	res, err := NewTuple(name, args)
	if err != nil {
		fmt.Println("SPACE-ERROR:", err)
	}
	
	*tuple = *res
	fmt.Println("Template (human-readable):", name, sum)
	fmt.Println("New tuple:", tuple)
	fmt.Println()
	
	return nil
}
