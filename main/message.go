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

	SenderGroup string
	ReceiverGroup string
	Flag bool
	Address string
	Port string
	Result metric
	Input string
}

func encode(sender string,receiver string,kind string, value string,typ string,sgp string,rgp string, flag bool, addr string, port string,m metric, i string) string{
	msg := &Message{Sender: sender, Receiver: receiver, Kind: kind, Value: value, Timestamp: getCurrentTimestamp(), Type:typ,SenderGroup: sgp,ReceiverGroup:rgp, Flag:flag,Address:addr,Port:port, Result:m, Input:i}
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

func IsMyMessage(self NodeInfo, message Message) bool{
	if (message.Receiver == self.NodeName){
		return true
	}
	return false
}