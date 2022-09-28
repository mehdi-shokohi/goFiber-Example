package goexJWT

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"

	"goex/db/redisHelper"
)

var (
	privateKey *rsa.PrivateKey
)

func New() func(*fiber.Ctx) error {
	if privateKey == nil {
		GetPrivateKey()
	}
	return jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    privateKey.Public(),
	})
}

func GetPrivateKey() (*rsa.PrivateKey, error) {

	if privateKey != nil {
		return privateKey, nil
	}
	rv := redisHelper.GetValue(context.Background(), "private_key")
	if rv.Val() != "" {
		var err error
		privateKey, err = ParseRsaPrivateKeyFromPemStr(rv.Val())
		if err == nil {

			return privateKey, nil
		}
	}
	rng := rand.Reader
	var err error
	privateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	pv := ExportRsaPrivateKeyAsPemStr(privateKey)
	fmt.Println(pv)

	redisHelper.SaveKey(context.Background(), "private_key", pv, time.Second*999999)

	return privateKey, nil
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
