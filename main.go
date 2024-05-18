package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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

	// go func() {
	// 	stdReader := bufio.NewReader(os.Stdin)
	// 	for {
	// 		fmt.Print("Enter message to broadcast: ")
	// 		msg, err := stdReader.ReadString('\n')
	// 		if err != nil {
	// 			fmt.Println("Error reading from stdin:", err)
	// 			return
	// 		}
	// 		node.BroadcastMessage(strings.TrimSpace(msg))
	// 	}
	// }()
	go func() {
		stdReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message to broadcast (transaction/block): ")
			msg, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				return
			}
			msg = strings.TrimSpace(msg)
			if msg == "transaction" {
				tx := Transaction{
					FromAddress:   "Alice",
					ToAddress:     "Bob",
					Amount:        5.00,
					Timestamp:     int(time.Now().Unix()),
					TransactionId: uuid.New().String(),
				}
				node.BroadcastTransaction(tx)
			} else if msg == "block" {
				block := Block{
					Transactions: nil,
					Timestamp:    time.Now().Unix(),
					Capacity:    5,
					PreviousHash: "previousHash",
				}
				node.BroadcastBlock(block)
			} else {
				fmt.Println("Unknown message type")
			}
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

	t1, err := NewTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t1)

	t2, err := NewTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t2)

	t3, err := NewTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t3)

	t4, err := NewTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	chain.AddTransactionToPool(t4)

	t5, err := NewTransaction(w.privateKey, w.publicKey, "0x123", 5.0)
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
