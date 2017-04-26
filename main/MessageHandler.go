package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
	"github.com/pebbe/zmq4"
	"strings"
)


func processRequestSend(node NodeInfo, self NodeSocket, input string) {
	m:= decode(input)
		 if m.Type=="Selected"{
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
			println(m.Value)
			num:=big.NewInt(i)
			metric = testPrime(*num)
			ms=metricString(metric)

			msg = encode(node.NodeName, m.Sender,m.Kind,ms,m.Job, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,num.String())

		} else if m.Kind == "Hash" {
			metric = crackHash(m.Value)
			ms=hmetricString(metric)

			msg = encode(node.NodeName, m.Sender,m.Kind,m.Job,ms, "Reply",node.NodeGroup,m.SenderGroup,node.NodeAddr,node.DataSendPort,metric,m.Value)

		}
		SendResult(self,node,decode(msg))
		updatemsg := encode(node.NodeName, "",m.Kind, ms, m.Job,"Update",node.NodeGroup,"","","",metric,"")
		nodeSend(updatemsg, self)
	}
}



/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(node NodeInfo,self NodeSocket, m string){
	fmt.Print(m+"\n")
	msg:=decode(m)

	if msg.Type =="Request" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		LeadNodeSend(m, self) // group node forwards the request to master node
	} else if msg.Type=="Metric" && (msg.Kind=="Prime"||msg.Kind=="Hash"){
		
		//reply to master node with the best node
		//update busy list
		//master finds the best node
		bestname, bestscore := getBestFreeScore(node.RepMets)
		setBusy(node.RepMets, bestname, msg.Job)
		var m metric
		retmsg := encode(node.NodeName, bestname, msg.Kind, msg.Job,strconv.Itoa(bestscore), "Metric", node.NodeGroup,"", node.NodeAddr, node.DataSendPort, m,"")
		LeadNodeSend(retmsg, self)
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

		nodeSend(m,self)
	}else if msg.Type=="Update" && (msg.Kind=="Prime"||msg.Kind=="Hash") {
		setFree(node.RepMets,msg.Sender)
		if msg.Kind == "Prime"{
			updateReputation(node.RepMets, msg.Result, msg.Sender, primeScorer)
		} else if msg.Kind == "Hash"{
			updateReputation(node.RepMets, msg.Result, msg.Sender, hashScorer)
		}
	}

}

/*Master node will call this function after it received a message from a node. It will use send to retransmit the node. */
func MasterNodeRec(node NodeInfo,self NodeSocket, m string){
	fmt.Print(m+"\n")
	msg := decode(m)
	var dummy metric
	if msg.Type=="Request" {
		message := encode(msg.Sender, msg.Receiver, "Metric", msg.Job, msg.Value, msg.Type, msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		nodeSend(message, self)
		bestnode := MasterNodeMet(node, self)

		//put the best node in msg.Receiver
		m := encode(msg.Sender, bestnode, msg.Kind, msg.Job, msg.Value, msg.Type, msg.SenderGroup, msg.ReceiverGroup, msg.Address, msg.Port, dummy, msg.Value)
		nodeSend(m, self)
	}
}

func MasterNodeMet(node NodeInfo,self NodeSocket) string {
	for {
		s := MQpop(self.recvq)
		if s != nil {
			message := fmt.Sprint(s)
			m := decode(message)
			if m.Type == "Metric" {
				//check if we received all metrics from GLs and set bestnode accordingly
			} else if m.Type == "Request" {
				MQpush(self.recvq, s)
			}
		}
	}
	return ""
}


func MessageHandler(node NodeInfo, self NodeSocket){

	s := MQpop(self.recvq)
	if s == nil{
		return
	}
	message := fmt.Sprint(s)
	m:= decode(message)
	if m.Receiver == node.NodeName {
		processRequestReceive(node, self, message)
	}else if  m.Sender == node.NodeName {
		processRequestSend(node, self, message)

	} else if self.leader == true {
			//println("Retransmitting " + message)
			LeadNodeRec(node, self, message)
	} else if self.master == true {
		MasterNodeRec(node,self,message)
	}else{
		return //drop message
	}
}

func startMessageHandler(node NodeInfo, self NodeSocket){
	for{
		time.Sleep(time.Millisecond*50)
		MessageHandler(node,self)
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
	fmt.Print(string(msg) + "\n")
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
<<<<<<< HEAD


func selectNode(job string){


}
=======
>>>>>>> c0149d461d5cd00e8182d13a93d0f8778e362d68
