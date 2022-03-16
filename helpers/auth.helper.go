package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"strings"
)

func HashPassword(password string) (*string, error) {
	randByte := make([]byte, 8)

	_, err := rand.Read(randByte)
	if err != nil {
		return nil, err
	}

	base64RandByte := base64.StdEncoding.EncodeToString(randByte)
	salt := []byte(base64RandByte)

	iter := 180000

	dk := pbkdf2.Key([]byte(password), salt, iter, 32, sha256.New)

	hashedPassword := fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iter, string(salt), base64.StdEncoding.EncodeToString(dk))

	return &hashedPassword, nil
}

func ComparePassword(userPassword string, password string) bool {
	splitted := strings.Split(userPassword, "$")

	salt := []byte(splitted[2])

	// saved password iteration value should be converted to int
	iter, _ := strconv.Atoi(splitted[1])

	dk := pbkdf2.Key([]byte(password), salt, iter, 32, sha256.New)

	hashedPassword := fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iter, splitted[2], base64.StdEncoding.EncodeToString(dk))

	if subtle.ConstantTimeCompare([]byte(userPassword), []byte(hashedPassword)) == 0 {
		return false
	}

	return true
}
