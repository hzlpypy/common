package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

var letterRunes = []rune("abcdefghijkmnpqrstuvwxyzABCDEFGHJKMNOPQRSTUVWXYZ23456789")

func NewUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func NewMD5(str string) string {
	if len(str) == 0 {
		return ""
	}
	sign := md5.New()
	sign.Write([]byte(str))
	strSign := sign.Sum(nil)
	strMD5 := hex.EncodeToString(strSign)
	return strMD5
}

//use database => created
func NewUnixtime() uint {
	return uint(time.Now().Unix())
}

//rand number
func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
