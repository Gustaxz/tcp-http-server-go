package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func formatTextResponse(msg string) []byte {
	return []byte(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg)

}

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

			if parts[1] == "/" {
				return []byte("HTTP/1.1 200 OK\r\n\r\n"), nil
			}

			paths := strings.Split(parts[1], "/")

			if len(paths) == 3 {
				if paths[1] == "echo" {
					return formatTextResponse(paths[2]), nil
				}
			}

			return []byte("HTTP/1.1 404 Not Found\r\n\r\n"), nil
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
