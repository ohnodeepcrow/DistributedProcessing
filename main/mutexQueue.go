package main

import (
	"container/list"
	"sync"
)

type mutexQueue struct{
	mlist	*list.List
	mutex	sync.Mutex
}

func newMutexQueue() mutexQueue {
	var mutex sync.Mutex
	var lst *list.List
	lst = list.New()
	var mlst mutexQueue
	mlst.mutex = mutex
	mlst.mlist = lst
	return mlst
}

func MQpush(ml mutexQueue, ele interface{}){
	ml.mutex.Lock()
	ml.mlist.PushBack(ele)
	//println("MQpush " + fmt.Sprint(ml.mlist.Front().Value))
	//println("MQpush: " + fmt.Sprint(ml.mlist.Len()))
	ml.mutex.Unlock()
}

func MQpop(ml mutexQueue) interface{}{
	ml.mutex.Lock()
	//println("MQpop " + fmt.Sprint(ml.mlist.Len()))
	if ml.mlist.Len() >= 1 {
		//println("MQPOP " + fmt.Sprint(ml.mlist.Front().Value))
		first := ml.mlist.Remove(ml.mlist.Front())
		//println("MQPOP " + fmt.Sprint(first))
		ml.mutex.Unlock()
		return first
	}
	ml.mutex.Unlock()
	return nil
}

func MQpopAll(mq mutexQueue) *list.List{
	ret := list.New()
	mq.mutex.Lock()
	for n := mq.mlist.Front(); n != nil; n = n.Next(){
		ret.PushBack(n.Value)
	}
	mq.mlist = mq.mlist.Init()
	mq.mutex.Unlock()
	return ret
}