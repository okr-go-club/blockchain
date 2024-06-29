package storage

import (
	"encoding/json"
	"fmt"

	"blockchain/chain"

	"github.com/dgraph-io/badger/v4"
)

const (
	blockPrefix       = "block_"
	transactionPrefix = "tx_"
	blockSeqKey       = "block_sequence"
	txSeqKey          = "block_sequence"
)

func keyHasPrefix(key, prefix string) bool {
	return len(key) >= len(prefix) && key[:len(prefix)] == prefix
}

func (bs *Storage) getNextSeq(txn *badger.Txn, seqKey string) (int64, error) {
	item, err := txn.Get([]byte(seqKey))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 1, nil
		}
		return 0, err
	}
	var seq int64
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &seq)
	})
	if err != nil {
		return 0, err
	}
	return seq + 1, nil
}

func (bs *Storage) setNextSeq(txn *badger.Txn, seqKey string, seq int64) error {
	seqData, err := json.Marshal(seq)
	if err != nil {
		return err
	}
	return txn.Set([]byte(seqKey), seqData)
}

type Storage struct {
	db *badger.DB
}

func NewBadgerStorage(path string) (*Storage, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (bs *Storage) Close() {
	bs.db.Close()
}

func (bs *Storage) Load(difficulty, maxBlockSize int, miningReward float64) (*chain.Blockchain, error) {
	blockchain := &chain.Blockchain{Storage: bs}

	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				key := string(item.Key())
				if keyHasPrefix(key, blockPrefix) {
					var block chain.Block
					if err := json.Unmarshal(val, &block); err != nil {
						return err
					}
					blockchain.Blocks = append(blockchain.Blocks, block)
				} else if keyHasPrefix(key, transactionPrefix) {
					var transaction chain.Transaction
					if err := json.Unmarshal(val, &transaction); err != nil {
						return err
					}
					blockchain.PendingTransactions = append(blockchain.PendingTransactions, transaction)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	blockchain.Difficulty = difficulty
	blockchain.MaxBlockSize = maxBlockSize
	blockchain.MiningReward = miningReward

	return blockchain, nil
}

func (bs *Storage) AddBlock(b chain.Block) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		seq, err := bs.getNextSeq(txn, blockSeqKey)
		if err != nil {
			return err
		}

		blockData, err := json.Marshal(b)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s%06d_%s", blockPrefix, seq, b.Hash)
		err = txn.Set([]byte(key), blockData)
		if err != nil {
			return err
		}

		return bs.setNextSeq(txn, blockSeqKey, seq)
	})

	return err
}

func (bs *Storage) AddTransaction(t chain.Transaction) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		seq, err := bs.getNextSeq(txn, txSeqKey)
		if err != nil {
			return err
		}

		txData, err := json.Marshal(t)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s%06d_%s", transactionPrefix, seq, t.TransactionId)
		err = txn.Set([]byte(key), txData)
		if err != nil {
			return err
		}

		return bs.setNextSeq(txn, txSeqKey, seq)
	})

	return err
}

func (storage *Storage) deleteByPrefix(prefix []byte) error {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := storage.db.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}

func (bs *Storage) Reset(chain *chain.Blockchain) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		err := bs.deleteByPrefix([]byte(blockPrefix))
		if err != nil {
			return err
		}
		err = bs.deleteByPrefix([]byte(transactionPrefix))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = bs.db.Update(func(txn *badger.Txn) error {
		for _, block := range chain.Blocks {
			err := bs.AddBlock(block)
			if err != nil {
				return err
			}
		}

		for _, transaction := range chain.PendingTransactions {
			err := bs.AddTransaction(transaction)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
