package main

import (
	"os"
	"github.com/pebbe/zmq4"
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
	println(self.SendPort)
	println(self.RecvPort)
	println(self.NodeGroup)
	cntxt,_ := zmq4.NewContext()
	if self.NodeName == "leader"{
		nsoc := establishLeader(cntxt, self)
		for {
			nodeSend(self.NodeName, nsoc.sendsock)
			nodeReceive(nsoc.recvsock)
		}
	} else {
		nsoc := establishMember(cntxt, self, getNodeInfo("leader", configs))
		for {
			nodeSend(self.NodeName, nsoc.sendsock)
			nodeReceive(nsoc.recvsock)
		}
	}


}

func check(e error) {
	if e != nil {
		panic(e)
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