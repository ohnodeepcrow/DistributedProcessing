package main

import (
	"encoding/json"
	"fmt"
)

type Message struct{
	Sender string
	Receiver string
	Kind string
	Value string
	Timestamp string
	Type string
}

func encode(sender string,receiver string,kind string, value string,typ string) string{
	msg := &Message{Sender: sender, Receiver: receiver, Kind: kind, Value: value, Timestamp: getCurrentTimestamp(), Type:typ}
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}else{
		return string(b)
	}

}

func decode(in string) Message{
	res := []byte(in)
	var m Message
	json.Unmarshal(res,&m)
	return m
}

func IsMyMessage(self NodeInfo, message Message){
	if (message.Receiver == self.NodeName){
		return true
	}
	return false
}