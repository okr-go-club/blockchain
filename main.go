package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MineStatusResponse struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

type MineResponse struct {
	Id string `json:"id"`
}

const (
	StatusPending    = "pending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	chain := InitBlockchain(6, 5, 5)
	miningLock := sync.Mutex{}
	statusesLock := sync.RWMutex{}
	listenAddress := flag.String("address", "localhost:8080", "Address to listen on")
	httpAddress := flag.String("http", "localhost:8090", "Address to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peers to connect to")
	miningStatuses := map[uuid.UUID]MineStatusResponse{}
	flag.Parse()

	node := NewNode(*listenAddress, strings.Split(*peers, ","))

	go node.StartServer(chain)

	for _, peer := range node.Peers {
		if peer != "" {
			go node.ConnectToPeer(peer, chain)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /blockchain/mine", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		lock := miningLock.TryLock()
		if !lock {
			fmt.Println("Already mine block")
			_, err := w.Write([]byte(fmt.Sprintf("Mining already started")))
			if err != nil {
				fmt.Println("Error while handle request", err)
			}
			return
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					statusesLock.Lock()
					miningStatuses[id] = MineStatusResponse{
						Status:  StatusFailed,
						Details: fmt.Sprintf("Panic: %v", r),
					}
					statusesLock.Unlock()
				}
				miningLock.Unlock()
			}()
			statusesLock.Lock()
			miningStatuses[id] = MineStatusResponse{Status: StatusPending}
			statusesLock.Unlock()
			chain.MinePendingTransactions("")
			statusesLock.Lock()
			miningStatuses[id] = MineStatusResponse{Status: StatusSuccessful}
			statusesLock.Unlock()
		}()
		err := json.NewEncoder(w).Encode(MineResponse{Id: id.String()})
		if err != nil {
			fmt.Println("Error while handle request", err)
		}
	})

	mux.HandleFunc("GET /blockchain/mine/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		rawId := r.PathValue("id")
		id, err := uuid.Parse(rawId)
		if err != nil {
			fmt.Println("Error while handle request", err)
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(fmt.Sprintf("Invalid id: %s", rawId)))
			if err != nil {
				fmt.Println("Error while handle request", err)
			}
			return
		}
		statusesLock.RLock()
		status := miningStatuses[id]
		statusesLock.RUnlock()
		if (status == MineStatusResponse{}) {
			http.NotFound(w, r)
			return
		}
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			fmt.Println("Error while handle request", err)
		}
	})

	server := http.Server{
		Addr:    *httpAddress,
		Handler: mux,
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
