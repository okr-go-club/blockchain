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

// func TestBlock_IsValid(t *testing.T) {
// 	type fields struct {
// 		Transactions []Transaction
// 		Timestamp    int64
// 		PreviousHash string
// 		Nonce        int
// 		Hash         string
// 		Capacity     int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   bool
// 	}{
// 		{
// 			name: "empty block is invalid",
// 			fields: fields{
// 				Transactions: []Transaction{},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "block with transactions is valid",
// 			fields: fields{
// 				Transactions: []Transaction{
// 					{
// 						FromAddress:   "123",
// 						ToAddress:     "456",
// 						Amount:        10.0,
// 						TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
// 						Timestamp:     1643723400,
// 						Signature:     b.Transactions[0].generateVerifiedSignature(),
// 					},
// 				},
// 			},
// 			want: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := &Block{
// 				Transactions: tt.fields.Transactions,
// 				Timestamp:    tt.fields.Timestamp,
// 				PreviousHash: tt.fields.PreviousHash,
// 				Nonce:        tt.fields.Nonce,
// 				Hash:         tt.fields.Hash,
// 				Capacity:     tt.fields.Capacity,
// 			}
// 			if got := b.IsValid(); got != tt.want {
// 				t.Errorf("Block.IsValid() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestTransaction_verifySignature(t *testing.T) {
	// Generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Marshal the private key to PEM
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509.MarshalPKCS8PrivateKey(privateKey)},
	)

	type fields struct {
		FromAddress   string
		ToAddress     string
		Amount        float64
		Timestamp     int
		TransactionId string
		Signature     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "valid signature",
			fields: fields{
				FromAddress:   "123",
				ToAddress:     "456",
				Amount:        10.0,
				Timestamp:     1643723400,
				TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
				Signature:     "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid signature",
			fields: fields{
				FromAddress:   "123",
				ToAddress:     "456",
				Amount:        10.0,
				Timestamp:     1643723400,
				TransactionId: "d6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
				Signature:     "wrong signature",
			},
			want:    false,
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
			}

			// Sign the transaction using the private key
			err = tr.Sign(string(privateKeyPEM))
			if err != nil {
				t.Fatal(err)
			}

			// Verify the signature
			got, err := tr.VerifySignature()
			if err != nil {
				t.Errorf("VerifySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
