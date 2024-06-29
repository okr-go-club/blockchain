package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Wallet struct {
	PrivateKey string
	PublicKey  string
}

func (w *Wallet) KeyGen() {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKeyPEMStr, err := PrivateKeyToPEMString(privateKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}
	w.PrivateKey = privateKeyPEMStr

	publicKeyPEMStr, err := PublicKeyToPEMString(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}
	w.PublicKey = publicKeyPEMStr
}

func PrivateKeyToPEMString(PrivateKey *ecdsa.PrivateKey) (string, error) {
	der, err := x509.MarshalECPrivateKey(PrivateKey)
	if err != nil {
		return "", err
	}

	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}
	pemData := pem.EncodeToMemory(pemBlock)

	return string(pemData), nil
}

func PublicKeyToPEMString(pubKey *ecdsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	pemData := pem.EncodeToMemory(pemBlock)

	return string(pemData), nil
}
