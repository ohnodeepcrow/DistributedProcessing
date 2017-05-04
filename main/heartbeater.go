package main

import (
	"time"
	"fmt"
	"sync"
	"strings"
)

var timeOutThreshold time.Duration = time.Millisecond * 2000

type heartbeat struct{
	Sender string
	Timestamp time.Time
}

func HandleTimeout(nm NodeMap, selfname string, soc NodeSocket, othername string) {
	self := nm.Nodes[selfname]
	other,exists := nm.Nodes[othername]
	if !exists{
		fmt.Println("DOESN'T EXIST: " + othername)
		return
	}

	if other.Leader && self.Master{
		delete(nm.Nodes, othername)
		counter--
		fmt.Println("Deleted " + othername + " from our group!")
		//The other one is below us... just clean up our entries
	} else if self.Leader{
		delete(nm.Nodes, othername)
		counter--
		fmt.Println("Deleted " + othername + " from our group!")
		//The other one is below us... just clean up our entries
	} else {
		panic("OUR LEADER WENT DOWN")
	}
}

func hbencode(sender string, timestamp time.Time) string{
	ret := "HBmsg|" + sender + "|"
	ret += timestamp.Format(time.RFC3339)
	return ret
}

func isHbString(in string) bool{
	if len(in) < 6{
		return false
	}
	if in[0:5] == "HBmsg"{
		return true
	}
	return false
}

func hbdecode(in string) heartbeat{
	var ret heartbeat
	arr := strings.Split(in, "|")
	ret.Sender = arr[1]
	ret.Timestamp,_ = time.Parse( time.RFC3339, arr[2])
	return ret
}

func heartbeatSender(socket NodeSocket, selfname string){
	sender := selfname
	oldtime := time.Now()

	for {
		if time.Since(oldtime) < timeOutThreshold/5{
			time.Sleep(timeOutThreshold/50)
			continue
		}
		oldtime = time.Now()
		h := hbencode(sender, oldtime)
		if socket.master {
			MQpush(socket.sendq, h)
		} else if socket.leader{
			MQpush(socket.sendq, h)
			MQpush(socket.lsendq, h)
		} else {
			MQpush(socket.sendq, h)
		}
	}
}

func heartbeatUpdater(soc NodeSocket, hbmap map[string]time.Time, mut sync.Mutex){
	for{
		hstr := MQpop(soc.hbqueue)
		if hstr == nil{
			time.Sleep(timeOutThreshold/50)
			continue
		}
		mut.Lock()
		hb := hbdecode(hstr.(string))
		_,there := hbmap[hb.Sender]
		if !there{
			fmt.Println("LEARNED ABOUT NODE " + hb.Sender + "!")
		}
		hbmap[hb.Sender] = hb.Timestamp
		mut.Unlock()
	}
}

func selfTimeout(selfname string, socket NodeSocket){
	var emptytime time.Time
	h := hbencode(selfname, emptytime)
	if socket.master {
		MQpush(socket.sendq, h)
	} else if socket.leader{
		MQpush(socket.sendq, h)
		MQpush(socket.lsendq, h)
	} else {
		MQpush(socket.sendq, h)
	}
}

func heartbeatChecker(selfname string, soc NodeSocket, hbmap map[string]time.Time, mut sync.Mutex){
	for{
		time.Sleep(timeOutThreshold/5)
		mut.Lock()
		for k,v := range hbmap{
			if time.Since(v) > timeOutThreshold{
				fmt.Println("TIMEOUT DETECTED FOR NODE " + k + "!")
				delete(hbmap, k)
				var dummy metric
				msg := encode(selfname, selfname, "", "", k, "TimeoutDetected", "", "", "", "", dummy, k)
				MQpush(soc.recvq, msg)
			}
		}
		mut.Unlock()
	}
}

func startHeartbeatService(selfname string, soc NodeSocket){
	hbmap := make(map[string]time.Time)
	var mut sync.Mutex
	go heartbeatSender(soc, selfname)
	go heartbeatUpdater(soc, hbmap, mut)
	go heartbeatChecker(selfname, soc, hbmap, mut)
}