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
	"strings"
	"time"

	"github.com/google/uuid"
)

type Block struct {
	Transactions []Transaction
	Timestamp    int64
	PreviousHash string
	Nonce        int
	Hash         string
	Capacity     int
}

func (b *Block) CalculateHash() string {
	// Объединяем содержимое блока, включая транзакции
	blockContent := ""
	for _, tx := range b.Transactions {
		blockContent += tx.GetDataString()
	}

	// Создаем byte array с blockContent
	blockBytes := []byte(fmt.Sprintf("%d%s%s%d", b.Timestamp, b.PreviousHash, blockContent, b.Nonce))

	// Получаем хэш блока с использованием SHA256
	hash := sha256.Sum256(blockBytes)
	// Кодируем хэш в строку
	hashHex := hex.EncodeToString(hash[:])

	// Возвращаем хэш в виде строки
	return hashHex
}

func (b *Block) MineBlock(difficulty int) {
	for {
		// Вычисляем хэш блока
		hash := b.CalculateHash()

		// Берем первые несколько битов (размера difficulty) из хэша
		prefix := strings.Repeat("0", difficulty)
		if strings.HasPrefix(hash, prefix) {
			// Nonce найден, блок майнится
			b.Hash = hash
			return
		} else {
			// Увеличиваем Nonce и пробуем снова
			b.Nonce++
		}
	}
}

func (b *Block) ValidateBlock() bool {
	calculatedHash := b.CalculateHash()

	if len(b.Transactions) > b.Capacity-1 {
		return false
	}

	for _, tx := range b.Transactions {
		if !tx.IsValid() {
			return false
		}
	}

	return b.Hash == calculatedHash
}

type Transaction struct {
	FromAddress   string
	ToAddress     string
	Amount        float64
	Timestamp     int
	TransactionId string
	Signature     string
	Fee           float64
}

func (t *Transaction) GetDataString() string {
	return t.FromAddress +
		t.ToAddress +
		fmt.Sprintf("%.2f", t.Amount) +
		strconv.Itoa(t.Timestamp) +
		t.TransactionId +
		fmt.Sprintf("%.2f", t.Fee)
}

func (t *Transaction) calculateHash() string {
	hash := sha256.Sum256([]byte(t.GetDataString()))
	return hex.EncodeToString(hash[:])
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

	hashed := sha256.Sum256(hash)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])
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

	hashed := sha256.Sum256(hash)
	valid := ecdsa.Verify(publicKey, hashed[:], sigStruct.R, sigStruct.S)
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

type Blockchain struct {
	Blocks              []Block
	PendingTransactions []Transaction
	Difficulty          int
}

func (b *Blockchain) AddBlock(block Block) {
	if len(b.Blocks) != 0 {
		block.PreviousHash = b.Blocks[len(b.Blocks)-1].Hash
	}
	b.Blocks = append(b.Blocks, block)
}

func (b *Blockchain) AddTransactionToPool(t Transaction) {
	b.PendingTransactions = append(b.PendingTransactions, t)
}

func (b *Blockchain) GetBalance(address string) float64 {
	var balance float64 = 0
	for _, block := range b.Blocks {
		for _, t := range block.Transactions {
			switch address {
			case t.ToAddress:
				balance += t.Amount
			case t.FromAddress:
				balance -= t.Amount
			default:
				balance += 0
			}
		}
	}
	return balance
}

func (b Blockchain) isValid() bool {
	previousHash := ""
	for index, block := range b.Blocks {
		if index == 0 {
			previousHash = block.Hash
			continue
		}
		if block.ValidateBlock() || block.PreviousHash != previousHash {
			return false
		}
		previousHash = block.Hash
	}
	return true
}

func (b *Blockchain) CreateBlock(capacity int) (Block, error) {
	if len(b.PendingTransactions) < capacity {
		return Block{}, errors.New("Pending transactions list is less than capacity")
	}
	transactions := b.PendingTransactions[0:capacity]
	b.PendingTransactions = b.PendingTransactions[capacity:]
	block := Block{
		Transactions: transactions,
		Timestamp:    time.Now().Unix(),
		Capacity:     capacity,
	}
	return block, nil
}

func InitBlockchain(difficulty int) Blockchain {
	blockchain := Blockchain{Difficulty: difficulty}
	genesisBlock := Block{
		Timestamp: time.Now().Unix(),
	}
	genesisBlock.MineBlock(blockchain.Difficulty)
	blockchain.AddBlock(genesisBlock)
	return blockchain
}

func main() {
	w := new(Wallet)
	w.keyGen()

	fmt.Println("PrivateKey is:", w.privateKey)
	fmt.Println("PublicKey is:", w.publicKey)

	t, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0, 0.1)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	fmt.Println("Transaction: ", t)
	hash := t.calculateHash()
	fmt.Println("Hash: ", hash)
	fmt.Println("Signature is valid: ", t.IsValid())
	b := InitBlockchain(5)
	fmt.Println(b)
	fmt.Println(b.isValid())
	b.AddTransactionToPool(Transaction{
		Timestamp: int(time.Now().Unix()),
	})
	b.AddTransactionToPool(Transaction{
		Timestamp: int(time.Now().Unix()),
	})
	block, err := b.CreateBlock(2)
	if err != nil {
		fmt.Println("Error while creating block:", err)
	}
	block.MineBlock(b.Difficulty)
	b.AddBlock(block)
	fmt.Println(b.isValid())
	b.AddBlock(Block{})
	fmt.Println(b.isValid())
	fmt.Println(b.GetBalance("0x123"))
}

// Utility functions
type Wallet struct {
	privateKey string
	publicKey  string
}

func (w *Wallet) keyGen() {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKeyPEMStr, err := privateKeyToPEMString(privateKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}
	w.privateKey = privateKeyPEMStr

	publicKeyPEMStr, err := publicKeyToPEMString(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error converting private key to PEM:", err)
		return
	}
	w.publicKey = publicKeyPEMStr
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
