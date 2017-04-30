package main

import (
"os"
zmq4 "github.com/pebbe/zmq4"
"strconv"
"fmt"
"sync"
"runtime"
)
func main(){
	if len(os.Args) < 3{
		println("Need to pass in config file location and node name!")
		return
	}

	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	wg.Add(2)

	configfile := os.Args[1]
	selfstr := os.Args[2]
	var configs Configs = ReadConfig(configfile)

	self := getNodeInfo(selfstr, configs)
	// currently we are assigning leader with their peer manually, need to figure out a better way to connect the peers.
	master := getNodeInfo("master", configs)
	fmt.Println("Running as "+self.NodeName)
	fmt.Println("IP: " + self.NodeAddr)
	fmt.Println("Port1: " + self.SendPort)
	fmt.Println("Port2: " + self.RecvPort)
	fmt.Println("Group: " + self.NodeGroup)
	fmt.Println("Effort Level: " + self.Effort)
	cntxt,_ := zmq4.NewContext()
	eff, _ := strconv.Atoi(self.Effort)
	setEffort(eff)

	setDict()

	var ns NodeSocket
	var myNodeMap NodeMap
	myNodeMap = initializeNode(self, master)

	self = myNodeMap.Nodes[selfstr]
	/*TODO:
		-If Leader: get group uptimes, group reputations, leader uptimes
			-Group uptimes happen when nodes join network
			-Group reputations happen when nodes join network
			-Leader uptimes happen via bootstrap message
		-If Root: get leader uptimes
			-Happens via leader join bootstrap
		-If Member: get group uptimes
			-Happens via bootstrap message
	*/
	println("Started Bootstrapping...")
	if self.NodeType == "master"{
		ns = establishMaster(cntxt, self)
	}else{
		ns  = BootStrap(cntxt,self,master)
		println("Finished Bootstrap!")
	}
	go startReceiver(ns)
	println("Started Receiver Thread")
	go startSender(ns)
	println("Started Sender Thread")
	go startMessageHandler(self.NodeName, myNodeMap, ns)
	println("Started Handler Thread")
	go startIO(cntxt, ns, self)
	println("Started IO")
	go startUI(ns, self)
	println("Started UI")
	wg.Wait()
}
