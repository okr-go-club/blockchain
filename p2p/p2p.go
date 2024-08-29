package p2p

import (
	"blockchain/chain"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Node struct {
	Address     string
	Peers       map[string]bool
	Connections map[string]net.Conn
	Mutex       sync.Mutex
}

func NewNode(address string, peers []string) *Node {
	peersMap := map[string]bool{}
	for _, peer := range peers {
		peersMap[peer] = true
	}
	return &Node{
		Address:     address,
		Peers:       peersMap,
		Connections: make(map[string]net.Conn),
	}
}

func (node *Node) StartServer(blockchain *chain.Blockchain) {
	listener, err := net.Listen("tcp", node.Address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing listener:", err)
		}
	}()

	fmt.Println("Server started on", node.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go node.HandleConnection(conn, blockchain)
	}
}

func ProcessMessage(message string, blockchain *chain.Blockchain) error {
	var msgMap map[string]interface{}
	err := json.Unmarshal([]byte(message), &msgMap)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return err
	}

	if msgType, ok := msgMap["type"]; ok {
		switch msgType {
		case "transaction":
			var tx chain.Transaction
			if txData, ok := msgMap["transaction"].(map[string]interface{}); ok {
				txJson, err := json.Marshal(txData)
				if err != nil {
					fmt.Println("Error marshalling transaction:", err)
					return err
				}

				err = json.Unmarshal(txJson, &tx)
				if err != nil {
					fmt.Println("Error unmarshalling transaction:", err)
					return err
				}
				fmt.Println("Received transaction:", tx)
				err = blockchain.AddTransactionToPool(tx)
				if err != nil {
					fmt.Println("Error adding transaction to pool:", err)
					return err
				}
				fmt.Println("Transaction pool:", blockchain.PendingTransactions)
			} else {
				fmt.Println("Invalid transaction data")
			}
		case "block":
			var block chain.Block
			if blockData, ok := msgMap["block"].(map[string]interface{}); ok {
				blockJson, err := json.Marshal(blockData)
				if err != nil {
					fmt.Println("Error marshalling block:", err)
					return err
				}
				err = json.Unmarshal(blockJson, &block)
				if err != nil {
					fmt.Println("Error unmarshalling block:", err)
					return err
				}
				fmt.Println("Received block:", block)
				blockchain.AddBlock(block)
				fmt.Println("Blocks:", blockchain.Blocks)
			} else {
				fmt.Println("Invalid block data")
			}
		default:
			fmt.Println("Unknown message type")
		}
	} else {
		fmt.Println("Message does not contain a type")
	}

	return nil
}

func (node *Node) HandleConnection(conn net.Conn, blockchain *chain.Blockchain) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()
	reader := bufio.NewReader(conn)

	// Read initial hello message
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Received Initial:", message)

	// Read peer address
	message, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Received peer address:", message)

	peerAddress := strings.TrimSpace(message)
	node.Peers[peerAddress] = true
	node.AddConnection(peerAddress, conn)

	// Get len of blockchain
	message, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	otherLenBlockchain, _ := strconv.Atoi(strings.TrimSpace(message))
	fmt.Printf(
		"Received len of blockhain: %d, peer address: %s",
		otherLenBlockchain,
		peerAddress,
	)

	// Keep the connection open to read messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			node.RemoveConnection(conn.RemoteAddr().String())
			return
		}
		fmt.Println("Received Message:", message)

		err = ProcessMessage(message, blockchain)
		if err != nil {
			fmt.Println("Error processing message:", err)
			return
		}
	}
}

func (node *Node) ConnectToPeer(address string, blockchain *chain.Blockchain) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}

	message := "Hello, Blockchain!\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Sent:", message)

	_, err = conn.Write([]byte(node.Address + "\n"))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Sent address:", node.Address)

	node.AddConnection(address, conn)
	fmt.Println("Connected to peer:", address)

	buf := new(bytes.Buffer)
	var num = int32(len(blockchain.Blocks))
	err = binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	buf.Write([]byte("\n"))
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Sent Length of my Blockchain:", len(blockchain.Blocks))

	go node.ReadData(conn, blockchain)
}

func (node *Node) AddConnection(peerAddress string, conn net.Conn) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	node.Connections[peerAddress] = conn
	fmt.Println("Connection added:", peerAddress)
}

func (node *Node) RemoveConnection(peerAddress string) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	if conn, ok := node.Connections[peerAddress]; ok {
		conn.Close()
		delete(node.Connections, peerAddress)
		node.Peers[peerAddress] = false
		fmt.Println("Connection removed:", peerAddress)
	} else {
		fmt.Println("No connection found for:", peerAddress)
	}
}

func (node *Node) ReadData(conn net.Conn, blockchain *chain.Blockchain) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			node.RemoveConnection(conn.RemoteAddr().String())
			return
		}
		fmt.Println("Received in ReadData:", message)
		err = ProcessMessage(message, blockchain)
		if err != nil {
			fmt.Println("Error processing message:", err)
			return
		}
	}
}

func (node *Node) BroadcastMessage(message string) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	for address, conn := range node.Connections {
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error writing to peer", address, ":", err)
		}
	}
}

func (node *Node) BroadcastTransaction(tx chain.Transaction) {
	txJson, err := json.Marshal(struct {
		Type        string            `json:"type"`
		Transaction chain.Transaction `json:"transaction"`
	}{"transaction", tx})
	if err != nil {
		fmt.Println("Error marshalling transaction:", err)
		return
	}
	node.BroadcastMessage(string(txJson))
}

func (node *Node) BroadcastBlock(block chain.Block) {
	blockJson, err := json.Marshal(struct {
		Type  string      `json:"type"`
		Block chain.Block `json:"block"`
	}{"block", block})
	if err != nil {
		fmt.Println("Error marshalling block:", err)
		return
	}
	node.BroadcastMessage(string(blockJson))
}
