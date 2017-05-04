package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
	"strings"
	_"context"
)

var counter int
var leader string
func processRequestSend(node NodeInfo, self NodeSocket, input string) {
	fmt.Println("DEBUG: " + input)
	m := decode(input)
	MQpush(self.dataq, m)
}


func processRequestReceive(nm NodeMap, selfstr string, self NodeSocket, input string) {
	node := nm.Nodes[selfstr]
	m:= decode(input)
	var msg string
	var ms string
	var metric metric
	if m.Type=="Selected" || m.Type == "Train"{
		//TODO: Actually process request/response portion --> need to make node 2 node connection
		if m.Kind == "Prime" {
			i,_:= strconv.ParseInt(m.Value,10,64)
			num:=big.NewInt(i)
			metric = testPrime(*num)
			ms=metricString(metric)
			nodeinf.PrimeMetric=updateReputation(nodeinf.PrimeMetric, metric, node.NodeName, primeScorer)
			metric.NodeInf = nodeinf
			msg = encode(node.NodeName, m.Sender,m.Kind,ms,m.Job, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Input)

		} else if m.Kind == "Hash" {
			metric = crackHash(m.Value)
			ms=hmetricString(metric)
			//fmt.Println(metric)
			nodeinf.HashMetric=updateReputation(nodeinf.HashMetric, metric, node.NodeName, hashScorer)
			metric.NodeInf = nodeinf
			msg = encode(node.NodeName, m.Sender,m.Kind,m.Job,ms, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Value)
		}

		nm.Nodes[selfstr] = nodeinf

		//fmt.Println("Hash:")
		//fmt.Println(nodeinf.HashMetric)
		//fmt.Println(node.HashMetric)
		//fmt.Println("Prime:")
		//fmt.Println(nodeinf.PrimeMetric)
		//fmt.Println(node.PrimeMetric)

		if m.Sender!=m.Receiver{
			SendResult(self,m.Result.NodeInf,decode(msg))
		} else if m.Type == "Selected"{
			MQpush(self.dataq, decode(msg))
		}
		updatemsg := encode(node.NodeName, m.Result.NodeInf.NodeName,m.Kind, ms, m.Job,"Update",node.NodeGroup,"","","",metric,"")
		nodeSend(updatemsg, self)
	}
}



