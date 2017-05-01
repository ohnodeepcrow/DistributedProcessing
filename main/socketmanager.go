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
	sendq	 mutexQueue
	lsendq   mutexQueue
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
	counter=0
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

	ds:=establishServer(self.NodeAddr, self.DataSendPort)


	fmt.Print("")
	var ret NodeSocket
	ret.leader = true
	ret.datasendsock=ds
	ret.master = false
	ret.sendsock = ssoc
	ret.recvsock = rsoc

	ret.datasendsock=ds

	ret.leaderrecvsock = lrsoc
	ret.leadersendsock = lssoc

	ret.recvq = newMutexQueue()
	ret.sendq = newMutexQueue()
	ret.lsendq = newMutexQueue()
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
	counter=0
	ds:=establishServer(self.NodeAddr, self.DataSendPort)
	var ret NodeSocket
	ret.datasendsock=ds
	ret.leader = false
	ret.master = true
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	ret.sendq = newMutexQueue()
	ret.lsendq = newMutexQueue()
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
	ds:=establishServer(self.NodeAddr, self.DataSendPort)
	ret.datasendsock=ds
	ret.leader = false
	ret.master = false
	ret.sendsock = ssoc
	ret.recvsock = rsoc
	ret.recvq = newMutexQueue()
	ret.sendq = newMutexQueue()
	ret.lsendq = newMutexQueue()
	ret.appq = newMutexQueue()
	ret.dataq = newMutexQueue()
	return ret
}

func establishClient(addr string, port string) *zmq4.Socket{
	context,_ := zmq4.NewContext()
	soc,_ := context.NewSocket(zmq4.REQ)
	socstr := "tcp://" + addr + ":" + port
	soc.Connect(socstr)
	return soc
}
func establishServer(addr string, port string)*zmq4.Socket{
	context,_ := zmq4.NewContext()
	soc,_ := context.NewSocket(zmq4.REP)
	socstr := "tcp://" + addr + ":" + port
	soc.Bind(socstr)
	return soc
}
func nodeSend(str string, soc NodeSocket){
	MQpush(soc.sendq,str)
	//_, err := soc.sendsock.Send(str, 0)
	//check(err)
	return

}
func LeadNodeSend(str string, soc NodeSocket) {
	//_, err := soc.leadersendsock.Send(str, 0)
	//check(err)
	MQpush(soc.lsendq,str)
	return
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
		tmp2,err := soc.datasendsock.Recv(zmq4.DONTWAIT)
		if tmp2 != "" {
			MQpush(soc.recvq, tmp2)
		}
		if err == syscall.EAGAIN {
			continue
		}
		time.Sleep(time.Millisecond*50)
	}
}
func startSender (soc NodeSocket){
	for {
		if  soc.leader == true{

			s := MQpop(soc.lsendq)
			if s != nil {
				msg := fmt.Sprint(s)
				_, err := soc.leadersendsock.Send(msg, 0)
				check(err)
			}
		}
		t := MQpop(soc.sendq)
		if t != nil{
			m := fmt.Sprint(t)
			_, err := soc.sendsock.Send(m, 0)
			check(err)
		}
		time.Sleep(time.Millisecond*50)
	}
}
func startReceiver(soc NodeSocket){
	nodeReceive(soc)
}

func BootStrap(context *zmq4.Context, self NodeInfo, master NodeInfo, nm NodeMap) NodeSocket{
	MasterAddr := master.NodeAddr
	MasterPort := master.DataSendPort
	var dummy metric
	dummy.NodeInf=self
	m1 := encode(self.NodeName, "", "",getCurrentTimestamp(),"","Boot","","","","",dummy,"")
	soc,_ := context.NewSocket(zmq4.REQ)
	socstr := "tcp://" + MasterAddr + ":" + MasterPort
	soc.Connect(socstr)
	soc.Send(m1,0)
	for {
		tmp,_ := soc.Recv(zmq4.DONTWAIT)
		if (tmp != ""){
			fmt.Print("Received from data soc: "+tmp + "\n")

			msg:=decode(tmp)
			if msg.Type=="Leader" {

				dummy.NodeInf=self
				//say Hi

				if msg.Address =="" && msg.Port==""{
					ns := establishLeader(context,self,master)
					m := encode(self.NodeName, "", "",getCurrentTimestamp(),"","Hi","","","","",dummy,"")
					LeadNodeSend(m,ns)
					return ns

				} else if msg.Address!=""{
					LeaderAddr := msg.Address
					LeaderPort := msg.Port

					soc,_ := context.NewSocket(zmq4.REQ)
					socstr := "tcp://" + LeaderAddr + ":" + LeaderPort
					soc.Connect(socstr)
					m := encode(self.NodeName, "", "",getCurrentTimestamp(),"","Connect","","","","",dummy,"")
					soc.Send(m,0)
					for{
						temp,_ := soc.Recv(zmq4.DONTWAIT)
						m:=decode(temp)
						if m.Type=="Accepted"{
							ns := establishMember(context, self, m.Result.NodeInf)
							m := encode(self.NodeName, "", "",getCurrentTimestamp(),"","Hi","","","","",dummy,"")
							nodeSend(m,ns)
							return ns

						}else if m.Type=="Rejected"{

						}
					}

				}
		}
		time.Sleep(time.Millisecond*50)
	}


}}