package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const templateArgonString = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"

type Argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewArgon2Hasher(time, memory, keyLen, saltLen uint32, threads uint8) *Argon2Hasher {
	return &Argon2Hasher{time, memory, threads, keyLen, saltLen}
}

func (ah *Argon2Hasher) HashPassword(password string) (string, error) {
	salt, err := ah.generateSalt()
	if err != nil {
		return "", errors.Wrap(err, "generate salt error")
	}
	hash := argon2.IDKey([]byte(password), salt, ah.time, ah.memory, ah.threads, ah.keyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	hashedPassword := fmt.Sprintf(templateArgonString, argon2.Version, ah.memory, ah.time, ah.threads, encodedSalt, encodedHash)

	return hashedPassword, nil
}

func (ah *Argon2Hasher) CheckPassword(password, hashedPassword string) (bool, error) {
	var version int
	var parallelism uint8
	var memory, iterations uint32
	var salt, hash string

	_, err := fmt.Sscanf(hashedPassword, templateArgonString, &version, &memory, &iterations, &parallelism, &salt, &hash)
	if err != nil {
		return false, errors.Wrap(err, "sscanf read to vars error")
	}

	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false, errors.Wrap(err, "error decode base64 salt")
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, errors.Wrap(err, "error decode base64 hash password")
	}

	newHash := argon2.IDKey([]byte(password), saltBytes, iterations, memory, parallelism, ah.keyLen)

	return subtle.ConstantTimeCompare(hashBytes, newHash) == 1, nil
}

func (ah *Argon2Hasher) generateSalt() ([]byte, error) {
	salt := make([]byte, ah.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
