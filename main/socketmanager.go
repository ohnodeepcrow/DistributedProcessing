package main

import (
	"github.com/pebbe/zmq4"
	"syscall"
)

type NodeSocket struct {
	recvq    mutexQueue
	sendsock *zmq4.Socket
	recvsock *zmq4.Socket
	leader   bool
}

func establishLeader(context *zmq4.Context, self NodeInfo) NodeSocket{
	ssoc, err := context.NewSocket(zmq4.PUB)
	check(err)
	rsoc, err := context.NewSocket(zmq4.PULL)
	check(err)
	socstr := "tcp://" + self.NodeAddr + ":" + self.SendPort
	err = ssoc.Bind(socstr)
	check(err)
	socstr = "tcp://" + self.NodeAddr + ":" + self.RecvPort
	err = rsoc.Bind(socstr)
	check(err)

	var ret NodeSocket
	ret.leader = true
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	return ret
}

func establishMember(context *zmq4.Context, self NodeInfo, ldr NodeInfo) NodeSocket{
	rsoc, err := context.NewSocket(zmq4.SUB)
	check(err)
	ssoc, err := context.NewSocket(zmq4.PUSH)
	check(err)
	socstr := "tcp://" + ldr.NodeAddr + ":" + ldr.RecvPort
	err = ssoc.Connect(socstr)
	check(err)
	socstr = "tcp://" + ldr.NodeAddr + ":" + ldr.SendPort
	rsoc.SetSubscribe("")
	err = rsoc.Connect(socstr)
	check(err)

	var ret NodeSocket
	ret.leader = false
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	return ret
}

func nodeSend(str string, soc NodeSocket) error{
	_, err := soc.sendsock.Send(str, 0)
	return err

}

func nodeReceive(soc NodeSocket){
	for {
		tmp,err := soc.recvsock.Recv(zmq4.DONTWAIT)
		if err == syscall.EAGAIN {
			continue
		}
		if tmp != "" {
			MQpush(soc.recvq, tmp)
		}
	}
}

func startReceiver(soc NodeSocket){
	nodeReceive(soc)
}