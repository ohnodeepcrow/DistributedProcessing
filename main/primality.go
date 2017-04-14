package main

import (
	"math/big"
	"math/rand"
)

func generateCandidate() *big.Int{
	tmp := rand.Int63()
	if tmp % 2 == 0{
		tmp -= 1
	}
	return big.NewInt(tmp)
}

func testPrime(){

}