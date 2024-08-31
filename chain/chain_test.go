package chain

import (
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
