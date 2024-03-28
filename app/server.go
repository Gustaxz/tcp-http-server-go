package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleRequest(msg []byte) ([]byte, error) {
	strMsg := string(msg)
	strMsgSplit := strings.Split(strMsg, "\r\n")

	for i, line := range strMsgSplit {

		switch i {
		case 0:
			parts := strings.Split(line, " ")
			fmt.Println("Method: ", parts[0])
			fmt.Println("Path: ", parts[1])
			fmt.Println("Protocol: ", parts[2])

			if parts[1] != "/" {
				fmt.Println("Invalid path")
				return []byte("HTTP/1.1 404 Not Found\r\n\r\n"), nil
			}

			return []byte("HTTP/1.1 200 OK\r\n\r\n"), nil
		}

	}

	return []byte("HTTP/1.1 200 OK\r\n\r\n"), nil
}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on Port 4221")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println("Connection closed")
			os.Exit(0)
		}
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(buf))
	msg, err := handleRequest(buf)
	if err != nil {
		fmt.Println("Error handling request: ", err.Error())
		os.Exit(1)
	}
	conn.Write(msg)
	fmt.Println("Response sent")
}
