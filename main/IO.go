package main

import (
	"os"
	"github.com/pebbe/zmq4"
	"time"
	"encoding/json"
	"bufio"
	"strings"
	"fmt"

)

func startIO(cntxt *zmq4.Context, self NodeSocket, nodeinfo NodeInfo){
	reader :=bufio.NewReader(os.Stdin)
	for {
		fmt.Print("(s)end/(r)eceive/(g)enerate")
		input, _ := reader.ReadString('\n')
		if input == "s" {
			fmt.Print("Enter Receiver:")

			fmt.Print("->\n")
			text, _ := reader.ReadString('\n')
			fmt.Print("Enter Kind and Value:")
			fmt.Print("->\n")
			text1, _ := reader.ReadString('\n')
			var a [2]string
			a[0] = strings.Split(text, " ")[0]
			a[1] = strings.Trim(strings.Split(text1, " ")[1], "\n")

			msg := &Message{sender: nodeinfo.NodeName, receiver: text, kind: a[0], value: a[1], timestamp: getCurrentTimestamp()}
			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return;
			}
			fmt.Print(string(b) + "\n")
			//nodeSend(soc, string(b))

		}
		else if input=="r" {

		}

		else if input=="g"{
			fmt.print(generateCandidate())
		}
	}
}
