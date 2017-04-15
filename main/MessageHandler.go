package main
import (
	"encoding/json"
)
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