package space

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"github.com/wisllayvitrio/ppd2014/index"
	"github.com/wisllayvitrio/ppd2014/encode"
	"github.com/wisllayvitrio/ppd2014/logger"
	"code.google.com/p/go-uuid/uuid"
)

const maxTupleSize int = 10
const tupleIndexDuration time.Duration = 1 * time.Minute
const waitIndexDuration time.Duration = 1 * time.Minute
var nilHashCode string = NilHash()

type TupleSpace struct {
	mutex sync.RWMutex		//cotrola o acesso concorrente permitindo apenas leituras simultaneas
	tupleIndex *index.Index //armazena a tupla que esta relacionada a cada indice
	waitList *WaitList 		//armazena a lista de espera por tuplas (linear)
	searchTable [maxTupleSize]*searchIndex
	l *logger.Logger
}

func NewTupleSpace() *TupleSpace {
	space := new(TupleSpace)
	space.tupleIndex = index.NewIndex(tupleIndexDuration)
	space.waitList = NewWaitList()

	for i := 1; i < maxTupleSize; i++ {
		space.searchTable[i] = NewSearchIndex(i)
	}

	var err error
	space.l, err = logger.NewLogger("./ppd2014_space_log.txt", time.Second)
	if err != nil {
		fmt.Println("ERROR creating logger:", err)
	}
	go space.l.LogStart()
	
	return space
}

type Request struct{
	Data Tuple
	Leasing time.Duration
}

func (space *TupleSpace) Write(tuple Request, dummy *Tuple) error {
	aux := time.Now()
	
	space.mutex.Lock()
	defer space.mutex.Unlock()

	searchSpace := space.searchTable[tuple.Data.Size()]

	//Procurar caras esperando
	//Se tiver alguem verifica o que ele espera
	wasTaken := space.waitList.sendNewTuple(tuple.Data)

	if wasTaken {
		return nil
	}

	//Arrumandos os indeices de busca
	tupleID := uuid.NewRandom().String()

	for i, arg := range tuple.Data.Args {
		hashcode := encode.Hash(arg)
		searchSpace.tupleIndex[i].Put(hashcode, tupleID, tuple.Leasing)
	}
	//Armazenando no espaço de tuplas
	space.tupleIndex.Put(tupleID, tuple.Data, tuple.Leasing)

	space.l.AddTime(false, time.Since(aux))

	return nil
}

func (space *TupleSpace) Read(template Request, tuple *Tuple) error {
	aux := time.Now()

	searchSpace := space.searchTable[template.Data.Size()]

	ret := space.searchTuple(template.Data, searchSpace, true)
	
	space.l.AddTime(true, time.Since(aux))

	if ret != nil {
		*tuple = *ret
		return nil
	}

	//Cria o estado de espera
	waitState := NewWaitState(template.Leasing, template.Data, make(chan Tuple, 1), false)
	space.waitList.Add(waitState)

	//Se nao encontrou espera pelo timeout
	select {
	case *tuple = <-waitState.wait:
		return nil
	case <-time.After(template.Leasing):
		return errors.New("Data not found")
	}
}

func (space *TupleSpace) Take(template Request, tuple *Tuple) error {
	aux := time.Now()

	searchSpace := space.searchTable[template.Data.Size()]
	
	ret := space.searchTuple(template.Data, searchSpace, true)
	
	space.l.AddTime(true, time.Since(aux))

	if ret != nil {
		*tuple = *ret
		return nil
	}

	//Cria o estado de espera
	waitState := NewWaitState(template.Leasing, template.Data, make(chan Tuple, 1), true)
	space.waitList.Add(waitState)

	//Se nao encontrou espera pelo timeout
	select {
	case *tuple = <-waitState.wait:
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
			list := searchSpace.tupleIndex[i].Get(hashcode, false)

			if len(hashList) == 0 {
				hashList = toStringList(list)
			} else {
				hashList = intersectList(hashList, toStringList(list))
			}
		}
	}

	return hashList
}

func (space *TupleSpace) searchTuple(template Tuple, searchSpace *searchIndex, remove bool) *Tuple {
	space.mutex.RLock()
	defer space.mutex.RUnlock()

	var hashList []string = searchHashes(template, searchSpace)
	
	//Procura por possiveis valores
	for _, hash := range hashList {
		tuples := space.tupleIndex.Get(hash, remove)

		if tuples == nil {
			continue
		}

		for _, value := range tuples {
			if ret, ok := value.(Tuple); ok {
				return &ret
			}
		}
	}

	return nil
}