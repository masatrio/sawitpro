package bcrypt

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		TestName    string
		Password    string
		ExpectError bool
	}{
		{"Short password", "short", false},
		{"Long password", "thisisalongpassword", false},
		{"Empty password", "", false},
		{"Empty password", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			hasher := &bcrypthash{
				hashFunc: bcrypt.GenerateFromPassword,
			}
			if tc.ExpectError {
				hasher.hashFunc = func(password []byte, cost int) ([]byte, error) {
					return []byte{}, errors.New("sample error")
				}
			}

			hash, err := hasher.HashPassword(tc.Password)

			if tc.ExpectError {
				assert.Error(t, err)
			} else {
				if err != nil {
					t.Errorf("Error hashing password: %v", err)
				}

				if len(hash) == 0 {
					t.Error("Expected non-empty hash, got an empty hash")
				}
			}

		})
	}
}

func TestNewHasher(t *testing.T) {
	hasher := NewHasher()

	bcryptHasher, ok := hasher.(*bcrypthash)
	assert.True(t, ok)

	expectedFunc := reflect.ValueOf(bcrypt.GenerateFromPassword)
	actualFunc := reflect.ValueOf(bcryptHasher.hashFunc)
	assert.Equal(t, expectedFunc.Pointer(), actualFunc.Pointer())
}

func TestCompareHashAndPassword(t *testing.T) {
	plainPassword := []byte("password123")

	hasher := NewHasher()

	hashedPassword, _ := hasher.HashPassword(string(plainPassword))

	err := hasher.CompareHashAndPassword([]byte(hashedPassword), plainPassword)

	assert.NoError(t, err)
}
