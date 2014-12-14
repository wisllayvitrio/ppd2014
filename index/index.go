package index

import(
	"sync"
)

/*
	Indice que aponta de uma string para uma lista de valores que estao associados a essa string
*/
type Index struct {
	mutex sync.RWMutex
	index map[string][]interface{} 
}

func (i *Index) Put(key string, value interface{}) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if list, exists := i.index[key]; exists {
		list = append(list, value)
	}
	else {
		list = make([]interface{})
		list = append(list, value)
	}
}

func (index *Index) Get(key string) []interface{} {
	i.mutex.RLock()
	defer i.mutex.Unlock()

	if list, exists := i.index[key]; exists {
		ret := make([]interface{})
		copy(ret, list)
		return ret
	} else {
		return nil
	}
}

func (index *Index) Remove(key string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	delete(i.index[key])
}