package main
import (
	"fmt"
	"os"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"time"
	"math/rand"
)
var run int
var file string


func crackHash(hashToCrack string) (metric) {

	foundHash := false



	if run==4{
		file ="dict1.txt"
	} else if run==1{
		file ="dict2.txt"
	}else if run==2{
		file ="dict3.txt"
	}else if run==3{
		file ="dict4.txt"
	}else{
		file="dict.txt"
	}
	fmt.Println(file)

	inFile, ioErr := os.Open(file)

	if ioErr != nil{
		fmt.Println(ioErr)
		var dummy metric
		dummy.hPerf=time.Millisecond*0
		dummy.Hash=""
		return dummy
	}

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var fileTextLine string

	tStart := time.Now()



	for scanner.Scan() {
		fileTextLine = scanner.Text()
		fileHash := getMD5HashForString(fileTextLine)

		if fileHash == hashToCrack {
			foundHash = true
			break
		}
	}

	tEnd := time.Now()

	time:=tEnd.Sub(tStart)

	var h metric
	if foundHash{
		h.hPerf=time
		h.Hash=fileTextLine
		return h
	}else{
		h.hPerf=time
		h.Hash=""
		return h
	}


}

func getMD5HashForString(userString string) string {
	hash:= md5.New()
	hash.Write([]byte(userString))

	return hex.EncodeToString(hash.Sum(nil))
}
func setDict(){
	rand.Seed(time.Now().UTC().UnixNano())
	run =rand.Intn(4-0) + 0
}

func generateHash ()string{
	hash:= md5.New()
	hash.Write([]byte(generateCandidate()))

	return hex.EncodeToString(hash.Sum(nil))
}