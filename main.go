package main

import (
	"blockchain/api"
	"blockchain/chain"
	"blockchain/p2p"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	blockchain := chain.InitBlockchain(5, 5, 5)
	listenAddress := flag.String("address", "localhost:8080", "Address to listen on")
	httpAddress := flag.String("http", "localhost:8090", "Address to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peers to connect to")
	flag.Parse()

	node := p2p.NewNode(*listenAddress, strings.Split(*peers, ","))
	handler := api.Handler{
		Blockchain:     blockchain,
		Node:           node,
		MiningStatuses: make(map[uuid.UUID]api.MineStatusResponse),
	}
	go node.StartServer(blockchain)

	for _, peer := range node.Peers {
		if peer != "" {
			go node.ConnectToPeer(peer, blockchain)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /blockchain/mine", handler.MineBlock)

	mux.HandleFunc("GET /blockchain/mine/{id}/status", handler.GetMiningStatus)

	mux.HandleFunc("GET /transactions/pool/", handler.GetTransactionPool)

	mux.HandleFunc("GET /blocks/pool/", handler.GetBlocksPool)

	server := http.Server{
		Addr:    *httpAddress,
		Handler: api.SetCORSHeaders(mux),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

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
				tx := chain.Transaction{
					FromAddress:   "Alice",
					ToAddress:     "Bob",
					Amount:        5.00,
					Timestamp:     int(time.Now().Unix()),
					TransactionId: uuid.New().String(),
				}
				node.BroadcastTransaction(tx)
			case "block":
				block := chain.Block{
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
	w := new(chain.Wallet)
	w.KeyGen()
	fmt.Println("Successfully generate wallet keys!")
	fmt.Print("\n\n")

	blockchain := chain.InitBlockchain(5, 5, 5)
	fmt.Println("Successfully initialized blockchain!")
	fmt.Println("Blockchain is valid: ", blockchain.IsValid())
	fmt.Print("\n\n")

	fmt.Println("Balance of 0x123 before mining:", blockchain.GetBalance("0x123"))
	fmt.Println("Adding transactions to the pool...")

	t1, err := chain.NewTransaction(w.PrivateKey, w.PublicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	blockchain.AddTransactionToPool(t1)

	t2, err := chain.NewTransaction(w.PrivateKey, w.PublicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	blockchain.AddTransactionToPool(t2)

	t3, err := chain.NewTransaction(w.PrivateKey, w.PublicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	blockchain.AddTransactionToPool(t3)

	t4, err := chain.NewTransaction(w.PrivateKey, w.PublicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	blockchain.AddTransactionToPool(t4)

	t5, err := chain.NewTransaction(w.PrivateKey, w.PublicKey, "0x123", 5.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}
	blockchain.AddTransactionToPool(t5)

	fmt.Println("Length of pending transactions:", len(blockchain.PendingTransactions))
	fmt.Print("\n\n")

	fmt.Println("Mining...")
	blockchain.MinePendingTransactions("0x123")
	fmt.Println("Mining successful. New block added to the chain!")
	fmt.Println("Blockchain is valid: ", blockchain.IsValid())
	fmt.Print("\n\n")

	fmt.Println("Balance of 0x123 after mining:", blockchain.GetBalance("0x123"))
	fmt.Print("\n\n")

	fmt.Println("Length of pending transactions after mining:", len(blockchain.PendingTransactions))
	fmt.Print("\n\n")

	fmt.Println("Adding invalid block to the chain...")
	blockchain.AddBlock(chain.Block{})
	fmt.Println("Blockchain is valid: ", blockchain.IsValid())
}
