package p2p

import (
	"blockchain/chain"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Node struct {
	Address     string
	Peers       []string
	Connections map[string]net.Conn
	Mutex       sync.Mutex
}

func NewNode(address string, peers []string) *Node {
	return &Node{
		Address:     address,
		Peers:       peers,
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
			err = json.Unmarshal([]byte(message), &tx)
			if err != nil {
				fmt.Println("Error unmarshalling transaction:", err)
				return err
			} else {
				fmt.Println("Received transaction:", tx)
				err := blockchain.AddTransactionToPool(tx)
				if err != nil {
					return err
				}
				fmt.Println("Transaction pool:", blockchain.PendingTransactions)
			}
		case "block":
			var block chain.Block
			err = json.Unmarshal([]byte(message), &block)
			if err != nil {
				fmt.Println("Error unmarshalling block:", err)
				return err
			} else {
				fmt.Println("Received block:", block)
				blockchain.AddBlock(block)
				fmt.Println("Blocks:", blockchain.Blocks)
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
	node.AddConnection(peerAddress, conn)

	// Notify the peer about itself
	selfMessage := "PEER:" + node.Address + "\n"
	_, err = conn.Write([]byte(selfMessage))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		node.RemoveConnection(conn.RemoteAddr().String())
		return
	}
	fmt.Println("Notified peer about self:", selfMessage)

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
		fmt.Println("Connection removed:", peerAddress)
	} else {
		fmt.Println("No connection found for:", peerAddress)
	}

	for i, peer := range node.Peers {
		if peer == peerAddress {
			node.Peers = append(node.Peers[:i], node.Peers[i+1:]...)
			fmt.Println("Peer removed from list:", peerAddress)
			break
		}
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
