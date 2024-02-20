package helper

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	userID := int64(1)
	expiration := time.Duration(1) * time.Hour

	token, err := CreateToken(userID, expiration)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if _, err := ValidateToken(token); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token, got an empty token")
	}
}

func TestValidateToken(t *testing.T) {
	// #TC 1 Expiry Token
	userID := int64(1)
	expiration := time.Microsecond
	tokenString, err := CreateToken(userID, expiration)
	if err != nil {
		t.Fatalf("Error creating token for testing: %v", err)
	}

	time.Sleep(time.Microsecond * 2)

	claims, err := ValidateToken(tokenString)

	if err == nil || !strings.Contains(err.Error(), "token is expired") {
		t.Errorf("Expected 'token has expired' error, got: %v", err)
	}

	if claims != nil {
		t.Errorf("Expected nil claims, got: %v", claims)
	}

	// #TC 2 Create a new token with a longer expiration time for testing
	expiration = time.Hour
	tokenString, err = CreateToken(userID, expiration)
	if err != nil {
		t.Fatalf("Error creating token for testing: %v", err)
	}

	claims, err = ValidateToken(tokenString)

	if err != nil {
		t.Errorf("Unexpected error validating token: %v", err)
	}

	if claims == nil {
		t.Error("Expected non-nil claims, got nil")
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
}

func TestInitializeKeys(t *testing.T) {
	// Clear keys before initializing
	privateKey = nil
	publicKey = nil

	// Initialize keys
	initializeKeys()

	// Check if privateKey and publicKey are initialized
	if privateKey == nil || publicKey == nil {
		t.Error("Expected privateKey and publicKey to be initialized, but one or both are nil")
	}

	// Encode keys to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	public, _ := x509.MarshalPKIXPublicKey(publicKey)

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: public,
	})

	// Set environment variables with base64-encoded keys
	os.Setenv("JWT_PRIVATE_KEY", base64.StdEncoding.EncodeToString(privateKeyPEM))
	os.Setenv("JWT_PUBLIC_KEY", base64.StdEncoding.EncodeToString(publicKeyPEM))

	// Clear keys again
	privateKey = nil
	publicKey = nil

	// Re-initialize keys
	initializeKeys()

	// Check if privateKey and publicKey are initialized after setting environment variables
	if privateKey == nil || publicKey == nil {
		t.Error("Expected privateKey and publicKey to be initialized after setting environment variables, but one or both are nil")
	}
}
