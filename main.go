package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var clients []net.Conn

func main() {
	var option int

	fmt.Println("Welcome to TTS")
	fmt.Println("Please choose an option:")
	fmt.Println("1. Start as Server")
	fmt.Println("2. Join a Server")

	fmt.Scanln(&option)
	switch option {
	case 1:
		beSever()
	case 2:
		beClient()
	}

}

func beSever() {
	var port string
	fmt.Println("Please input port:")
	fmt.Scanln(&port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		println("Listen failed!")
		return
	}
	defer listener.Close()

	println("Server started. Waiting for someone......")

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Something wrong!!! (connect)")
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	clients = append(clients, conn)

	reader := bufio.NewReader(conn)
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n")

	fmt.Printf("%s joined the chat...\n", username)
	for _, client := range clients {
		client.Write([]byte(fmt.Sprintf("%s joined the chat...\n", username)))
	}

	for {
		message, err := reader.ReadString('\n')
		message = strings.Trim(message, "\n")
		if err != nil {
			return
		}
		if message == "/exit" {
			for i := 0; i < len(clients); i++ {

				if clients[i] == conn {
					clients = append(clients[:i], clients[i+1:]...)
					i--
					continue
				}
				clients[i].Write([]byte(fmt.Sprintf("%s exited the chat\n", username)))
			}

			fmt.Printf("%s exited\n", username)
			break
		}

		fmt.Printf("%s:%s\n", username, message)
		for _, client := range clients {
			client.Write([]byte(fmt.Sprintf("%s: %s\n", username, message)))
		}
	}
}

func beClient() {
	println("Please enter the server address:")
	reader := bufio.NewReader(os.Stdin)
	address, _ := reader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.Trim(address, "\r\n"))
	if err != nil {
		println("Connect failed")
		return
	}
	defer conn.Close()

	println("Set your username:")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\r\n")
	conn.Write([]byte(username + "\n"))

	go receiveMessage(conn)

	for {
		message, _ := reader.ReadString('\n')
		message = strings.Trim(message, "\r\n")

		fmt.Print("\033[A\033[K")
		conn.Write([]byte(message + "\n"))

		if message == "/exit" {
			break
		}
	}
}

func receiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			println("Receive message failed")
			break
		}
		fmt.Print(message)
	}

}
