package main
import (
	"fmt"
	"os"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"time"
)


func crackHash(hashToCrack string) (metric) {
	inFile, ioErr := os.Open("dict.txt")
	foundHash := false

	fmt.Println("hi")

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