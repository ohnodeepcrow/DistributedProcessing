package main

import (
	"math/big"
	"strconv"
	"fmt"
	"encoding/json"
)

type Message struct{
	Sender string
	Receiver string
	Kind string
	Value string
	Timestamp string
	Type string
}

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
	panic("Node doesn't exist!")
}

func metricString (m metric) string{
	a:=strconv.FormatBool(m.IsPrime)
	b:=strconv.Itoa(m.Perf)
	c:=a+"\n"+b
	return c
}