package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
	"github.com/pebbe/zmq4"
	"strings"
	_"context"
)

var counter int
var leader string
func processRequestSend(node NodeInfo, self NodeSocket, input string) {
	m:= decode(input)
		 if m.Type=="Selected" && m.Sender!=m.Receiver{
			ReceiveResult(self,m)
		}
}


func processRequestReceive(nm NodeMap, selfstr string, self NodeSocket, input string) {
	node := nm.Nodes[selfstr]
	m:= decode(input)
	var msg string
	var ms string
	var metric metric
	if m.Type=="Selected"{
		if m.Kind == "Prime" {
			i,_:= strconv.ParseInt(m.Value,10,64)
			num:=big.NewInt(i)
			metric = testPrime(*num)
			ms=metricString(metric)
			nodeinf.PrimeMetric=updateReputation(nodeinf.PrimeMetric, metric, node.NodeName, primeScorer)
			metric.NodeInf = nodeinf
			msg = encode(node.NodeName, m.Sender,m.Kind,ms,m.Job, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,num.String())

		} else if m.Kind == "Hash" {
			metric = crackHash(m.Value)
			ms=hmetricString(metric)
			//fmt.Println(metric)
			nodeinf.HashMetric=updateReputation(nodeinf.HashMetric, metric, node.NodeName, hashScorer)
			metric.NodeInf = nodeinf
			msg = encode(node.NodeName, m.Sender,m.Kind,m.Job,ms, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Value)
		}

		nm.Nodes[selfstr] = nodeinf

		fmt.Println("Hash:")
		fmt.Println(nodeinf.HashMetric)
		fmt.Println(node.HashMetric)
		fmt.Println("Prime:")
		fmt.Println(nodeinf.PrimeMetric)
		fmt.Println(node.PrimeMetric)

		if m.Sender!=m.Receiver{
			SendResult(self,node,decode(msg))
		}
		updatemsg := encode(node.NodeName, "",m.Kind, ms, m.Job,"Update",node.NodeGroup,"","","",metric,"")
		nodeSend(updatemsg, self)
	}
}



/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(selfname string, nm NodeMap, selfsoc NodeSocket, m string){
	node := nm.Nodes[selfname]
	fmt.Print(m+"\n")
	msg:=decode(m)
	var dummy metric

	if msg.Type =="Request" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		LeadNodeSend(m, selfsoc) // group node forwards the request to master node
	} else if msg.Type=="Metric" && (msg.Kind=="Prime"||msg.Kind=="Hash"){
		
		//reply to master node with the best node
		//update busy list
		//master finds the best node
		bestname, bestscore := getBestFreeScore(nm, msg.Kind)
		setBusy(nm, bestname, msg.Job)
		var m metric
		t,exists := nm.Nodes[bestname]
		if exists {
			m.NodeInf = t
		}
		retmsg := encode(node.NodeName, bestname, msg.Kind, msg.Job,strconv.Itoa(bestscore), "Metric", node.NodeGroup,"", node.NodeAddr, node.DataSendPort, m,"")
		LeadNodeSend(retmsg, selfsoc)
	} else if msg.Type=="Selected" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		//update busy list

		kids := getChildren(nm)
		nparr := strings.Split(msg.Input, ";")
		for _,kid := range kids{
			for _,np := range nparr{
				if kid == np && getBusyJob(nm, np) == msg.Job{
					setFree(nm, kid)
				}
			}
		}

		nodeSend(m, selfsoc)
	}else if msg.Type=="Update" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		setFree(nm,msg.Sender)
		snodeinfo := msg.Result.NodeInf
		if msg.Kind == "Prime"{
			snodeinfo.PrimeMetric =  updateReputation(snodeinfo.PrimeMetric, msg.Result, msg.Sender, primeScorer)
		} else if msg.Kind == "Hash"{
			snodeinfo.HashMetric = updateReputation(snodeinfo.HashMetric, msg.Result, msg.Sender, hashScorer)
		}
		oldni := nm.Nodes[msg.Sender]
		oldni.HashMetric = snodeinfo.HashMetric
		oldni.PrimeMetric = snodeinfo.PrimeMetric
		nm.Nodes[msg.Sender] = oldni
		fmt.Println(nm.Nodes["bob"].NodeName)
		fmt.Println(nm.Nodes["bob"].PrimeMetric)
		fmt.Println(nm.Nodes["bob"].HashMetric)
	}else if msg.Type=="Connect" {
		if counter<1{
			counter++
			updateNodeInfo(nm, msg.Sender, msg.Result.NodeInf)
			dummy.NodeInf=nm.Nodes[node.NodeName]
			retmsg := encode(node.NodeName, "", "", "","", "Accepted", "","", "", "", dummy,"")
			selfsoc.datasendsock.Send(retmsg,0)

		}else {
			retmsg := encode(node.NodeName, "", "", "","", "Rejected", "","", "", "", dummy,"")
			selfsoc.datasendsock.Send(retmsg,0)
		}

	}else if msg.Type=="Hi" {
			for k,v := range nm.Nodes {
				if v.Leader == false && v.Master == false{
					var dummy metric
					dummy.NodeInf = v
					up := encode(k, "", "", "", "", "UpdateUptime", "", "", "", "", dummy, "")
					nodeSend(up,selfsoc)
				}
			}
				up := encode("", "", "", "", "", "End", "", "", "", "", dummy, "")
				nodeSend(up,selfsoc)

	}else if msg.Type=="Bye"{
		clearNodeInfo(nm, msg.Sender)
		var dummy metric
		var dummyni NodeInfo
		dummyni.Uptime = node.Uptime
		dummy.NodeInf = dummyni
		retmsg := encode(node.NodeName, "", "", "","", "UpdateUptime", "","", "", "", dummy,"")
		nodeSend(retmsg, selfsoc)

	}

}

