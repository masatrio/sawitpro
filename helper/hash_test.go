package helper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
			if tc.ExpectError {
				hasherFunc = func(password []byte, cost int) ([]byte, error) {
					return []byte{}, errors.New("sample error")
				}
			}

			hash, err := HashPassword(tc.Password)

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
