package main

import (
	"os"
	"io/ioutil"
	"strings"
	"fmt"
)

type Configs struct{
	Nodes		[]NodeInfo
}

type NodeInfo struct {
	NodeName	string
	NodeType	string
	NodeGroup	string
	NodeAddr	string
	SendPort	string
	RecvPort	string
	LeaderSendPort	string
	LeaderRecvPort	string
	Effort 		string
	DataRecvPort 	string
	DataSendPort	string
}

// Reads info from config file
func ReadConfig(configfile string) Configs {

	_, err := os.Stat(configfile)
	check(err)

	dat, err := ioutil.ReadFile(configfile)
	confstr := string(dat[:len(dat)])
	check(err)

	return ParseConfigString(confstr)
}

func ParseConfigString(raw string) Configs{
	var retconf Configs
	retconf.Nodes = []NodeInfo{}
	splitraw := strings.Split(raw,"\n")
	startind := -1
	endind := -1
	typestr := "none"
	for ind,lin := range splitraw {
		if lin == "Node" {
			typestr = "node"
			startind = ind
		} else if lin == "End" {
			endind = ind
		}
		if startind > 0 && endind > 0{
			if typestr == "node" {
				tmp := ParseNode(splitraw[startind:endind])
				retconf.Nodes = append(retconf.Nodes, tmp)
				startind = -1
				endind = -1
			}
		}
	}
	return retconf
}

func ParseNode(nodelines []string) NodeInfo{

	var retnode NodeInfo
	var err error
	for _,lin := range nodelines {
		if strings.Contains(lin,"NodeName"){
			retnode.NodeName = strings.Split(lin,"=")[1]
		} else if strings.Contains(lin,"NodeAddr"){
			retnode.NodeAddr = strings.Split(lin,"=")[1]
		} else if strings.Contains(lin,"NodeGroup"){
			retnode.NodeGroup = strings.Split(lin,"=")[1]
		} else if strings.Contains(lin,"SendPort"){
			retnode.SendPort =  strings.Split(lin,"=")[1]
			check(err)
		} else if strings.Contains(lin,"RecvPort"){
			retnode.RecvPort =  strings.Split(lin,"=")[1]
			check(err)
		} else if strings.Contains(lin,"Effort"){

			retnode.Effort =  strings.Split(lin,"=")[1]
			check(err)
		}else if strings.Contains(lin,"NodeType"){
			retnode.NodeType =  strings.Split(lin,"=")[1]
			check(err)
		}else if strings.Contains(lin,"LeaderRecvPort"){
			retnode.LeaderRecvPort =  strings.Split(lin,"=")[1]
			check(err)
		}else if strings.Contains(lin,"LeaderSendPort"){

			retnode.LeaderSendPort =  strings.Split(lin,"=")[1]
			fmt.Print(retnode.LeaderSendPort)
			check(err)
		}else if strings.Contains(lin,"lsp"){
			retnode.LeaderSendPort =  strings.Split(lin,"=")[1]

			check(err)
		}else if strings.Contains(lin,"lcp"){
			retnode.LeaderRecvPort =  strings.Split(lin,"=")[1]

			check(err)
		}else if strings.Contains(lin,"drp"){
			retnode.DataRecvPort =  strings.Split(lin,"=")[1]

			check(err)
		}else if strings.Contains(lin,"dsp"){
			retnode.DataSendPort =  strings.Split(lin,"=")[1]

			check(err)
		}




}
	return retnode
}