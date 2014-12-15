package space

import (
	"sync"
	"container/list"
	"time"
)

type WaitList struct {
	mutex sync.RWMutex
	waitList *list.List
}

type WaitState struct {
	expireTime time.Time
	template Tuple
	wait chan Tuple
	isToTake bool
}

func NewWaitState(timeout time.Duration, template Tuple, waitChan chan Tuple, isToTake bool) WaitState {
	return WaitState {
		time.Now().Add(timeout),
		template,
		waitChan,
		isToTake,
	}
}

func NewWaitList() *WaitList {
	w := new(WaitList)
	w.waitList = list.New()
	return w
}

func (w *WaitList) Add(state WaitState) {
	w.mutex.Lock()
	w.waitList.PushBack(state)
	defer w.mutex.Unlock()
}

func (w *WaitList) sendNewTuple(tuple Tuple) bool {
	w.mutex.Lock()
	defer w.mutex.Lock()

	now := time.Now()
	
	for e := w.waitList.Front(); e != nil; {
		if state, ok := e.Value.(WaitState); ok {
			if state.expireTime.Before(now) {
				var aux *list.Element
				aux, e = e, e.Next()
				w.waitList.Remove(aux)
				continue
			}

			if tuple.Match(state.template) {
				state.wait <- tuple
				var aux *list.Element
				aux, e = e, e.Next()
				w.waitList.Remove(aux)
				
				if state.isToTake {
					return true
				}

				continue
			} 

			e = e.Next()
		}
	}

	return false
}
