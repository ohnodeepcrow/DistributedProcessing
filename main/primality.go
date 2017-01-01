package main

import (
	"math/big"
	"math/rand"
	"fmt"
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
func testPrime(num big.Int,effort int){
	isPrime := num.ProbablyPrime(effort)
	fmt.Println(isPrime)
}