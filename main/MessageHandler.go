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
			msg := &Message{Sender: self, Receiver: m.Sender, Kind: m.Kind, Value: ms, Timestamp: getCurrentTimestamp(), Type:"Reply"}
			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			fmt.Print(string(b) + "\n")
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
	if (test.Type == "req"){
		msg := &Message{Sender: test.Sender, Receiver: test.Receiver, Kind: test.Kind, Value: test.Value, Timestamp: getCurrentTimestamp(),Type: "req"}
		b, err := json.Marshal(msg)
		check(err)
		//nodeSend(string(b))

	}
}