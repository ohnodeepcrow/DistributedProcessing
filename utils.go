package main

import "time"

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