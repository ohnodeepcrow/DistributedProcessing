package main

import (
	"time"
	"strconv"
	"math/big"
	"fmt"
)

type metric struct {
	Perf int
	IsPrime bool
	hPerf time.Duration
	Val string
	NodeInf NodeInfo
}

type Reputation struct {
	Score		int
	Count		int
	Correct 	int
}

//Map that maps node names to the first time that node was seen
type NodeMap struct{
	Nodes map[string]NodeInfo
}

func initializeNode(self NodeInfo, master NodeInfo) NodeMap{
	nm := newNodeMap(self.NodeName)
	nm.Nodes[master.NodeName] = master
	self.PrimeMetric, self.HashMetric, self.Busy = newRepMetrics()
	return nm
}

func newNodeInfo(name string, ut time.Time) NodeInfo{
	var ret NodeInfo
	ret.NodeName = name
	ret.Uptime = ut
	//TODO: ensure all needed fields are filled in
	return ret
}

func newMasterBoard(self NodeInfo) map[string]int{
	var ret map[string]int
	ret = make(map[string]int)
	ret[self.NodeName]=0
	return ret
}

func newNodeMap(selfname string)NodeMap{
	var ret NodeMap
	ret.Nodes = make(map[string]NodeInfo)
	var tmp NodeInfo
	tmp.Uptime = time.Now()
	ret.Nodes[selfname] = tmp
	return ret
}

func newRepMetrics() (Reputation, Reputation, string){
	var tmp1 Reputation
	var tmp2 Reputation
	tmp2.Correct = 0
	tmp2.Count = 0
	tmp2.Score = 0
	tmp1.Correct = 0
	tmp1.Count = 0
	tmp1.Score = 0
	return tmp1, tmp2, ""
}

func getChildren(nm NodeMap) []string{
	keys := make([]string, 0, len(nm.Nodes))
	for k,v := range nm.Nodes {
		if (!v.Leader) && (!v.Master){
			keys = append(keys, k)
		}
	}
	return keys
}

func getLeaders(nm NodeMap) []string{
	keys := make([]string, 0, len(nm.Nodes))
	for k,v := range nm.Nodes {
		if (!v.Master){
			keys = append(keys, k)
		}
	}
	return keys
}

func getBestFreeScore(nm NodeMap, probtype string) (string, int){
	bestscore := -1
	bestname := "NOTFOUND!"
	kids := getChildren(nm)
	for _,k := range kids{
		tmp := nm.Nodes[k]
		//println("TMP: " + tmp.NodeName)
		//println(tmp.HashMetric.Score)
		//println(tmp.PrimeMetric.Score)
		if tmp.Busy == ""{
			if probtype == "Hash"{
				if tmp.HashMetric.Score > bestscore{
					bestscore = tmp.HashMetric.Score
					bestname = tmp.NodeName
				}
			} else if probtype == "Prime"{
				if tmp.PrimeMetric.Score > bestscore{
					bestscore = tmp.PrimeMetric.Score
					bestname = tmp.NodeName
				}
			}
		}
	}
	return bestname, bestscore
}

func getBusyJob(nm NodeMap, nodename string) string{
	return nm.Nodes[nodename].Busy
}

func setBusy(nm NodeMap, nodename string, jid string){
	tmp := nm.Nodes[nodename]
	tmp.Busy = jid
	nm.Nodes[nodename] = tmp
}

func setFree(nm NodeMap, nodename string){
	tmp := nm.Nodes[nodename]
	tmp.Busy = ""
	nm.Nodes[nodename] = tmp
}

//If a node doesn't currently have a node's nodeinfo, add one
func updateNodeInfo(nm NodeMap, name string, newinfo NodeInfo) bool{
	_, ok := nm.Nodes[name]
	if !ok{
		nm.Nodes[name] = newinfo
		return true
	}
	return false
}

//Remove the nodeinfo associated with a node
func clearNodeInfo(nm NodeMap, name string){
	delete(nm.Nodes, name)
}

func getLongestUptime(nm NodeMap) (string,time.Time){
	longest := time.Now()
	lname := ""
	for k,v := range nm.Nodes{
		if longest.After(v.Uptime){
			longest = v.Uptime
			lname = k
		}
	}
	return lname, nm.Nodes[lname].Uptime
}

//Scorer should take in the current reputation and the new result and update the reputation as a result
func updateReputation(repmets Reputation, newmet metric, node string, scorer func(nm metric, rp Reputation) Reputation) Reputation{
	return scorer(newmet, repmets)
}

//The score for hashing is the average time it takes to generate a collision
//It doesn't use correctness currently
func hashScorer(met metric, rep Reputation) Reputation{
	rep.Count += 1
	newscore := rep.Score / rep.Count
	newscore += int(met.hPerf)
	//fmt.Print(newscore)
	//fmt.Println(" HASHSCORER debug")
	newscore = newscore/ rep.Count
	rep.Score = newscore
	return rep
}

//the score for primality is the average number of correct assessments out of 100,000
//The score is score = correct/count
func primeScorer(met metric, rep Reputation) Reputation{
	i, _ := strconv.ParseInt(met.Val,10,64)
	test := big.NewInt(i)
	if met.IsPrime == testPrime(*test).IsPrime{
		rep.Correct += 1
	}
	rep.Count += 1
	rep.Score = (rep.Correct*1.0) / (rep.Count *1.0)* 100
	return rep
}

func getRepBoard(self NodeSocket,nodeinf NodeInfo, nm NodeMap) (map[string]int, map[string]int){
	var dummy metric
	c:=0
	bo := encode("Master","","Prime","","","BoardLeader","","","","",dummy,"")
	nodeSend(bo,self)
	size := len(getLeaders(nm))

	fmt.Print("SIZE: ")
	fmt.Println(size)

	fmt.Print("LEADERS: ")
	fmt.Println(getLeaders(nm))

	var a map[string]int
	var a1 map[string]int
	a = make(map[string]int)
	a1 = make(map[string]int)

	var BoardH map[string]int
	var BoardP map[string]int

	BoardH = make(map[string]int)
	BoardP = make(map[string]int)

	for {
		s := MQpop(self.recvq)
		if s != nil {
			message := fmt.Sprint(s)
			m := decode(message)
			if m.Type == "BoardReply" {
				a = decodeRep(m.Value)
				a1 = decodeRep(m.SenderGroup)
				for k, v := range a {
					BoardH[k] = v
				}
				for k, v := range a1 {
					BoardP[k] = v
				}
				c++
				if (c>=size){
					fmt.Println("BREAKING")
					break
				}

			}else  {
				MQpush(self.recvq, s)
			}

		}

	}
	return BoardP,BoardH
}