package space

import (
	"container/list"
	"time"
)

type WaitList struct {
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
	w.waitList.PushBack(state)
}

func (w *WaitList) sendNewTuple(tuple Tuple) bool {
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
