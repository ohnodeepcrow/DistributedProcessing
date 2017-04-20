package main

import (
	"github.com/pebbe/zmq4"
	"syscall"
	"time"
	"fmt"
)

type NodeSocket struct {
	appq	 mutexQueue
	recvq    mutexQueue
	sendsock *zmq4.Socket
	recvsock *zmq4.Socket
	leadersendsock *zmq4.Socket
	leaderrecvsock *zmq4.Socket
	leader   bool
}

func establishLeader(context *zmq4.Context, self NodeInfo, peer NodeInfo) NodeSocket{
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

	lssoc, err := context.NewSocket(zmq4.PUB)
	check(err)
	lrsoc, err := context.NewSocket(zmq4.SUB)
	check(err)

	lsocstr := "tcp://" + self.NodeAddr + ":" + self.LeaderSendPort
	fmt.Print("self.LeaderSendPort "+self.LeaderSendPort +"\n")
	err = lssoc.Bind(lsocstr)
	check(err)
	/*
	lsocstr = "tcp://" + self.NodeAddr + ":" + self.LeaderRecvPort
	fmt.Print("self.LeaderRecvPort "+self.LeaderRecvPort +"\n")
	err = lrsoc.Bind(lsocstr)
	check(err)

	lsocstr = "tcp://" + peer.NodeAddr + ":" + peer.LeaderRecvPort
	fmt.Print("peer.LeaderRecvPort "+ peer.LeaderRecvPort +"\n")
	err = lssoc.Connect(lsocstr)
	check(err)
	*/
	lrsoc.SetSubscribe("")
	lrsocstr := "tcp://" + peer.NodeAddr + ":" + peer.LeaderSendPort
	fmt.Print("peer.LeaderSendPort "+peer.LeaderSendPort +"\n")
	err = lrsoc.Connect(lrsocstr)
	check(err)



	fmt.Print("")
	var ret NodeSocket
	ret.leader = true
	ret.sendsock = ssoc
	ret.recvsock = rsoc

	ret.leaderrecvsock = lrsoc
	ret.leadersendsock = lssoc

	ret.recvq = newMutexQueue()
	ret.appq = newMutexQueue()
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
	ret.appq = newMutexQueue()
	return ret
}

func nodeSend(str string, soc NodeSocket) error{
	_, err := soc.sendsock.Send(str, 0)
	check(err)
	return err

}
func LeadNodeSend(str string, soc NodeSocket) error{
	_, err := soc.leadersendsock.Send(str, 0)
	check(err)
	return err
}

func nodeReceive(soc NodeSocket){
	for {
		tmp,err := soc.recvsock.Recv(zmq4.DONTWAIT)
		if soc.leader == true{
			tmp1,err := soc.leaderrecvsock.Recv(zmq4.DONTWAIT)

			if tmp1 != "" {
				fmt.Print("from leader receive sock : "+(fmt.Sprint(tmp1))+"\n")
				MQpush(soc.recvq, tmp1)
			}
			if err == syscall.EAGAIN {
				continue
			}
		}
		if err == syscall.EAGAIN {
			continue
		}
		if tmp != "" {
			MQpush(soc.recvq, tmp)
		}
		time.Sleep(time.Millisecond*50)
	}
}

func startReceiver(soc NodeSocket){
	nodeReceive(soc)

}