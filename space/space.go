package space

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"github.com/wisllayvitrio/ppd2014/index"
	"github.com/wisllayvitrio/ppd2014/encode"
	"code.google.com/p/go-uuid/uuid"
)

const maxTupleSize int = 10
const tupleIndexDuration time.Duration = 1 * time.Minute
const waitIndexDuration time.Duration = 1 * time.Minute
var nilHashCode string = NilHash()

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

type TupleSpace struct {
	mutex sync.RWMutex		//cotrola o acesso concorrente permitindo apenas leituras simultaneas
	tupleIndex *index.Index //armazena a tupla que esta relacionada a cada indice
	waitList *WaitList 		//armazena a lista de espera por tuplas (linear)
	searchTable [maxTupleSize]*searchIndex
}

type searchIndex struct {
	numAttributes int
	tupleIndex []*index.Index
}

type WaitList struct {
	list []waitState
}

type waitState struct {
	expireTime time.Time
	template Tuple
	wait chan Tuple
	isToTake bool 
}

func NewWaitList() *WaitList {
	waitList := new(WaitList)
	waitList.list = make([]waitState, 0)
	return waitList
}

func NewSearchIndex(numAttributes int) *searchIndex {
	search := &searchIndex{
		numAttributes,
		make([]*index.Index, numAttributes),
	}

	for i := 0; i < numAttributes; i++ {
		search.tupleIndex[i] = index.NewIndex(tupleIndexDuration)
	}

	return search
}

func NewTupleSpace() *TupleSpace {
	space := new(TupleSpace)
	space.tupleIndex = index.NewIndex(tupleIndexDuration)
	space.waitList = NewWaitList()

	for i := 1; i < maxTupleSize; i++ {
		space.searchTable[i] = NewSearchIndex(i)
	}

	return space
}

type Request struct{
	Data Tuple
	Leasing time.Duration
}

func (space *TupleSpace) Write(tuple Request, dummy *Tuple) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)
	space.mutex.Lock()
	defer space.mutex.Unlock()

	searchSpace := space.searchTable[tuple.Data.Size()]
	tupleID := uuid.NewRandom().String()

	//Procurar caras esperando
	//Se tiver alguem verifica o que ele espera

	//Arrumandos os indeices de busca
	for i, arg := range tuple.Data.Args {
		hashcode := encode.Hash(arg)
		searchSpace.tupleIndex[i].Put(hashcode, tupleID, tuple.Leasing)
	}

	//Armazenando no espaço de tuplas
	space.tupleIndex.Put(tupleID, tuple.Data, tuple.Leasing)

	return nil
}

func (space *TupleSpace) Read(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Read Called!")
	fmt.Println("Template provided:", template)

	space.mutex.RLock()
	defer space.mutex.RUnlock()

	searchSpace := space.searchTable[template.Data.Size()]

	var hashList []string = searchHashes(template.Data, searchSpace)

	//Procura por possiveis valores
	for _, hash := range hashList {
		tuples := space.tupleIndex.Get(hash)

		if tuples == nil {
			continue
		}

		for _, value := range tuples {
			if ret, ok := value.(Tuple); ok {
				tuple = &ret
				return nil
			}
		}
	}
	
	//Cria o estado de espera
	waitChan := make(chan Tuple, 1)
	wait := waitState{
		time.Now().Add(template.Leasing),
		template.Data,
		waitChan,
		false,
	}

	//Se nao encontrou espera pelo timeout
	select {
	case *tuple = <-waitChan:
		return nil
	case <-time.After(template.Leasing):
		return errors.New("Data not found")
	}
}

func (space *TupleSpace) Take(template Request, tuple *Tuple) error {
	fmt.Println("TupleSpace.Take Called!")
	fmt.Println("Template provided:", template)
	fmt.Println("Return tuple:", tuple)
	//fmt.Println()
	
	space.mutex.RLock()
	defer space.mutex.RUnlock()

	searchSpace := space.searchTable[template.Data.Size()]

	var hashList []string = searchHashes(template.Data, searchSpace)
	
	//Procura por possiveis valores
	for _, hash := range hashList {
		tuples := space.tupleIndex.Take(hash)

		if tuples == nil {
			continue
		}

		for _, value := range tuples {
			if ret, ok := value.(Tuple); ok {
				tuple = &ret
				return nil
			}
		}
	}
	
	//Cria o estado de espera
	waitChan := make(chan Tuple, 1)
	wait := waitState{
		time.Now().Add(template.Leasing),
		template.Data,
		waitChan,
		true,
	}

	//Se nao encontrou espera pelo timeout
	select {
	case *tuple = <-waitChan:
		return nil
	case <-time.After(template.Leasing):
		return errors.New("Data not found")
	}
}

func toStringList(list []interface{}) []string {
	stringList := make([]string, len(list))

	for i, value := range list {
		if str, ok := value.(string); ok {
			stringList[i] = str
		}
	}

	return stringList
}

func intersectList(list, otherList []string) []string {
	inter := make([]string, 0)
	hashMap := make(map[string]bool)

	for _, hash := range list {
		hashMap[hash] = true
	}

	for _, hash := range otherList {
		if _, exists := hashMap[hash]; exists {
			inter = append(inter, hash)
		}
	}

	return inter
}

func searchHashes(template Tuple, searchSpace *searchIndex) []string {
	hashList := make([]string, 0)

	for i, arg := range template.Args {
		hashcode := encode.Hash(arg)

		if hashcode != nilHashCode {
			list := searchSpace.tupleIndex[i].Get(hashcode)

			if len(hashList) == 0 {
				hashList = toStringList(list)
			} else {
				hashList = intersectList(hashList, toStringList(list))
			}
		}
	}

	return hashList
}