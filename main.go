package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGET(w, r)
	default:
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	var filename = "transactions_pool.json"
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error during reading file:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonFile)

	if err != nil {
		fmt.Println("Error during writing response:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func StartWebServer(port string) {
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Fatal server error:", err)
	}
}

func main() {
	// Регистрируем обработчик для роута "/"
	http.HandleFunc("/transactions/pool/", handler)

	go StartWebServer(":8888")

	chain := InitBlockchain(5, 5, 5)

	listenAddress := flag.String("address", "localhost:8080", "Address to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peers to connect to")
	flag.Parse()

	node := NewNode(*listenAddress, strings.Split(*peers, ","))

	go node.StartServer(chain)

	for _, peer := range node.Peers {
		if peer != "" {
			go node.ConnectToPeer(peer, chain)
		}
	}

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
			switch msg {
			case "transaction":
				tx := Transaction{
					FromAddress:   "Alice",
					ToAddress:     "Bob",
					Amount:        5.00,
					Timestamp:     int(time.Now().Unix()),
					TransactionId: uuid.New().String(),
				}
				node.BroadcastTransaction(tx)
			case "block":
				block := Block{
					Transactions: nil,
					Timestamp:    time.Now().Unix(),
					Capacity:     5,
					PreviousHash: "previousHash",
				}
				node.BroadcastBlock(block)
			default:
				fmt.Println("Unknown message type")
			}
		}
	}()

	select {}
}

//nolint:all
func test() {
	w := new(Wallet)
	w.KeyGen()
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
