package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

// EncryptPassword 密码加密
func EncryptPassword(password, salt string) string {
	d := []byte(password + salt)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func MD5(str string) string {
	d := []byte(str)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
func RandStr(n int) string {
	b := make([]rune, n)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
