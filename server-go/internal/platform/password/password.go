package password

import "golang.org/x/crypto/bcrypt"

// Hash 使用 bcrypt(cost=10) 哈希密码。
func Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	return string(b), err
}

// Verify 校验明文与哈希是否匹配。
func Verify(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
