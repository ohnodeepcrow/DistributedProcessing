package main

import (
	"os"
	zmq4 "github.com/pebbe/zmq4"
	"strconv"
	"fmt"
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
	fmt.Println("Running as "+self.NodeName)
	fmt.Println("IP: " + self.NodeAddr)
	fmt.Println("Port1: " + self.SendPort)
	fmt.Println("Port2: " + self.RecvPort)
	fmt.Println("Group: " + self.NodeGroup)
	fmt.Println("Effort Level: " + self.Effort)
	cntxt,_ := zmq4.NewContext()
	eff, _ := strconv.Atoi(self.Effort)
	setEffort(eff)
	var ns NodeSocket
	if self.NodeName == "leader"{
		ns = establishLeader(cntxt, self)
	} else {
		ns = establishMember(cntxt, self, leader)
	}
	startIO(cntxt, ns, self)
}