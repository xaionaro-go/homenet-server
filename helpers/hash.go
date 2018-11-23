package helpers

import (
	"crypto/sha512"
	"math/rand"

	"github.com/Sirupsen/logrus"
)

const (
	saltSize = int(8)
)

var (
	defaultSalt = []byte("homenets")
)

func Hash(in []byte) []byte {
	salt := make([]byte, saltSize)
	n, err := rand.Read(salt)
	if n != saltSize || err != nil {
		logrus.Errorf("Cannot generate a random salt: n == %d; err == %s", n, err)
		salt = defaultSalt
	}
	key := make([]byte, saltSize+len(in))
	copy(key, salt)
	key = append(key, in...)
	sum := sha512.Sum512(key)
	hash := append(salt, sum[:]...)
	return hash
}

func CheckHash(hash, check []byte) bool {
	salt := hash[:saltSize]
	oldSum := hash[saltSize:]
	key := make([]byte, saltSize+len(check))
	copy(key, salt)
	key = append(key, check...)
	sum := sha512.Sum512(key)
	return string(sum[:]) == string(oldSum)
}
