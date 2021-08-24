package jwe_auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

type KeyStorageType string

var (
	StorageConfig KeyStorageType = "config"
	StorageVault  KeyStorageType = "vault"
)

type (
	Key struct {
		Value string
		Type  KeyStorageType
	}
	KeyStorage interface {
		GetPrivate() *rsa.PrivateKey
	}
)

func (k *Key) GetPrivate() (*rsa.PrivateKey, error) {
	if k.Type == StorageConfig {
		return generateRsaPrivateKeyFromPemString(k.Value)
	}
	return nil, errors.New("private key from vault not supported at the moment")
}

func generateRsaPrivateKeyFromPemString(privatePem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri, nil
}
