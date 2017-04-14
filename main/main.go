package main

import (
	"os"
	zmq4 "github.com/pebbe/zmq4"
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
	leader := getNodeInfo("leader", configs)
	println("Running as " + self.NodeName)
	println("IP: " + self.NodeAddr)
	println("Port1: " + self.SendPort)
	println("Port2: " + self.RecvPort)
	println("Group: " + self.NodeGroup)
	cntxt,_ := zmq4.NewContext()
	var ns NodeSocket
	if self.NodeName == "leader"{
		ns = establishLeader(cntxt, self)
	} else {
		ns = establishMember(cntxt, self, leader)
	}
	startIO(cntxt, ns, self)
}