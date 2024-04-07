package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Block struct {
	Transactions []Transaction
	Timestamp    int64
	PreviousHash string
	Nonce        int
	Hash         string
	Capacity     int
}

func (b *Block) CalculateHash() (string, error) {
	// Объединяем содержимое блока, включая транзакции
	blockContent := ""
	for _, tx := range b.Transactions {
		blockContent += tx.Data
	}

	// Создаем byte array с blockContent
	blockBytes := []byte(fmt.Sprintf("%d%s%s%d", b.Timestamp, b.PreviousHash, blockContent, b.Nonce))

	// Получаем хэш блока с использованием SHA256
	hash := sha256.Sum256(blockBytes)
	// Кодируем хэш в строку
	hashHex := hex.EncodeToString(hash[:])

	// Возвращаем хэш в виде строки
	return hashHex, nil
}

func (b *Block) MineBlock(difficulty int) {
	for {
		// Вычисляем хэш блока
		hash, _ := b.CalculateHash()

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
