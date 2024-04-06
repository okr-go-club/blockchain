package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	FromAddress   string
	ToAddress     string
	Amount        float64
	Timestamp     int
	TransactionId string
	Signature     string
	Fee           float64
}

func (t *Transaction) calculateHash() string {
	data := t.FromAddress +
		t.ToAddress +
		fmt.Sprintf("%.2f", t.Amount) +
		strconv.Itoa(t.Timestamp) +
		t.TransactionId +
		fmt.Sprintf("%.2f", t.Fee)

	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (t *Transaction) IsValid() bool {
	// Check if the transaction is from the system
	if t.FromAddress == "" {
		return true
	}

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
	pemBlock, _ := pem.Decode([]byte(privateKeyPEMStr))
	if pemBlock == nil {
		return errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		return err
	}

	hash, err := hex.DecodeString(t.calculateHash())
	if err != nil {
		return err
	}

	hasher := sha256.New()
	_, err = hasher.Write(hash)
	if err != nil {
		return err
	}
	hashed := hasher.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return err
	}

	signature, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		return err
	}

	t.Signature = base64.StdEncoding.EncodeToString(signature)
	return nil
}

func (t *Transaction) verifySignature() (bool, error) {
	pemBlock, _ := pem.Decode([]byte(t.FromAddress))
	if pemBlock == nil {
		return false, errors.New("failed to parse PEM block containing the key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return false, err
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("not ECDSA public key")
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(t.Signature)
	if err != nil {
		return false, err
	}

	var sigStruct struct {
		R, S *big.Int
	}

	_, err = asn1.Unmarshal(signatureBytes, &sigStruct)
	if err != nil {
		return false, err
	}

	hash, err := hex.DecodeString(t.calculateHash())
	if err != nil {
		return false, err
	}

	hasher := sha256.New()
	_, err = hasher.Write(hash)
	if err != nil {
		return false, err
	}
	hashed := hasher.Sum(nil)

	valid := ecdsa.Verify(publicKey, hashed, sigStruct.R, sigStruct.S)
	if !valid {
		return false, errors.New("signature verification failed")
	}

	return true, nil
}

func createTransaction(privateKey, fromAddress, toAddress string, amount, fee float64) (Transaction, error) {
	t := Transaction{
		FromAddress:   fromAddress,
		ToAddress:     toAddress,
		Amount:        amount,
		Timestamp:     int(time.Now().Unix()),
		TransactionId: uuid.New().String(),
		Fee:           fee,
	}
	err := t.Sign(privateKey)
	if err != nil {
		return Transaction{}, err
	}
	return t, nil
}

func main() {
	w := new(Wallet)
    w.keyGen()

	privateKeyPEMStr, err := privateKeyToPEMString(w.privateKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}

	publicKeyPEMStr, err := publicKeyToPEMString(w.publicKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}

    fmt.Println("PrivateKey is:", w.privateKey) 
    fmt.Println("PublicKey is:", w.publicKey) 

	t, err := createTransaction(privateKeyPEMStr, publicKeyPEMStr, "0x123", 5.0, 0.1)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	fmt.Println("Transaction: ", t)
	hash := t.calculateHash()
	fmt.Println("Hash: ", hash)
	fmt.Println("Signature is valid: ", t.IsValid())
}

// Utility functions
type Wallet struct {
	privateKey	*ecdsa.PrivateKey
	publicKey	*ecdsa.PublicKey
}

func (w *Wallet) keyGen()  {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) 
    w.privateKey = privateKey
    w.publicKey = &w.privateKey.PublicKey
}

func privateKeyToPEMString(privKey *ecdsa.PrivateKey) (string, error) {
    der, err := x509.MarshalECPrivateKey(privKey)
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

func publicKeyToPEMString(pubKey *ecdsa.PublicKey) (string, error) {
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
