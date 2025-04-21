package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// 私钥生成
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	// Just as a demo, generate a new private/public key pair on each run. See note above.
	rng := rand.Reader
	var err error
	privateKey, err := rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	return privateKey, nil
}

// 验证密码
func VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
