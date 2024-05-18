package main

import (
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

func (node *Node) StartServer(chain *Blockchain) {
	listener, err := net.Listen("tcp", node.Address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started on", node.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go node.HandleConnection(conn, chain)
	}
}

func ProcessMessage (message string, chain *Blockchain) error {
	var msgMap map[string]interface{}
	err := json.Unmarshal([]byte(message), &msgMap)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return err
	}

	if msgType, ok := msgMap["type"]; ok {
		switch msgType {
		case "transaction":
			var tx Transaction
			err = json.Unmarshal([]byte(message), &tx)
			if err != nil {
				fmt.Println("Error unmarshalling transaction:", err)
				return err
			} else {
				fmt.Println("Received transaction:", tx)
				chain.AddTransactionToPool(tx)
				fmt.Println("Transaction pool:", chain.PendingTransactions)
			}
		case "block":
			var block Block
			err = json.Unmarshal([]byte(message), &block)
			if err != nil {
				fmt.Println("Error unmarshalling block:", err)
				return err
			} else {
				fmt.Println("Received block:", block)
				chain.AddBlock(block)
				fmt.Println("Blocks:", chain.Blocks)
			}
		default:
			fmt.Println("Unknown message type")
		}
	} else {
		fmt.Println("Message does not contain a type")
	}

	return nil
}

func (node *Node) HandleConnection(conn net.Conn, chain *Blockchain) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    // Read initial hello message
    message, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    fmt.Println("Received Initial:", message)

    // Read peer address
    message, err = reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    fmt.Println("Received peer address:", message)

    peerAddress := strings.TrimSpace(message)
    node.Mutex.Lock()
    node.Connections[peerAddress] = conn
    node.Mutex.Unlock()

    // Notify the peer about itself
    selfMessage := "PEER:" + node.Address + "\n"
    conn.Write([]byte(selfMessage))
    fmt.Println("Notified peer about self:", selfMessage)

    // Keep the connection open to read messages
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading from connection:", err)
            return
        }
        fmt.Println("Received Message:", message)

        err = ProcessMessage(message, chain)
		if err != nil {
			fmt.Println("Error processing message:", err)
			return
		}
    }
}

func (node *Node) ConnectToPeer(address string, chain *Blockchain) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}

	message := "Hello, Blockchain!\n"
	conn.Write([]byte(message))
	fmt.Println("Sent:", message)

	conn.Write([]byte(node.Address + "\n"))
	fmt.Println("Sent address:", node.Address)

	node.Mutex.Lock()
	node.Connections[address] = conn
	node.Mutex.Unlock()
	fmt.Println("Connected to peer:", address)

	go node.ReadData(conn, chain)
}

func (node *Node) ReadData(conn net.Conn, chain *Blockchain) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received in ReadData:", message)
		ProcessMessage(message, chain)
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

func (node *Node) BroadcastTransaction(tx Transaction) {
	txJson, err := json.Marshal(struct {
		Type        string      `json:"type"`
		Transaction Transaction `json:"transaction"`
	}{"transaction", tx})
	if err != nil {
		fmt.Println("Error marshalling transaction:", err)
		return
	}
	node.BroadcastMessage(string(txJson))
}

func (node *Node) BroadcastBlock(block Block) {
	blockJson, err := json.Marshal(struct {
		Type  string `json:"type"`
		Block Block  `json:"block"`
	}{"block", block})
	if err != nil {
		fmt.Println("Error marshalling block:", err)
		return
	}
	node.BroadcastMessage(string(blockJson))
}
