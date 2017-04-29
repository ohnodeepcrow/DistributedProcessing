package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
	"github.com/pebbe/zmq4"
	"strings"
)

var counter int
var leader string
func processRequestSend(node NodeInfo, self NodeSocket, input string) {
	m:= decode(input)
		 if m.Type=="Selected" && m.Sender!=m.Receiver{
			ReceiveResult(self,m)
		}
}


func processRequestReceive(node NodeInfo, self NodeSocket, input string) {
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
			nodeinf.RepMets.PrimeMetrics=updateReputation(node.RepMets.PrimeMetrics, metric, node.NodeName, primeScorer)
			msg = encode(node.NodeName, m.Sender,m.Kind,ms,m.Job, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,num.String())

		} else if m.Kind == "Hash" {
			metric = crackHash(m.Value)
			ms=hmetricString(metric)
			nodeinf.RepMets.HashMetrics=updateReputation(node.RepMets.HashMetrics, metric, node.NodeName, hashScorer)
			msg = encode(node.NodeName, m.Sender,m.Kind,m.Job,ms, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Value)

		}

		fmt.Println("Hash:")
		fmt.Println(nodeinf.RepMets.HashMetrics[nodeinf.NodeName])
		fmt.Println("Prime:")
		fmt.Println(nodeinf.RepMets.PrimeMetrics[nodeinf.NodeName])

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

	if msg.Type =="Request" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		LeadNodeSend(m, selfsoc) // group node forwards the request to master node
	} else if msg.Type=="Metric" && (msg.Kind=="Prime"||msg.Kind=="Hash"){
		
		//reply to master node with the best node
		//update busy list
		//master finds the best node
		bestname, bestscore := getBestFreeScore(node.RepMets, msg.Kind)
		setBusy(node.RepMets, bestname, msg.Job)
		var m metric
		retmsg := encode(node.NodeName, bestname, msg.Kind, msg.Job,strconv.Itoa(bestscore), "Metric", node.NodeGroup,"", node.NodeAddr, node.DataSendPort, m,"")
		LeadNodeSend(retmsg, selfsoc)
	} else if msg.Type=="Selected" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		//update busy list

		kids := getChildren(node.RepMets)
		nparr := strings.Split(msg.Input, ";")
		for _,kid := range kids{
			for _,np := range nparr{
				if kid == np && getBusyJob(node.RepMets, np) == msg.Job{
					setFree(node.RepMets, kid)
				}
			}
		}

		nodeSend(m, selfsoc)
	}else if msg.Type=="Update" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		setFree(node.RepMets,msg.Sender)
		if msg.Kind == "Prime"{
			updateReputation(node.RepMets.PrimeMetrics, msg.Result, msg.Sender, primeScorer)
		} else if msg.Kind == "Hash"{
			updateReputation(node.RepMets.HashMetrics, msg.Result, msg.Sender, hashScorer)
		}
	}else if msg.Type=="Hi" {
		updateUptime(nm, msg.Sender, msg.Result.NodeInf.Uptime)
	}else if msg.Type=="Bye"{
		clearUptime(nm, msg.Sender)
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
		message := encode(msg.Sender, msg.Receiver, "Metric", msg.Job, msg.Value, msg.Type, msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		nodeSend(message, self)
		bestnode := MasterNodeMet(node, self,message)

		//put the best node in msg.Receiver
		m := encode(msg.Sender, bestnode, msg.Kind, msg.Job, msg.Value, msg.Type, msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		LeadNodeSend(m, self)
	} else if msg.Type=="Hi" {
		updateUptime(nm, msg.Sender, msg.Result.NodeInf.Uptime)
	}else if msg.Type=="Bye"{
		clearUptime(nm, msg.Sender)
		var dummy metric
		var dummyni NodeInfo
		dummyni.Uptime = node.Uptime
		dummy.NodeInf = dummyni
		retmsg := encode(node.NodeName, "", "", "","", "UpdateUptime", "","", "", "", dummy,"")
		LeadNodeSend(retmsg, self)

	}else if msg.Type=="Boot"{
		if counter==0 || counter%7==0{
			//establish the node as leader
			nm.Nodes[msg.Sender]=msg.Result.NodeInf
			leader=msg.Sender
			counter++
			//return to the node that it is a leader
		}else {
			//establish member assign the last leader

		}

	}
}

func MasterNodeMet(node NodeInfo,self NodeSocket, msg string) string {
	var counter map[string]bool
	size := len(getChildren(node.RepMets))
	c := 0
	maxRep := -1
	bestNode := ""
	for _,child := range getChildren(node.RepMets){
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
						c = c + 1
						if (c==size){
							return bestNode
						}
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
		processRequestReceive(selfnode, selfsoc, message)
	}else if  m.Sender == selfname {
		processRequestSend(selfnode, selfsoc, message)

	} else if selfsoc.leader == true {
			//println("Retransmitting " + message)
			LeadNodeRec(selfname, nm, selfsoc, message)
	} else if selfsoc.master == true {
		MasterNodeRec(selfnode,nm,selfsoc,message)
	}else if m.Kind == "UpdateUptime"{
		updateUptime(nm, m.Sender, m.Result.NodeInf.Uptime)
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
	soc :=establishClient(addr,port,self)
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
	send_sock:= establishServer(addr,port,self)
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

