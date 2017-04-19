package main

import (
	"time"
	"strconv"
	"math/big"
)

type metric struct {
	Perf int
	IsPrime bool
	hPerf time.Duration
	Val string
}

//Maps node name/ID to Reputation
type RepMetrics struct {
	CurrentMetrics map[string]Reputation
}

type Reputation struct {
	Score		int
	Count		int
	Correct 	int
}

//Scorer should take in the current reputation and the new result and update the reputation as a result
func updateReputation(repmets RepMetrics, newmet metric, node string, scorer func(nm metric, rp Reputation)) bool{
	rep, ok := repmets.CurrentMetrics[node]
	if !ok{
		return false
	}
	scorer(newmet, rep)
	return true
}

//The score for hashing is the average time it takes to generate a collision
//It doesn't use correctness currently
func hashScorer(met metric, rep Reputation){
	rep.Count += 1
	newscore := rep.Score / rep.Count
	newscore += int(met.hPerf)
	newscore = newscore/ rep.Count
	rep.Score = newscore
}

//the score for primality is the average number of correct assessments out of 100,000
//The score is score = correct/count
func primeScorer(met metric, rep Reputation){
	i, _ := strconv.ParseInt(met.Val,10,64)
	test := big.NewInt(i)
	if met.IsPrime == testPrime(*test).IsPrime{
		rep.Correct += 1
	}
	rep.Count += 1
	rep.Score = rep.Correct / rep.Count
}