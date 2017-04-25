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
	dataq 	 mutexQueue // only store the computation result
	sendsock *zmq4.Socket //leaf multicast send
	recvsock *zmq4.Socket //leaf multicast receive
	leadersendsock *zmq4.Socket //leader multicast send
	leaderrecvsock *zmq4.Socket //leader multicast receive
	datasendsock *zmq4.Socket //p2p data send
	datarecvsock *zmq4.Socket //p2p receive
	leader   bool //Am I a group leader?
	master bool //Am I the root node?
}

func establishLeader(context *zmq4.Context, self NodeInfo, master NodeInfo) NodeSocket{
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

	lrsoc, err := context.NewSocket(zmq4.SUB)
	check(err)
	lssoc, err := context.NewSocket(zmq4.PUSH)
	check(err)
	socstr = "tcp://" + master.NodeAddr + ":" + master.RecvPort
	err = lssoc.Connect(socstr)
	check(err)
	socstr = "tcp://" + master.NodeAddr + ":" + master.SendPort
	lrsoc.SetSubscribe("")
	err = lrsoc.Connect(socstr)
	check(err)


	fmt.Print("")
	var ret NodeSocket
	ret.leader = true
	ret.master = false
	ret.sendsock = ssoc
	ret.recvsock = rsoc

	ret.leaderrecvsock = lrsoc
	ret.leadersendsock = lssoc

	ret.recvq = newMutexQueue()
	ret.appq = newMutexQueue()
	ret.dataq = newMutexQueue()
	return ret
}
func establishMaster (context *zmq4.Context, self NodeInfo) NodeSocket{
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
	ret.leader = false
	ret.master = true
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	ret.appq = newMutexQueue()
	ret.dataq = newMutexQueue()
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
	ret.master = false
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	ret.appq = newMutexQueue()
	ret.dataq = newMutexQueue()
	return ret
}

func establishClient(addr string, port string, socket NodeSocket) *zmq4.Socket{
	context,_ := zmq4.NewContext()
	soc,_ := context.NewSocket(zmq4.REQ)
	socstr := "tcp://" + addr + ":" + port
	soc.Connect(socstr)
	return soc
}
func establishServer(addr string, port string, socket NodeSocket)*zmq4.Socket{
	context,_ := zmq4.NewContext()
	soc,_ := context.NewSocket(zmq4.REP)
	socstr := "tcp://" + addr + ":" + port
	soc.Bind(socstr)
	return soc
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