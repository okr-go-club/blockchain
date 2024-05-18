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
	Transactions []Transaction `json:"transactions"`
	Timestamp    int64         `json:"timestamp"`
	PreviousHash string        `json:"previousHash"`
	Nonce        int           `json:"nonce"`
	Hash         string        `json:"hash"`
	Capacity     int           `json:"capacity"`
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

func (b *Block) IsValid() bool {
	calculatedHash := b.CalculateHash()

	if len(b.Transactions) > b.Capacity {
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
	FromAddress   string  `json:"fromAddress"`
	ToAddress     string  `json:"toAddress"`
	Amount        float64 `json:"amount"`
	Timestamp     int     `json:"timestamp"`
	TransactionId string  `json:"transactionId"`
	Signature     string  `json:"signature"`
}

func (t *Transaction) GetDataString() string {
	return t.FromAddress +
		t.ToAddress +
		fmt.Sprintf("%.2f", t.Amount) +
		strconv.Itoa(t.Timestamp) +
		t.TransactionId
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

func NewTransaction(privateKey, fromAddress, toAddress string, amount float64) (Transaction, error) {
	t := Transaction{
		FromAddress:   fromAddress,
		ToAddress:     toAddress,
		Amount:        amount,
		Timestamp:     int(time.Now().Unix()),
		TransactionId: uuid.New().String(),
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
	MaxBlockSize        int
	MiningReward        float64
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

func (b Blockchain) IsValid() bool {
	previousHash := ""
	for index, block := range b.Blocks {
		if !block.IsValid() {
			return false
		}
		if index == 0 {
			previousHash = block.Hash
			continue
		}
		if block.PreviousHash != previousHash {
			return false
		}
		previousHash = block.Hash
	}
	return true
}

func (chain *Blockchain) MinePendingTransactions(minerAddress string) {
	currentPoolSize := len(chain.PendingTransactions)
	var transactions []Transaction

	if currentPoolSize < chain.MaxBlockSize {
		transactions = chain.PendingTransactions[0:currentPoolSize]
		chain.PendingTransactions = chain.PendingTransactions[currentPoolSize:]
	} else {
		transactions = chain.PendingTransactions[0 : chain.MaxBlockSize-1]
		chain.PendingTransactions = chain.PendingTransactions[chain.MaxBlockSize-1:]
	}

	rewardTx := Transaction{
		FromAddress:   "",
		ToAddress:     minerAddress,
		Amount:        chain.MiningReward,
		Timestamp:     int(time.Now().Unix()),
		TransactionId: uuid.New().String(),
	}
	transactions = append(transactions, rewardTx)

	block := Block{
		Transactions: transactions,
		Timestamp:    time.Now().Unix(),
		Capacity:     chain.MaxBlockSize,
		PreviousHash: chain.Blocks[len(chain.Blocks)-1].Hash,
	}
	block.Hash = block.CalculateHash()

	block.MineBlock(chain.Difficulty)
	chain.AddBlock(block)
}

func InitBlockchain(difficulty, maxBlockSize int, miningReward float64) Blockchain {
	blockchain := Blockchain{Difficulty: difficulty, MaxBlockSize: maxBlockSize, MiningReward: miningReward}
	genesisBlock := Block{
		Timestamp: time.Now().Unix(),
	}
	genesisBlock.MineBlock(blockchain.Difficulty)
	blockchain.AddBlock(genesisBlock)
	return blockchain
}

type Wallet struct {
	privateKey string
	publicKey  string
}

func (w *Wallet) KeyGen() {
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
