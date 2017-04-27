package main

import (
	"os"
	"github.com/pebbe/zmq4"
	"bufio"
	"strings"
	"fmt"
)


func startIO(cntxt *zmq4.Context, self NodeSocket, nodeinfo NodeInfo){

	reader :=bufio.NewReader(os.Stdin)
	for {	fmt.Print("\n")
		fmt.Print("(s)end/(r)eceive/(g)enerate\n")
		input, _ := reader.ReadString('\n')
		input= strings.Trim(input, "\n")
		if input == "s" {
			fmt.Print("Enter processing node and group:\n")
			text, _ := reader.ReadString('\n')
			var b [2]string
			b[0] = strings.Split(text, " ")[0]
			b[1] = strings.Trim(strings.Split(text, " ")[1], "\n")
			fmt.Print("Enter Test and Value:\n")
			text1, _ := reader.ReadString('\n')
			var a [2]string
			a[0] = strings.Split(text1, " ")[0]
			a[1] = strings.Trim(strings.Split(text1, " ")[1], "\n")


			var dummy metric
			msg := encode(nodeinfo.NodeName, b[0], a[0],a[1],"Request",nodeinfo.NodeGroup,b[1],"","",dummy,a[1])



			nodeSend(string(msg), self) // if leader--> group, if member-->leader



		} else if input=="r" {
			ml := MQpopAll(self.appq)
			if ml.Front() == nil{
				fmt.Println("No Messages!")
			}
			for n :=  ml.Front(); n != nil ; n = n.Next(){
				test := n.Value.(Message)
				/*
				fmt.Println("====Results====")
				fmt.Println("Test: " + test.Kind)
				fmt.Println(test.Value)
				fmt.Println("Processed By: " + test.Sender)
				fmt.Println()
				*/
				go ReceiveResult(self,test)
			}
		} else if input=="g"{
			fmt.Println(generateCandidate())
		}
	}
}

