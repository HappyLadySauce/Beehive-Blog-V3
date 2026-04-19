package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plaintext password with bcrypt.
// HashPassword 使用 bcrypt 对明文密码做哈希。
func HashPassword(password string, cost int) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// VerifyPassword compares a plaintext password to its hash.
// VerifyPassword 对比明文密码与哈希是否匹配。
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
