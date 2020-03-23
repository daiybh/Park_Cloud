package main

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text + GetRandomSalt()))

	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(lenx int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenx; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