/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(selfname string, nm NodeMap, selfsoc NodeSocket, m string){
	node := nm.Nodes[selfname]
	fmt.Print(m+"\n")
	msg:=decode(m)
	var dummy metric
	var r map[string]int
	var r1 map[string]int

	if (msg.Type =="Request" || msg.Type=="Board") && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		LeadNodeSend(m, selfsoc) // group node forwards the request to master node
	}else if (msg.Type=="BoardReply")  {
		nodeSend(m, selfsoc) // group node forwards the request to master node
	}else if(msg.Type=="BoardLeader"  )  {
		fmt.Print("inside leader")
		for _,child := range getChildren(nm){
			fmt.Print(child)
			r[child]=nm.Nodes[child].PrimeMetric.Score
			r1[child]=nm.Nodes[child].HashMetric.Score
		}
		s:=encodeRep(r)
		s1:=encodeRep(r1)
		var m metric
		retmsg := encode(node.NodeName,"", msg.Kind, msg.Job,s, "BoardRequest", s1,"", "", node.DataSendPort, m,"")
		fmt.Print("** ")
		fmt.Print(retmsg)
		fmt.Print("\n")
		LeadNodeSend(retmsg, selfsoc)
	}else if msg.Type=="Metric" && (msg.Kind=="Prime"||msg.Kind=="Hash"){

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
	} else if (msg.Type=="Selected" || msg.Type == "Train") && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		//update busy list
		//TODO: Ensure we only pass message on if we're the leader of one of the nodes (requester or requestee)
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
			selfsoc.bootstrapsoc.Send(retmsg,0)

		}else {
			retmsg := encode(node.NodeName, "", "", "","", "Rejected", "","", "", "", dummy,"")
			selfsoc.bootstrapsoc.Send(retmsg,0)
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
		bestnode, notpicked := MasterNodeMet(nm, node, self,message)
		fmt.Println("BESTNODE: " + bestnode.NodeName)
		fmt.Println("NOTPICKED: " + notpicked)
		//Need to fill the dummy metric out with info about the best node
		dummy.NodeInf = msg.Result.NodeInf
		//put the best node name in msg.Receiver, and the NodeInfo in metric
		m := encode(msg.Sender, bestnode.NodeName, msg.Kind, msg.Job, msg.Value, "Selected", msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, notpicked)
		nodeSend(m, self)
	} else if msg.Type=="Hi" {
		updateNodeInfo(nm, msg.Sender, msg.Result.NodeInf)
		for k,v := range nm.Nodes{
			if v.Leader==true && v.Master == false{
				var dummy metric
				dummy.NodeInf = v
				up := encode(k, "", "", "","", "UpdateUptime", "","", "", "", dummy,"")
				nodeSend(up,self)
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
				msg := encode(leader,"","","","","Leader","","",nodeinfo.NodeAddr,nodeinfo.BootstrapPort,dummy,"")
				self.bootstrapsoc.Send(msg,0)
				//Req/Rep needs to send/receive
				//This means we need 1 send and 1 recv before the next send
				self.bootstrapsoc.Recv(0)
			}
		}
		endmsg := encode(leader,"","","","","Leader","","","","",dummy,"")
		self.bootstrapsoc.Send(endmsg,0)
	}else if msg.Type=="Board"{
		c:=0
		counter := make(map[string]bool)
		fmt.Print("##\n")
		bo := encode("Master","","Prime","","","BoardLeader","","","","",dummy,"")
		nodeSend(bo,self)
		size := len(getLeaders(nm)) - 1
		for _,child := range getLeaders(nm){
			counter[child] = false
		}
		var a map[string]int
		var a1 map[string]int
		for {
			s := MQpop(self.recvq)
				if s != nil {
					message := fmt.Sprint(s)
					m := decode(message)
					if m.Type == "BoardRequest" {

							counter[m.Sender]=true
							a = decodeRep(m.Value)
							a1 = decodeRep(m.Value)
							for k, v := range a {
								BoardH[k] = v
							}
							for k, v := range a1 {
								BoardP[k] = v
							}
							c++
							if (c==size){
								bo := encode("Master","","","",encodeRep(a1),"BoardReply",encodeRep(a),"","","",dummy,"")
								LeadNodeSend(bo,self)
							}

					}else  {
						MQpush(self.recvq, s)
					}

				}

		}
		}

}

func MasterNodeMet(nm NodeMap, node NodeInfo,self NodeSocket, msg string) (NodeInfo,string) {
	counter := make(map[string]bool)
	size := len(getChildren(nm)) - 1
	c := 0
	maxRep := -1
	var bestNode NodeInfo
	bestNode.NodeName = ""
	notpicked := ""
	nplist := make([]string, size+1, size+1)
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
						if(bestNode.NodeName != "") {
							nplist = append(nplist, bestNode.NodeName)
						}
						bestNode = m.Result.NodeInf
						maxRep = a
					} else {
						nplist = append(nplist, m.Receiver)
					}
					c = c + 1
					if (c==size){
						notpicked = stringulate(nplist)
						return bestNode, notpicked
					}


				}
				//check if we received all metrics from GLs and set bestnode accordingly
			} else if m.Type == "Request" {
				MQpush(self.recvq, s)
			}
		}
	}
	return bestNode,notpicked
}


func MessageHandler(selfname string, nm NodeMap, selfsoc NodeSocket){

	selfnode := nm.Nodes[selfname]

	s := MQpop(selfsoc.recvq)
	if s == nil{
		return
	}
	message := fmt.Sprint(s)
	m:= decode(message)
	if m.Receiver == selfname && m.Type == "Selected"{
		processRequestReceive(nm, selfname, selfsoc, message)
	} else if m.Receiver == selfname && m.Type == "Reply"{
		processRequestSend(nm.Nodes[selfname], selfsoc, message)
	} else if selfsoc.leader == true {
			println("Retransmitting " + message)
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

/*
func ReceiveResult(self NodeSocket,msg Message){
	addr := msg.Result.NodeInf.NodeAddr
	port := msg.Result.NodeInf.DataSendPort
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
*/

func SendResult(self NodeSocket, node NodeInfo, m Message){
	addr:=node.NodeAddr
	port:=node.DataSendPort
	send_sock:= establishClient(addr,port)
	msg := encode(m.Sender, m.Receiver,m.Kind,m.Job,m.Value,m.Type,m.ReceiverGroup,m.SenderGroup,m.Address,m.Port,m.Result,m.Input)
	send_sock.Send(msg,0)
	send_sock.Disconnect(addr + ":" + port)
	send_sock.Close()
}

