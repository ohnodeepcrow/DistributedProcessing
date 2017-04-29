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
	leader1 := getNodeInfo("leader1", configs)
	leader2 := getNodeInfo("leader2", configs)
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
	var myleader NodeInfo
	if self.NodeType == "leader"{
		var ma NodeInfo
		ma = master
		ns = establishLeader(cntxt,self,ma)
	} else if self.NodeType == "master"{
		ns = establishMaster(cntxt, self)
	}else{
		if (self.NodeGroup == "group1") {
			myleader = leader1
		}else if (self.NodeGroup == "group2"){
			myleader = leader2
		}
		ns = establishMember(cntxt, self, myleader)

		var dummy metric
		dummy.NodeInf=self
		//say Hi
		msg := encode(self.NodeName, "", "",getCurrentTimestamp(),"","Hi","","","","",dummy,"")
		nodeSend(string(msg), ns)
	}
	go startReceiver(ns)
	go startSender(ns)
	go startMessageHandler(self.NodeName, myNodeMap, ns)
	go startIO(cntxt, ns, self)
	go startUI(ns, self)

	wg.Wait()
}