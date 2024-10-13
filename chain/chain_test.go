package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestBlock_CalculateHash(t *testing.T) {
	type fields struct {
		Transactions []Transaction
		Timestamp    int64
		PreviousHash string
		Nonce        int
		Hash         string
		Capacity     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty block",
			fields: fields{
				Transactions: []Transaction{},
				Timestamp:    1643723400,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			want: "fc5ff35565676ff04776ffdf7b19fee278ee03df89dcdf1f58f7f2557656cd69",
		},
		{
			name: "block with one transaction",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723400,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			want: "0c35717686ad9c1f665b2baf5c81044405ccea62b7b7e28cee554dc0a2c3f3c8",
		},
		{
			name: "block with multiple transactions",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
					{
						FromAddress:   "456",
						ToAddress:     "789",
						Amount:        20.0,
						TransactionId: "d6f1c4e7-9d7e-11eb-a8b3-0242ac130004",
						Timestamp:     1643723401,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723402,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			want: "c4bfabe66b3273f529920bde233a1ba280572f77db876f954c0fa2484415ad1f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				Transactions: tt.fields.Transactions,
				Timestamp:    tt.fields.Timestamp,
				PreviousHash: tt.fields.PreviousHash,
				Nonce:        tt.fields.Nonce,
				Hash:         tt.fields.Hash,
				Capacity:     tt.fields.Capacity,
			}
			if got := b.CalculateHash(); got != tt.want {
				t.Errorf("Block.CalculateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_MineBlock(t *testing.T) {
	type fields struct {
		Transactions []Transaction
		Timestamp    int64
		PreviousHash string
		Nonce        int
		Hash         string
		Capacity     int
	}
	type args struct {
		difficulty int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "simple mining",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723402,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			args: args{
				difficulty: 2,
			},
			wantErr: false,
		},
		{
			name: "mining with multiple transactions",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
					{
						FromAddress:   "456",
						ToAddress:     "789",
						Amount:        20.0,
						TransactionId: "d6f1c4e7-9d7e-11eb-a8b3-0242ac130004",
						Timestamp:     1643723401,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723402,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			args: args{
				difficulty: 3,
			},
			wantErr: false,
		},
		{
			name: "mining with high difficulty",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723402,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			args: args{
				difficulty: 6,
			},
			wantErr: false,
		},
		{
			name: "mining with invalid difficulty",
			fields: fields{
				Transactions: []Transaction{
					{
						FromAddress:   "123",
						ToAddress:     "456",
						Amount:        10.0,
						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
						Timestamp:     1643723400,
						Signature:     "sig",
					},
				},
				Timestamp:    1643723402,
				PreviousHash: "",
				Nonce:        0,
				Hash:         "",
				Capacity:     10,
			},
			args: args{
				difficulty: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				Transactions: tt.fields.Transactions,
				Timestamp:    tt.fields.Timestamp,
				PreviousHash: tt.fields.PreviousHash,
				Nonce:        tt.fields.Nonce,
				Hash:         tt.fields.Hash,
				Capacity:     tt.fields.Capacity,
			}
			err := b.MineBlock(tt.args.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("Block.MineBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_Sign(t *testing.T) {
	type fields struct {
		FromAddress   string
		ToAddress     string
		Amount        float64
		Timestamp     int
		TransactionId string
		Signature     string
	}
	type args struct {
		PrivateKeyPEMStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid signature",
			fields: fields{
				FromAddress:   "sender_address",
				ToAddress:     "recipient_address",
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "valid_transaction_id",
			},
			args: args{
				PrivateKeyPEMStr: func() string {
					privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
					privateKeyBytes, _ := x509.MarshalECPrivateKey(privateKey)
					privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})
					return string(privateKeyPEM)
				}(),
			},
			wantErr: false,
		},
		{
			name: "Invalid private key",
			fields: fields{
				FromAddress:   "sender_address",
				ToAddress:     "recipient_address",
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "invalid_key_transaction_id",
			},
			args: args{
				PrivateKeyPEMStr: "invalid_private_key",
			},
			wantErr: true,
		},
		{
			name: "Empty private key",
			fields: fields{
				FromAddress:   "sender_address",
				ToAddress:     "recipient_address",
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "empty_key_transaction_id",
			},
			args: args{
				PrivateKeyPEMStr: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transaction{
				FromAddress:   tt.fields.FromAddress,
				ToAddress:     tt.fields.ToAddress,
				Amount:        tt.fields.Amount,
				Timestamp:     tt.fields.Timestamp,
				TransactionId: tt.fields.TransactionId,
				Signature:     tt.fields.Signature,
			}
			if err := tr.Sign(tt.args.PrivateKeyPEMStr); (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_IsValid(t *testing.T) {
	// Generate a valid private key for testing
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKeyBytes, _ := x509.MarshalECPrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})
	privateKeyStr := string(privateKeyPEM)

	// Generate a valid public key for testing
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	publicKeyStr := string(publicKeyPEM)

	tests := []struct {
		name        string
		transaction Transaction
		wantValid   bool
	}{
		{
			name: "Valid transaction",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "valid_transaction_id",
				}
				tr.Sign(privateKeyStr)
				return tr
			}(),
			wantValid: true,
		},
		{
			name: "Invalid signature",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "invalid_signature_id",
					Signature:     "invalid_signature",
				}
				return tr
			}(),
			wantValid: false,
		},
		{
			name: "Missing FromAddress",
			transaction: Transaction{
				ToAddress:     "recipient_address",
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "missing_from_address_id",
				Signature:     "some_signature",
			},
			wantValid: true,
		},
		{
			name: "Missing ToAddress",
			transaction: Transaction{
				FromAddress:   publicKeyStr,
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "missing_to_address_id",
				Signature:     "some_signature",
			},
			wantValid: false,
		},
		{
			name: "Negative Amount",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        -100.0,
					Timestamp:     1643723400,
					TransactionId: "negative_amount_id",
				}
				tr.Sign(privateKeyStr)
				return tr
			}(),
			wantValid: true,
		},
		{
			name: "Zero Amount",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        0.0,
					Timestamp:     1643723400,
					TransactionId: "zero_amount_id",
				}
				tr.Sign(privateKeyStr)
				return tr
			}(),
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.transaction.IsValid(); got != tt.wantValid {
				t.Errorf("Transaction.IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

func TestTransaction_verifySignature(t *testing.T) {
	// Generate a valid key pair for testing
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKeyBytes, _ := x509.MarshalECPrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})
	privateKeyStr := string(privateKeyPEM)

	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	publicKeyStr := string(publicKeyPEM)

	// Generate another key pair for testing invalid signatures
	anotherPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	anotherPrivateKeyBytes, _ := x509.MarshalECPrivateKey(anotherPrivateKey)
	anotherPrivateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: anotherPrivateKeyBytes})
	anotherPrivateKeyStr := string(anotherPrivateKeyPEM)

	tests := []struct {
		name        string
		transaction Transaction
		wantValid   bool
		wantErr     bool
	}{
		{
			name: "Valid signature",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "valid_signature_id",
				}
				tr.Sign(privateKeyStr)
				return tr
			}(),
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "Invalid signature (signed with different key)",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "invalid_signature_id",
				}
				tr.Sign(anotherPrivateKeyStr)
				return tr
			}(),
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "Empty signature",
			transaction: Transaction{
				FromAddress:   publicKeyStr,
				ToAddress:     "recipient_address",
				Amount:        100.0,
				Timestamp:     1643723400,
				TransactionId: "empty_signature_id",
				Signature:     "",
			},
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "Invalid FromAddress (not a valid PEM)",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   "not_a_valid_pem",
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "invalid_from_address_id",
				}
				tr.Sign(privateKeyStr)
				return tr
			}(),
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "Modified transaction data after signing",
			transaction: func() Transaction {
				tr := Transaction{
					FromAddress:   publicKeyStr,
					ToAddress:     "recipient_address",
					Amount:        100.0,
					Timestamp:     1643723400,
					TransactionId: "modified_data_id",
				}
				tr.Sign(privateKeyStr)
				tr.Amount = 200.0 // Modify the amount after signing
				return tr
			}(),
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, err := tt.transaction.verifySignature()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.verifySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValid != tt.wantValid {
				t.Errorf("Transaction.verifySignature() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}