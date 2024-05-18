package main

import (
	"bufio"
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

func (node *Node) StartServer() {
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
		go node.HandleConnection(conn)
	}
}

func (node *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Received:", message)

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

	selfMessage := "PEER:" + node.Address + "\n"
	conn.Write([]byte(selfMessage))
	fmt.Println("Notified peer about self:", selfMessage)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received:", message)
	}
}

func (node *Node) ConnectToPeer(address string) {
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

	go node.ReadData(conn)
}

func (node *Node) ReadData(conn net.Conn) {
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
