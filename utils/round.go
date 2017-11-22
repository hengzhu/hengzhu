package utils

import (
	"math/rand"
	"time"
	"encoding/binary"
	"bytes"
)

func Range(m int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(m)
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