/*Master node will call this function after it received a message from a node. It will use send to retransmit the node. */
func MasterNodeRec(node NodeInfo,nm NodeMap,self NodeSocket, m string){

	fmt.Print(m+"\n")
	msg := decode(m)
	var dummy metric
	if msg.Type=="Request" {
		message := encode(msg.Sender, msg.Receiver, msg.Kind, msg.Job, msg.Value, "Metric", msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		nodeSend(message, self)
		bestnode := MasterNodeMet(nm, node, self,message)
		//TODO: Fill out the reply message in a useful way
		//TODO: WE NEED A SPECIAL MESSAGE TYPE FOR THIS... CAN'T USE METRIC

		//put the best node in msg.Receiver
		m := encode(msg.Sender, bestnode, msg.Kind, msg.Job, msg.Value, msg.Type, msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		LeadNodeSend(m, self)
	} else if msg.Type=="Hi" {
		updateNodeInfo(nm, msg.Sender, msg.Result.NodeInf)
		for k,v := range nm.Nodes{
			if v.Leader==true && v.Master == false{
				var dummy metric
				dummy.NodeInf = v
				up := encode(k, "", "", "","", "UpdateUptime", "","", "", "", dummy,"")
				LeadNodeSend(up,self)
			}
		}

	}else if msg.Type=="Bye"{
		clearNodeInfo(nm, msg.Sender)
		var dummy metric
		var dummyni NodeInfo
		dummyni.Uptime = node.Uptime
		dummy.NodeInf = dummyni
		retmsg := encode(node.NodeName, "", "", "","", "UpdateUptime", "","", "", "", dummy,"")
		LeadNodeSend(retmsg, self)

	}else if msg.Type=="Boot"{

		for leader,nodeinfo := range nm.Nodes{

			if (leader != node.NodeName){
				msg := encode(leader,"","","","","Leader","","",nodeinfo.NodeAddr,nodeinfo.DataSendPort,dummy,"")
				self.datasendsock.Send(msg,0)
				//Req/Rep needs to send/receive
				//This means we need 1 send and 1 recv before the next send
				self.datasendsock.Recv(0)
			}
		}
		endmsg := encode(leader,"","","","","Leader","","","","",dummy,"")
		self.datasendsock.Send(endmsg,0)
	}
}

func MasterNodeMet(nm NodeMap, node NodeInfo,self NodeSocket, msg string) string {
	counter := make(map[string]bool)
	size := len(getChildren(nm)) - 1
	c := 0
	maxRep := -1
	bestNode := ""
	for _,child := range getChildren(nm){
		counter[child] = false
	}
	for {
		s := MQpop(self.recvq)

		if s != nil {
			message := fmt.Sprint(s)
			m := decode(message)
			req := decode(msg)
			if m.Type == "Metric" {
				if (m.Job == req.Job){
					counter[m.Sender]=true
					a,_ := strconv.Atoi(m.Value)
					if (a > maxRep){
						bestNode = m.Receiver
						maxRep = a
					}
					c = c + 1
					if (c==size){
						return bestNode
					}


				}
				//check if we received all metrics from GLs and set bestnode accordingly
			} else if m.Type == "Request" {
				MQpush(self.recvq, s)
			}
		}
	}
	return ""
}


func MessageHandler(selfname string, nm NodeMap, selfsoc NodeSocket){

	selfnode := nm.Nodes[selfname]

	s := MQpop(selfsoc.recvq)
	if s == nil{
		return
	}
	message := fmt.Sprint(s)
	m:= decode(message)
	if m.Receiver == selfname {
		processRequestReceive(nm, selfname, selfsoc, message)
	}else if  m.Sender == selfname {
		processRequestSend(selfnode, selfsoc, message)

	} else if selfsoc.leader == true {
			//println("Retransmitting " + message)
			LeadNodeRec(selfname, nm, selfsoc, message)
	} else if selfsoc.master == true {
		MasterNodeRec(selfnode,nm,selfsoc,message)
	}else if m.Type == "UpdateUptime"{
		updateNodeInfo(nm, m.Sender, m.Result.NodeInf)
	}else {
		return //drop message
	}
}

func startMessageHandler(selfname string, nm NodeMap, selfsoc NodeSocket){
	for{

		time.Sleep(time.Millisecond*50)
		MessageHandler(selfname, nm, selfsoc)
	}
}

func ReceiveResult(self NodeSocket,msg Message){
	addr:=msg.Address
	port:=msg.Port
	soc :=establishClient(addr,port)
	soc.Send("req",0)
	for {
		tmp,_ := soc.Recv(zmq4.DONTWAIT)
		if (tmp != ""){
			//send the result back through data socket.
			fmt.Print("Received from data soc: "+tmp + "\n")
			MQpush(self.dataq, tmp)
			return
		}
		time.Sleep(time.Millisecond*50)
	}



}
func SendResult(self NodeSocket, node NodeInfo, m Message){
	addr:=node.NodeAddr
	port:=node.DataSendPort
	send_sock:= establishServer(addr,port)
	t,_ :=(send_sock.GetIdentity())
	fmt.Print(t +"\n")
	i,_:= strconv.ParseInt(m.Value,10,64)
	num:=big.NewInt(i)
	metric := testPrime(*num)
	ms:=metricString(metric)
	msg := encode(node.NodeName, m.Sender,m.Kind,m.Job,ms, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Value)
	fmt.Print(string(msg)+ "\n")
	for {
		signal,_ := send_sock.Recv(zmq4.DONTWAIT)
		if (signal != ""){
			//send the result back through data socket.
			fmt.Print(signal+"\n")
			send_sock.Send(msg,0)
			return
		}
		time.Sleep(time.Millisecond*50)
	}
}

