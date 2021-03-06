package main

import (
	"time"
	"strconv"
)

func getCurrentTimestamp() string{
	return time.Now().Format("15:04:05")
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

func metricString (m metric) string{
	a:=strconv.FormatBool(m.IsPrime)
	b:=strconv.Itoa(m.Perf)
	c:="Prime: "+a+"\n"+"Effort: "+b
	return c
}

func hmetricString (m metric) string{
	a:=m.Val
	b:=m.hPerf.String()
	c:="Preimage: "+a+"\n"+"Performance Hit: "+b
	return c
}

//takes in a list of strings
//outputs that list with no duplicates as a semicolon-delimited string
func stringulate(inlist []string) string{
	//fmt.Print("STRINGULATE: ")
	//fmt.Println(inlist)
	size := len(inlist)
	ret := ""
	for i:=0; i < size; i++{
		dup := false
		for j := i+1; j < size; j++{
			if inlist[i] == inlist[j]{
				dup = true
				break
			}
		}
		if !dup && inlist[i] != ""{
			if ret == ""{
				ret = inlist[i]
			} else {
				ret = ret + ";" + inlist[i]
			}
		}
	}
	return ret
}