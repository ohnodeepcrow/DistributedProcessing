package main
import (
	"encoding/json"
	"strconv"
	"math/big"
	"fmt"
)


func processRequest(self string, m Message) {
	if m.Receiver == self{
		if m.Kind == "Prime" && m.Type== "Request"{
			i,_:= strconv.ParseInt(m.Value,10,64)
			num:=big.NewInt(i)
			metric := testPrime(*num)
			ms:=metricString(metric)
			msg := encode(self, m.Sender,m.Kind,ms, "Reply")
			fmt.Print(string(msg) + "\n")
			//send
		} else if m.Kind == "Prime" && m.Type== "Reply" {
			//print the reply
		}

	}
}

/*Lead node will call this function after it received a message from a node. It will use send to retransmit the node. */
func LeadNodeRec(input string){
	res := []byte(input)
	var test Message

	json.Unmarshal(res,&test)
	if test.Type == "Request"{
		msg := encode(test.Sender, test.Receiver, test.Kind, test.Value,"Request")
		nodeSend(string(msg))

	}
}