package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	once       sync.Once
)

func init() {
	once.Do(initializeKeys)
}

func initializeKeys() {
	privateKeyEnv := os.Getenv("JWT_PRIVATE_KEY")
	publicKeyEnv := os.Getenv("JWT_PUBLIC_KEY")

	if privateKeyEnv != "" && publicKeyEnv != "" {
		var err error
		privateKey, publicKey, err = decodeKeys(privateKeyEnv, publicKeyEnv)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		privateKey, publicKey, err = generateKeys()
		if err != nil {
			panic(err)
		}
	}
}

func decodeKeys(privateKeyEnv, publicKeyEnv string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyEnv)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyEnv)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func generateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyGen, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKeyGen),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKeyGen.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func CreateToken(userID int64, expiresIn time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(expiresIn)

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
