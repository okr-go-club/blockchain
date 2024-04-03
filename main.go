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
	pem, _ := pem.Decode([]byte(privateKeyPEMStr))
	if pem == nil {
		return errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes)
	if err != nil {
		return err
	}

	hash, err := hex.DecodeString(t.calculateHash())
	if err != nil {
		return err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return err
	}

	t.Signature = base64.StdEncoding.EncodeToString(signature)
	return nil
}

func (t *Transaction) verifySignature() (bool, error) {
	pem, _ := pem.Decode([]byte(t.FromAddress))
	if pem == nil {
		return false, errors.New("failed to parse PEM block containing the key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(pem.Bytes)
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

	hash, err := hex.DecodeString(t.calculateHash())
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signature)
	if err != nil {
		return false, err
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

type Blockchain struct {
	transactions []Transaction
}

func (b *Blockchain) addTransaction(t Transaction) {
	b.transactions = append(b.transactions, t)
}

func (b *Blockchain) getBalance(address string) float64 {
	var balance float64 = 0
	for _, t := range b.transactions {
		switch address {
		case t.ToAddress:
			balance += t.Amount
		case t.FromAddress:
			balance -= t.Amount
		default:
			balance += 0
		}
	}
	return balance
}

func main() {
	privateKeyPEMStr, publicKeyPEMStr, err := GenerateRSAKeys(4096)
	if err != nil {
		fmt.Println("Error generating RSA keys:", err)
		return
	}

	t, err := createTransaction(privateKeyPEMStr, publicKeyPEMStr, "0x123", 5.0, 0.1)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	fmt.Println("Transaction: ", t)
	hash := t.calculateHash()
	fmt.Println("Hash: ", hash)
	fmt.Println("Signature is valid: ", t.IsValid())
	b := Blockchain{}
	b.addTransaction(t)
	fmt.Println(b)
	fmt.Println(b.getBalance("0x123"))
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
