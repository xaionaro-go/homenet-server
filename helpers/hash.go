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
	sum := sha512.Sum512(append(salt, in...))
	return append(salt, sum[:]...)
}

func CheckHash(hash, check []byte) bool {
	salt := hash[:8]
	sum := sha512.Sum512(append(salt, check...))
	return string(sum[:]) == string(hash[8:])
}
