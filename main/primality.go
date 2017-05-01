package main

import (
	"math/big"
	"math/rand"
	_ "fmt"
	"time"
)


var effort int

//generate candidate int64 for primality
func generateCandidate() *big.Int{
	tmp := rand.Int63()
	return big.NewInt(tmp)
}

func setEffort(i int){
	effort=i
	rand.Seed(time.Now().UTC().UnixNano())
}

//for node processing
func testPrime(num big.Int) metric{
	var m metric
	run:=rand.Intn(effort-0) + 0
	m.Perf= run
	isPrime := num.ProbablyPrime(run)
	m.IsPrime=isPrime
	//fmt.Println(isPrime)
	return m
}

//for verification of node result
func verifyPrime(num big.Int) bool{
	isPrime := num.ProbablyPrime(2000)
	return isPrime
}

func trainPrime(nm NodeMap, self NodeSocket,nodeinfo NodeInfo){
	for i:=0;i<10 ;i++  {
		var m metric
		msg := encode(nodeinfo.NodeName, nodeinfo.NodeName,"Prime",generateCandidate().String(),getCurrentTimestamp(), "Selected",nodeinfo.NodeGroup,"","","",m,"")

		processRequestReceive(nm, nodeinfo.NodeName, self ,msg )
	}
}