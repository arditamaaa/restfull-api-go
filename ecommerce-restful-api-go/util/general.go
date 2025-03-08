package util

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
)

func StrToUint64(str string) uint64 {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func RandStringBytes(n int) string {
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func EncodeSha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	encoded := fmt.Sprintf("%x", hash.Sum(nil))
	return encoded
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
