package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/zicare/rgm/config"
	"golang.org/x/crypto/bcrypt"
)

type ICrypto interface {

	// Encode string
	Encode(plain string) string

	// Validates if Encode(plain) corresponds to encoded
	Compare(plain, encoded string) bool
}

type crypto struct{}

func (crypto) Encode(plain string) string {

	pepper := config.Config().GetString("pepper")
	encoded, _ := bcrypt.GenerateFromPassword([]byte([]byte(plain+pepper)), bcrypt.DefaultCost)
	return string(encoded)
}

func (crypto) Compare(plain, encoded string) bool {

	pepper := config.Config().GetString("pepper")
	if err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(plain+pepper)); err != nil {
		return false
	}
	return true
}

func Crypto() crypto {

	return crypto{}
}

func Hash(src string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	return strings.TrimRight(base64.StdEncoding.EncodeToString(h.Sum(nil)), "=")
}
