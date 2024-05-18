package main

import (
	"bufio"
	"flag"
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

func main() {
	// Parse command-line arguments
	listenAddress := flag.String("address", "localhost:8080", "Address to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peers to connect to")
	flag.Parse()

	node := Node{
		Address:     *listenAddress,
		Peers:       strings.Split(*peers, ","),
		Connections: make(map[string]net.Conn),
	}

	go startServer(&node)

	// Connect to peers
	for _, peer := range node.Peers {
		if peer != "" {
			go connectToPeer(&node, peer)
		}
	}

	// Read user input to broadcast messages
	go func() {
		stdReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message to broadcast: ")
			msg, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				return
			}
			node.broadcastMessage(strings.TrimSpace(msg))
		}
	}()

	select {} // Keep the main function running
}

func startServer(node *Node) {
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
		go handleConnection(conn, node)
	}
}

func handleConnection(conn net.Conn, node *Node) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read initial hello message
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Received:", message)

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
		fmt.Println("Received:", message)
	}
}

func connectToPeer(node *Node, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}

	// Send initial hello message
	message := "Hello, Blockchain!\n"
	conn.Write([]byte(message))
	fmt.Println("Sent:", message)

	// Send node's address
	conn.Write([]byte(node.Address + "\n"))
	fmt.Println("Sent address:", node.Address)

	node.Mutex.Lock()
	node.Connections[address] = conn
	node.Mutex.Unlock()
	fmt.Println("Connected to peer:", address)

	// Keep the connection open to read messages
	go readData(conn)
}

func readData(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received:", message)
	}
}

func (node *Node) broadcastMessage(message string) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	for address, conn := range node.Connections {
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error writing to peer", address, ":", err)
		}
	}
}
