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
type Message struct{
	sender string
	receiver string
	kind string
	value string
	timestamp string
}
func main(){
	if len(os.Args) < 3{
		println("Need to pass in config file location and node name!")
		return
	}
	configfile := os.Args[1]
	selfstr := os.Args[2]
	var configs Configs = ReadConfig(configfile)
	self := getNodeInfo(selfstr, configs)
	println(self.NodeName)
	println(self.NodeAddr)
	println(self.SendPort)
	println(self.RecvPort)
	println(self.NodeGroup)
	cntxt,_ := zmq4.NewContext()
	reader :=bufio.NewReader(os.Stdin)
	for {	fmt.Print("Enter Receiver:")
		fmt.Print("->\n")
		text, _ := reader.ReadString('\n')
		fmt.Print("Enter Kind and Value:")
		fmt.Print("->\n")
		text1, _ := reader.ReadString('\n')
		var a [2]string
		a[0] = strings.Split(text, " ")[0]
		a[1] = strings.Trim(strings.Split(text1, " ")[1],"\n")

		msg := &Message{sender: self.NodeName, receiver: text,kind:a[0],value:a[1], timestamp:time.Now().Format("15:04:05")}
		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return;
		}
		fmt.Print(string(b)+"\n")
		//nodeSend(soc, string(b))
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
