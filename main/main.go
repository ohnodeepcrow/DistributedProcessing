package main

import (
	"os"
	"github.com/pebbe/zmq4"
	"time"
	"strconv"
)

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
	if self.NodeName == "leader"{
		soc := configPub(cntxt, self)
		sendPub(soc)
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

func sendPub(soc *zmq4.Socket){
	for i :=0 ; ;i++  {
		out := strconv.Itoa(i)
		soc.Send(out, 0)
		time.Sleep(time.Duration(time.Second)*2)
	}
}

func configSub(context *zmq4.Context) *zmq4.Socket{
	soc, err := context.NewSocket(zmq4.SUB)
	check(err)
	soc.Connect("tcp://127.0.0.1:13370")
	soc.SetSubscribe("")
	return soc
}

func receiveSub(soc *zmq4.Socket){
	oldtmp := ""
	for{
		tmp,_ := soc.Recv(0)
		if tmp != oldtmp {
			println(tmp)
			oldtmp = tmp
		}

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