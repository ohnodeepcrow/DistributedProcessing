package main

import (
	"math/big"
	"math/rand"
	_ "fmt"
)


var effort int
func generateCandidate() *big.Int{
	tmp := rand.Int63()
	if tmp % 2 == 0{
		tmp -= 1
	}
	return big.NewInt(tmp)
}

func setEffort(i int){
	effort=i
}
func testPrime(num big.Int) metric{
	var m metric
	m.Perf= effort
	isPrime := num.ProbablyPrime(effort)
	m.IsPrime=isPrime
	//fmt.Println(isPrime)
	return m
}