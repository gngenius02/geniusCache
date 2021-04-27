package main

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func timeStamp() string {
	return time.Now().Format("15:04:05.000")
}

func getRandBytes() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
