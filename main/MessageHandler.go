package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
	"github.com/pebbe/zmq4"
)


func processRequest(node NodeInfo, self NodeSocket, input string) {
	m:= decode(input)
		if m.Kind == "Prime" && m.Type== "Request"{
			i,_:= strconv.ParseInt(m.Value,10,64)
			println(m.Value)
			num:=big.NewInt(i)
			metric := testPrime(*num)
			ms:=metricString(metric)

			msg := encode(node.NodeName, m.Sender,m.Kind,ms, "Reply",node.NodeGroup,m.SenderGroup,false,node.NodeAddr,node.DataSendPort,metric,num.String())

			//fmt.Print(string(msg) + "\n")
			nodeSend(msg,self)
		} else if m.Kind == "Hash" && m.Type== "Request"{
			metric := crackHash(m.Value)
			ms:=hmetricString(metric)

			msg := encode(node.NodeName, m.Sender,m.Kind,ms, "Reply",node.NodeGroup,m.SenderGroup,false,node.NodeAddr,node.DataSendPort,metric,m.Value)

			//fmt.Print(string(msg) + "\n")

			nodeSend(msg,self)
		}
			else if m.Type== "Reply" {
			MQpush(self.appq, m)
		}
}


/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(node NodeInfo,self NodeSocket, m string){
	fmt.Print(m+"\n")
	msg:=decode(m)

	if (msg.ReceiverGroup != node.NodeGroup){
		if (msg.Flag == false) {
			LeadNodeSend(m, self) // lead node send this message to other lead node.
		}
		if(msg.Flag== true){
			//set the msg.receiver node as busy and inform all nodes
			nodeSend(m,self)
		}
	} else if (msg.ReceiverGroup == node.NodeGroup){
		nodeSend(m, self)
	}else if(msg.Kind=="Metric"){
		//reply to master node with the best node
	}

}

/*Master node will call this function after it received a message from a node. It will use send to retransmit the node. */
func MasterNodeRec(self NodeSocket, m string){
	fmt.Print(m+"\n")
	msg := decode(m)
	var dummy metric
	message := encode(msg.Sender,msg.Receiver,"Metric",msg.Value,msg.Type,msg.SenderGroup,msg.ReceiverGroup,false,msg.Address,msg.Port,dummy,msg.Value)
	nodeSend(message,self)
	bestnode :=selectNode()
	//put the best node in msg.Receiver
	message := encode(msg.Sender,bestnode,msg.Kind,msg.Value,msg.Type,msg.SenderGroup,msg.ReceiverGroup,true,msg.Address,msg.Port,dummy,msg.Value)
	nodeSend(message,self)

}

func MessageHandler(node NodeInfo, self NodeSocket){

	s := MQpop(self.recvq)
	if s == nil{
		return
	}
	message := fmt.Sprint(s)
	m:= decode(message)
	if m.Receiver == node.NodeName {
		processRequest(node, self, message)
		SendResult(self,node,m)
	} else if self.leader == true {
			//println("Retransmitting " + message)
			LeadNodeRec(node, self, message)
	} else if self.master == true {
		MasterNodeRec(self,message)
	}else {
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
	msg := encode(node.NodeName, m.Sender,m.Kind,ms, "Reply",node.NodeGroup,m.SenderGroup,false,node.NodeAddr,node.DataSendPort,metric,m.Value)
	fmt.Print(string(msg) + "\n")
	for {
		signal,_ := send_sock.Recv(zmq4.DONTWAIT)
		if (signal != ""){
			//send the result back through data socket.
			fmt.Print(signal+"\n")
			send_sock.Send(msg,0)

		}
		time.Sleep(time.Millisecond*50)
	}
}

func selectNode() string{
	//master finds the best node
}