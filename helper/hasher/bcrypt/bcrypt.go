package bcrypt

import (
	"github.com/sawitpro/UserService/helper/hasher"
	"golang.org/x/crypto/bcrypt"
)

type bcrypthash struct {
	hashFunc func(password []byte, cost int) ([]byte, error)
}

func NewHasher() hasher.PasswordHasher {
	return &bcrypthash{
		hashFunc: bcrypt.GenerateFromPassword,
	}
}

func (h *bcrypthash) HashPassword(password string) (string, error) {
	hash, err := h.hashFunc([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *bcrypthash) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
