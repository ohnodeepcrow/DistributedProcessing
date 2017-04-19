package main

import (
	"math/big"
	"math/rand"
	"time"
)


var effort int
func generateCandidate() *big.Int{
	rand.Seed(time.Now().UTC().UnixNano())
	tmp := rand.Int63()
	if tmp % 2 == 0{
		tmp -= 1
	}
	return big.NewInt(tmp)
}

func setEffort(i int){
	effort=i
	rand.Seed(time.Now().UTC().UnixNano())
}
func testPrime(num big.Int) metric{
	var m metric
	m.Val = num.String()
	run:=rand.Intn(effort-0) + 0
	m.Perf= run
	isPrime := num.ProbablyPrime(run)
	m.IsPrime=isPrime
	//fmt.Println(isPrime)
	return m
}
