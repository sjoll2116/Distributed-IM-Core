package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword 采用 bcrypt 对明文密码进行加密处理，采用默认的工作因子(代价)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 核验给定的明文密码与存储哈希密文是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
