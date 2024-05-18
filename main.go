package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	listenAddress := flag.String("address", "localhost:8080", "Address to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peers to connect to")
	flag.Parse()

	node := NewNode(*listenAddress, strings.Split(*peers, ","))

	go node.StartServer()

	for _, peer := range node.Peers {
		if peer != "" {
			go node.ConnectToPeer(peer)
		}
	}

	go func() {
		stdReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message to broadcast: ")
			msg, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				return
			}
			node.BroadcastMessage(strings.TrimSpace(msg))
		}
	}()

	select {}
}

func test() {
	w := new(Wallet)
	w.keyGen()
	fmt.Println("Successfully generate wallet keys!")
	fmt.Print("\n\n")

	chain := InitBlockchain(5, 5, 5)
	fmt.Println("Successfully initialized blockchain!")
	fmt.Println("Blockchain is valid: ", chain.IsValid())
	fmt.Print("\n\n")

	fmt.Println("Balance of 0x123 before mining:", chain.GetBalance("0x123"))
	fmt.Println("Adding transactions to the pool...")

	t1, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t1)

	t2, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t2)

	t3, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t3)

	t4, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t4)

	t5, err := createTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t5)

	fmt.Println("Length of pending transactions:", len(chain.PendingTransactions))
	fmt.Print("\n\n")

	fmt.Println("Mining...")
	chain.MinePendingTransactions("0x123")
	fmt.Println("Mining successful. New block added to the chain!")
	fmt.Println("Blockchain is valid: ", chain.IsValid())
	fmt.Print("\n\n")

	fmt.Println("Balance of 0x123 after mining:", chain.GetBalance("0x123"))
	fmt.Print("\n\n")
	
	fmt.Println("Length of pending transactions after mining:", len(chain.PendingTransactions))
	fmt.Print("\n\n")

	fmt.Println("Adding invalid block to the chain...")
	chain.AddBlock(Block{})
	fmt.Println("Blockchain is valid: ", chain.IsValid())
}
