package api

import (
	"blockchain/chain"
	"blockchain/p2p"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

func SetCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, req)
	})
}

type Handler struct {
	Blockchain       *chain.Blockchain
	Node             *p2p.Node
	BlockchainRWLock sync.RWMutex
	MiningLock       sync.Mutex
	StatusesRWLock   sync.RWMutex
	MiningStatuses   map[uuid.UUID]MineStatusResponse
}

type MineResponse struct {
	Id string `json:"id"`
}

type MineStatusResponse struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

type AddTransactionRequest struct {
	PrivateKey string  `json:"privateKey"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount,string"`
}

const (
	StatusPending    = "pending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

// @Success 200 {object} api.MineResponse
// @Router /blockchain/mine [post]
func (h *Handler) MineBlock(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	lock := h.MiningLock.TryLock()
	if !lock {
		fmt.Println("Already mine block")
		_, err := w.Write([]byte("Mining already started"))
		if err != nil {
			fmt.Println("Error while handle request", err)
		}
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.StatusesRWLock.Lock()
				h.MiningStatuses[id] = MineStatusResponse{
					Status:  StatusFailed,
					Details: fmt.Sprintf("Panic: %v", r),
				}
				h.StatusesRWLock.Unlock()
			}
			h.MiningLock.Unlock()
		}()
		h.StatusesRWLock.Lock()
		h.MiningStatuses[id] = MineStatusResponse{Status: StatusPending}
		h.StatusesRWLock.Unlock()
		err := h.Blockchain.MinePendingTransactions("")
		if err != nil {
			h.StatusesRWLock.Lock()
			h.MiningStatuses[id] = MineStatusResponse{
				Status:  StatusFailed,
				Details: fmt.Sprintf("Error: %v", err),
			}
			h.StatusesRWLock.Unlock()
		}
		h.StatusesRWLock.Lock()
		h.MiningStatuses[id] = MineStatusResponse{Status: StatusSuccessful}
		h.StatusesRWLock.Unlock()
	}()
	err := json.NewEncoder(w).Encode(MineResponse{Id: id.String()})
	if err != nil {
		fmt.Println("Error while handle request", err)
	}
}

// @Param id path string true "Mining process ID"
// @Success 200 {array} api.MineStatusResponse
// @Router /blockchain/mine/{id} [get]
func (h *Handler) GetMiningStatus(w http.ResponseWriter, r *http.Request) {
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
	h.StatusesRWLock.RLock()
	status := h.MiningStatuses[id]
	h.StatusesRWLock.RUnlock()
	if (status == MineStatusResponse{}) {
		http.NotFound(w, r)
		return
	}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		fmt.Println("Error while handle request", err)
	}
}

// @Success 200 {array} chain.Transaction
// @Router /blocks/pool/ [get]
func (h *Handler) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	h.BlockchainRWLock.RLock()
	transactions := h.Blockchain.PendingTransactions
	jsonTransactions, _ := json.Marshal(transactions)
	h.BlockchainRWLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonTransactions)

	if err != nil {
		fmt.Println("Error during writing response:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

// @Param request body api.AddTransactionRequest true "query params"
// @Success 200
// @Router /transactions [post]
func (h *Handler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	var request AddTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("Error while handle request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	privateKey := "-----BEGIN EC PRIVATE KEY-----\n" + request.PrivateKey + "\n-----END EC PRIVATE KEY-----"

	transaction, err := chain.NewTransaction(privateKey, request.From, request.To, request.Amount)
	if err != nil {
		fmt.Println("Error while create transaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.BlockchainRWLock.RLock()
	err = h.Blockchain.AddTransactionToPool(transaction)
	h.BlockchainRWLock.RUnlock()
	if err != nil {
		fmt.Println("Error while add transaction to pool", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go h.Node.BroadcastTransaction(transaction)
	w.WriteHeader(http.StatusOK)
}

// @Success 200 {object} chain.Blockchain
// @Router /blocks/pool [get]
func (h *Handler) GetBlocksPool(w http.ResponseWriter, r *http.Request) {
	h.BlockchainRWLock.RLock()
	blocks := h.Blockchain
	jsonBlocks, _ := json.Marshal(blocks)
	h.BlockchainRWLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonBlocks)

	if err != nil {
		fmt.Println("Error during writing response:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}
