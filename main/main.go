package main

import (
	"os"
)

func main(){
	if len(os.Args) < 2{
		println("Need to pass in config file location!")
		return
	}
	configfile := os.Args[1]
	var configs Configs = ReadConfig(configfile)
	println(configs.Nodes[1].NodeName)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}