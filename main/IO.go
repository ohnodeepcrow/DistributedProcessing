package main

import (
	"os"
	"github.com/pebbe/zmq4"
	"encoding/json"
	"bufio"
	"strings"
	"fmt"

)

func startIO(cntxt *zmq4.Context, self NodeSocket, nodeinfo NodeInfo){
	reader :=bufio.NewReader(os.Stdin)
	for {
		fmt.Print("(s)end/(r)eceive/(g)enerate\n")
		input, _ := reader.ReadString('\n')
		input= strings.Trim(input, "\n")
		if input == "s" {
			fmt.Print("Enter Receiver:\n")
			text, _ := reader.ReadString('\n')
			text= strings.Trim(text, "\n")
			fmt.Print("Enter Kind and Value:\n")
			text1, _ := reader.ReadString('\n')
			var a [2]string
			a[0] = strings.Split(text1, " ")[0]
			a[1] = strings.Trim(strings.Split(text1, " ")[1], "\n")

			msg := &Message{Sender: nodeinfo.NodeName, Receiver: text, Kind: a[0], Value: a[1], Timestamp: getCurrentTimestamp(), Type:"Request"}
			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return;
			}
			fmt.Print(string(b) + "\n")
			//nodeSend(soc, string(b))

		} else if input=="r" {
			nodeReceive(self)
			ml := MQpopAll(self.recvq)
			if ml.Front() == nil{
				fmt.Println("No Messages!")
			}
			for n :=  ml.Front(); n != nil ; n = n.Next(){
				s := fmt.Sprint(n.Value)
				if s == "" {
					continue
				}
				res := []byte(s)
				var test Message

				json.Unmarshal(res,&test)

				fmt.Print("===============Receive Message==========")
				fmt.Print("kind: " + test.Kind)
				fmt.Print("value: " + test.Value)
				fmt.Print("sender: " + test.Sender)
				fmt.Println()
			}
		} else if input=="g"{
			fmt.Println(generateCandidate())
		}
	}
}
