package main

import (
	"container/list"
	"sync"
)

type mutexQueue struct{
	mlist	list.List
	mutex	sync.Mutex
}

func newMutexQueue() mutexQueue {
	var mutex sync.Mutex
	var lst list.List
	lst = *lst.Init()
	var mlst mutexQueue
	mlst.mutex = mutex
	mlst.mlist = lst
	return mlst
}

func MQpush(ml mutexQueue, ele interface{}){
	ml.mutex.Lock()
	ml.mlist.PushBack(ele)
	ml.mutex.Unlock()
}

func MQpop(ml mutexQueue) *interface{}{
	var first *interface{}
	first = nil
	ml.mutex.Lock()
	if ml.mlist.Len() > 1 {
		*first = ml.mlist.Remove(ml.mlist.Front())
	}
	ml.mutex.Unlock()
	return first
}

func MQpopAll(mq mutexQueue) *list.List{
	ret := list.New()
	mq.mutex.Lock()
	for n := mq.mlist.Front(); n != nil; n = n.Next(){
		ret.PushBack(n.Value)
	}
	mq.mlist = *mq.mlist.Init()
	mq.mutex.Unlock()
	return ret
}