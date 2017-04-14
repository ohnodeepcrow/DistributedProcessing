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
	Kind string
	Value string
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
	println(self.NodePort)
	println(self.NodeGroup)
	cntxt,_ := zmq4.NewContext()
	reader :=bufio.NewReader(os.Stdin)
	if self.NodeName == "leader"{
		soc := configPub(cntxt, self)
		for {
			fmt.Print("kind and value:")
			fmt.Print("->\n")
			text, _ := reader.ReadString('\n')
			var a [2]string
			a[0] = strings.Split(text, " ")[0]
			a[1] = strings.Trim(strings.Split(text, " ")[1],"\n")
			msg := &Message{Kind:"prime",Value:a[1]}
			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return;
			}
			fmt.Print(string(b)+"\n")
			sendPub(soc, string(b))
		}
	} else {
		soc := configSub(cntxt)
		receiveSub(soc)
	}


}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func configPub(context *zmq4.Context, self NodeInfo) *zmq4.Socket{
	soc, err := context.NewSocket(zmq4.PUB)
	check(err)
	endpoint := "tcp://" + self.NodeAddr + ":" + self.NodePort
	//epgm://[IP of local interface];[multicast group IP]:[multicast port]
	soc.Bind(endpoint)
	return soc
}

func sendPub(soc *zmq4.Socket, msg string){


	soc.Send(msg, 0)
	time.Sleep(time.Duration(time.Second)*2)

}

func configSub(context *zmq4.Context) *zmq4.Socket{
	soc, err := context.NewSocket(zmq4.SUB)
	check(err)
	soc.Connect("tcp://127.0.0.1:13370")
	soc.SetSubscribe("")
	return soc
}

func receiveSub(soc *zmq4.Socket){

	for{
		tmp,_ := soc.Recv(0)

		res := []byte(tmp)
		var test Message
		json.Unmarshal(res,&test)
		fmt.Print("Kind: "+test.Kind +"\n")
		fmt.Print("Value: "+test.Value +"\n")

	}
}

func getNodeInfo(self string, config Configs) NodeInfo{
	for i := 0; i < len(config.Nodes); i++ {
		if config.Nodes[i].NodeName == self{
			return config.Nodes[i]
		}
	}
	panic("Node doesn't exist!")
}