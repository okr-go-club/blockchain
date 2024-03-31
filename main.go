package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
)

/*
Transaction
Properties: FromAddress, ToAddress, Amount, Timestamp, TransactionId, Signature, Fee

Behavior:
	+ Calculates its own hash based on transaction content.
	+ Signs itself with the sender's private key.
	+ Validates itself (signature and data integrity).
*/

type Transaction struct {
	FromAddress string
	ToAddress string
	Amount float64
	Timestamp int
	TransactionId string
	Signature string
	Fee float64
}

func (t *Transaction) CalculateHash() string {
	record := t.FromAddress +
		t.ToAddress +
		fmt.Sprintf("%.2f", t.Amount) +
		strconv.Itoa(t.Timestamp) +
		t.TransactionId +
		fmt.Sprintf("%.2f", t.Fee)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (t *Transaction) IsValid() bool {
	if t.FromAddress == "" { return true }  // Transaction from the system
	if t.Signature == "" {
		fmt.Println("Transaction is not signed")
		return false
	}

	isValid, err := t.verifySignature()
	if err != nil {
		fmt.Println("Error verifying signature:", err)
		return false
	}

	return isValid
}

func (t *Transaction) Sign(privateKeyPEMStr string) error {
    block, _ := pem.Decode([]byte(privateKeyPEMStr))
    if block == nil {
        return errors.New("failed to parse PEM block containing the key")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return err
    }

    hasher := sha256.New()
    hasher.Write([]byte(t.CalculateHash()))
    hash := hasher.Sum(nil)

    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
    if err != nil {
        return err
    }

    t.Signature = base64.StdEncoding.EncodeToString(signature)
    return nil
}

func (t *Transaction) verifySignature() (bool, error) {
    block, _ := pem.Decode([]byte(t.FromAddress))
    if block == nil {
        return false, errors.New("failed to parse PEM block containing the key")
    }

    publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return false, err
    }

    publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
    if !ok {
        return false, errors.New("not RSA public key")
    }

    signature, err := base64.StdEncoding.DecodeString(t.Signature)
    if err != nil {
        return false, err
    }

    hasher := sha256.New()
    hasher.Write([]byte(t.CalculateHash()))
    hash := hasher.Sum(nil)

    err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signature)
    if err != nil {
        return false, err
    }

    return true, nil
}

func main() {
	privateKeyPEMStr, publicKeyPEMStr, err := GenerateRSAKeys(4096)
	if err != nil {
		fmt.Println("Error generating RSA keys:", err)
		return
	}
	t := Transaction{
		FromAddress: publicKeyPEMStr,
		ToAddress: "0x456",
		Amount: 10.0,
		Timestamp: 1234567890,
		TransactionId: "abc123",
		Fee: 0.1,
	}
    fmt.Println("Transaction: ", t)
	hash := t.CalculateHash()
    fmt.Println("Hash: ", hash)

	t.Sign(privateKeyPEMStr)
	isValid := t.IsValid()
	fmt.Println("Signature is valid: ", isValid)
}

// Utility functions
func GenerateRSAKeys(bits int) (privateKeyPEMStr, publicKeyPEMStr string, err error) {
    privatekey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return "", "", err
    }

    // convert to DER format for PEM encoding
    privateDER := x509.MarshalPKCS1PrivateKey(privatekey)

    // Create a PEM block for the private key
    privateKeyBlock := &pem.Block{
        Type:  "RSA Private Key",
        Bytes: privateDER,
    }

    // Encode the private key to PEM format and convert to string
    privateKeyPEMStr = string(pem.EncodeToMemory(privateKeyBlock))

    // Extract the public key from the private key, marshal to DER format
    publicDER, err := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
    if err != nil {
        return "", "", err
    }

    // Create a PEM block for the public key
    publicKeyBlock := &pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicDER,
    }

    // Encode the public key to PEM format and convert to string
    publicKeyPEMStr = string(pem.EncodeToMemory(publicKeyBlock))

    return privateKeyPEMStr, publicKeyPEMStr, nil
}
