package index

import(
	"sync"
	"time"
)

/*
	Indice que aponta de uma string para uma lista de valores que estao associados a essa string
*/
type Index struct {
	cleanInterval time.Duration
	mutex sync.RWMutex
	index map[string][]entry
}

type entry struct {
	data interface{}
	expireTime time.Time
}

func NewIndex(cleanInterval time.Duration) *Index{
	i := new(Index)
	i.cleanInterval = cleanInterval
	i.index = make(map[string][]entry)
	go i.clean()
	return i
}

func (i *Index) clean() {
	for {
		<- time.After(i.cleanInterval)
		i.mutex.Lock()
		defer i.mutex.Unlock()

		time := time.Now()
		for key, list := range i.index {
			newList := make([]entry, 0)

			for _, value := range list {
				if value.expireTime.After(time) {
					newList = append(newList, value)
				}
			}

			if len(newList) == 0 {
				delete(i.index, key)
			} else {
				i.index[key] = newList
			}
		}

		i.mutex.Unlock()
	}
}

func (i *Index) Put(key string, value interface{}, duration time.Duration) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	expireTime := time.Now().Add(duration)

	var newList []entry

	if list, exists := i.index[key]; exists {
		newList = append(list, entry{value, expireTime})
	} else {
		newList = []entry{entry{value, expireTime}}
	}

	i.index[key] = newList
}

func (i *Index) Get(key string) []interface {} {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	time := time.Now()

	if list, exists := i.index[key]; exists {
		ret := make([]interface{}, 0)

		for _, value := range list {
			if value.expireTime.Before(time) { 
				ret = append(ret, value.data)
			}
		}

		return ret
	} else {
		return nil
	}
}

func (i *Index) Take(key string) []interface{} {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	time := time.Now()

	if list, exists := i.index[key]; exists {
		delete(i.index, key)
		
		ret := make([]interface{}, 0)

		for _, value := range list {
			if value.expireTime.Before(time) { 
				ret = append(ret, value.data)
			}
		}

		return ret
	} else {
		return nil
	}
}

func (i *Index) RemoveAll(key string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	delete(i.index, key)
}