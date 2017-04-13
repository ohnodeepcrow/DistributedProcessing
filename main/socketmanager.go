package main

import (
	"github.com/pebbe/zmq4"
	"time"
)

type NodeSocket struct {
	sendsock		*zmq4.Socket
	recvsock		*zmq4.Socket
	leader			bool
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
	return ret
}



func nodeSend(str string, soc *zmq4.Socket){
	out := str + " | " + time.Now().Format("15:04:05")
	soc.Send(out, 0)
}

func nodeReceive(soc *zmq4.Socket){
	oldtmp := ""
	for i := 0; i < 1000000; i++{
		tmp,_ := soc.Recv(zmq4.DONTWAIT)
		if tmp != oldtmp {
			println("Received: " + tmp)
			oldtmp = tmp
		}

	}
}