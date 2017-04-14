package main

import "time"

func getCurrentTimestamp() string{
	return time.Now().Format("15:04:05")
}