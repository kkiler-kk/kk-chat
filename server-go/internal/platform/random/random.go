package random

import (
	"crypto/rand"
	"math/big"
)

// Digits 返回 n 位随机数字字符串（crypto/rand）。
func Digits(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(10))
		b[i] = digits[idx.Int64()]
	}
	return string(b)
}
