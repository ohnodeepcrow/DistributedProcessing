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
	Job string
	Type string

	SenderGroup string
	ReceiverGroup string
	Address string
	Port string
	Result metric
	Input string
}

func encode(sender string,receiver string,kind string, job string,value string,typ string,sgp string,rgp string, addr string, port string,m metric, i string) string{
	msg := &Message{Sender: sender, Receiver: receiver, Kind: kind, Value: value, Job: job, Type:typ,SenderGroup: sgp,ReceiverGroup:rgp,Address:addr,Port:port, Result:m, Input:i}
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

func encodeRep(rep map[string]int) string{
	b, err := json.Marshal(rep)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}else{
		return string(b)
	}

}
func decodeRep(rep string) map[string]int{
	res := []byte(rep)
	var m map[string]int
	json.Unmarshal(res,&m)
	return m
}