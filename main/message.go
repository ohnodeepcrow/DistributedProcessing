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
	Result metric
	Input string
}

func encode(sender string,receiver string,kind string, value string,typ string, m metric, i string) string{
	msg := &Message{Sender: sender, Receiver: receiver, Kind: kind, Value: value, Timestamp: getCurrentTimestamp(), Type:typ, Result:m, Input:i}
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