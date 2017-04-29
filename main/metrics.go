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

//Maps node name/ID to Reputation and busy status
type RepMetrics struct {
	HashMetrics map[string]Reputation
	PrimeMetrics map[string]Reputation
	Busy map[string]string
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
	self.RepMets = newRepMetrics(self.NodeName)
	return nm
}

func newNodeInfo(name string, ut time.Time) NodeInfo{
	var ret NodeInfo
	ret.NodeName = name
	ret.Uptime = ut
	//TODO: ensure all needed fields are filled in
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

func newRepMetrics(name string)RepMetrics{
	var ret RepMetrics
	ret.PrimeMetrics = make(map[string]Reputation)
	ret.HashMetrics = make(map[string]Reputation)
	ret.Busy = make(map[string]string)
	var tmp Reputation
	tmp.Correct = 0
	tmp.Count = 0
	tmp.Score = 0
	ret.PrimeMetrics[name] = tmp
	ret.HashMetrics[name] = tmp
	return ret
}

func getChildren(metrics RepMetrics) []string{
	keys := make([]string, len(metrics.Busy))

	i := 0
	for k,_ := range metrics.Busy {
		keys[i] = k
		i++
	}
	return keys
}

func getBestFreeScore(metrics RepMetrics, probtype string) (string, int){
	bestscore := 0
	bestname := ""
	var probmap map[string]Reputation
	probmap = metrics.PrimeMetrics
	if probtype == "Hash"{
		probmap = metrics.HashMetrics
	}
	for k,v := range probmap{
		if (v.Score > bestscore) && (metrics.Busy[k] == ""){
			bestname = k
			bestscore = v.Score
		}
	}
	return bestname, bestscore
}

func getBusyJob(metrics RepMetrics, nodename string) string{
	return metrics.Busy[nodename]
}

func setBusy(metrics RepMetrics, nodename string, jid string){
	metrics.Busy[nodename] = jid
}

func setFree(metrics RepMetrics, nodename string){
	metrics.Busy[nodename] = ""
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
func updateReputation(repmets map[string]Reputation, newmet metric, node string, scorer func(nm metric, rp Reputation) Reputation) map[string]Reputation{
	rep, ok := repmets[node]
	if !ok{
		return nil
	}
	rep = scorer(newmet, rep)
	repmets[node] = rep
	return repmets
}

//The score for hashing is the average time it takes to generate a collision
//It doesn't use correctness currently
func hashScorer(met metric, rep Reputation) Reputation{
	fmt.Println(met.hPerf)
	fmt.Println("debug")
	rep.Count += 1
	newscore := rep.Score / rep.Count
	newscore += int(met.hPerf)
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