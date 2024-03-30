package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/status"
)

type HttpRequest struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
}

func formatResponse(msg string, status status.HttpStatus) []byte {
	return []byte(
		"HTTP/1.1 " + fmt.Sprint(status.StatusCode) + " " + status.StatusMsg + "\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg)
}

func formatTextResponse(msg string) []byte {
	fmtMsg := strings.Trim(msg, " ")
	return []byte(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(fmtMsg)) + "\r\n\r\n" + fmtMsg)

}
func handleRequest(msg []byte) ([]byte, error) {
	strMsg := string(msg)
	strMsgSplit := strings.Split(strMsg, "\r\n")
	httpReq := HttpRequest{}
	headers := make(map[string]string)

	if len(strMsgSplit) < 1 {
		return formatResponse("Bad Request", status.BadRequest), nil
	}

	for i, line := range strMsgSplit {

		if i == 0 {
			parts := strings.Split(line, " ")
			httpReq.Method = parts[0]
			httpReq.Path = parts[1]
			httpReq.Protocol = parts[2]
		}

		if i > 0 && strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			headers[parts[0]] = parts[1]
		}

	}

	pathParts := strings.Split(httpReq.Path, "/")

	if httpReq.Path == "/" {
		return status.FormatStatus(status.OK), nil
	}

	if pathParts[1] == "echo" {
		echo := strings.Replace(httpReq.Path, "/echo/", "", 1)
		return formatTextResponse(echo), nil
	}

	if pathParts[1] == "user-agent" {
		return formatTextResponse(headers["User-Agent"]), nil
	}

	return status.FormatStatus(status.NotFound), nil
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println("Connection closed")
			os.Exit(0)
		}
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println(string(buf))
	msg, err := handleRequest(buf)
	if err != nil {
		fmt.Println("Error handling request: ", err.Error())
		os.Exit(1)
	}
	conn.Write(msg)
}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on Port 4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Accepted connection from ", conn.RemoteAddr())
		go handleConnection(conn)
	}

}
