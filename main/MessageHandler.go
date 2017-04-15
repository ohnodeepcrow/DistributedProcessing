package main
import (
	"strconv"
	"math/big"
	"fmt"
	"time"
)


func processRequest(node NodeInfo, self NodeSocket, input string) {
	m:= decode(input)
		if m.Kind == "Prime" && m.Type== "Request"{
			i,_:= strconv.ParseInt(m.Value,10,64)
			num:=big.NewInt(i)
			metric := testPrime(*num)
			ms:=metricString(metric)
			msg := encode(node.NodeName, m.Sender,m.Kind,ms, "Reply")
			fmt.Print(string(msg) + "\n")
			nodeSend(msg,self)
		} else if m.Kind == "Prime" && m.Type== "Reply" {
			MQpush(self.appq, m)
		}

}

/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(self NodeSocket, m string){
	nodeSend(m,self)
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
	} else if self.leader == true {
			//println("Retransmitting " + message)
			LeadNodeRec(self, message)
	} else {
		return //drop message
	}
}

func startMessageHandler(node NodeInfo, self NodeSocket){
	for{
		time.Sleep(time.Millisecond*50)
		MessageHandler(node,self)
	}
}