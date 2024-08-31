package main

import (
	"blockchain/chain"
	"fmt"
)

func main() {
	wallet := chain.Wallet{}
	wallet.KeyGen()
	fmt.Println(wallet.PrivateKey)
	fmt.Println(wallet.PublicKey)
	t := chain.Transaction{}
	err := t.Sign(wallet.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}
}
