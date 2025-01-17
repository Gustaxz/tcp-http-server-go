package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/status"
)

// Path: app/server.go

type HttpRequest struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

func formatResponse(msg string, status status.HttpStatus) []byte {
	return []byte(
		"HTTP/1.1 " + fmt.Sprint(status.StatusCode) + " " + status.StatusMsg + "\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg)
}

func formatOcStreamResponse(msg string, status status.HttpStatus) []byte {
	return []byte(
		"HTTP/1.1 " + fmt.Sprint(status.StatusCode) + " " + status.StatusMsg + "\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg)
}

func formatTextResponse(msg string) []byte {
	fmtMsg := strings.Trim(msg, " ")
	return []byte(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(fmtMsg)) + "\r\n\r\n" + fmtMsg)

}

func handleGetFile(path string, directory string) ([]byte, bool, error) {
	fullPath := directory + path
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}

		return nil, false, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, false, err
	}

	fileSize := fileInfo.Size()
	fileBytes := make([]byte, fileSize)
	_, err = file.Read(fileBytes)
	if err != nil {
		return nil, false, err
	}
	return fileBytes, true, nil
}

func handlePostFile(path string, directory string, data []byte) error {
	fmt.Println("File size: ", len(data))
	fullPath := directory + path
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes.Trim(data, "\x00"))
	if err != nil {
		return err
	}

	return nil
}

func handleRequest(msg []byte, directory string) ([]byte, error) {
	strMsg := string(msg)
	strMsgSplit := strings.Split(strMsg, "\r\n")
	httpReq := HttpRequest{}
	headers := make(map[string]string)
	httpReq.Body = strMsgSplit[len(strMsgSplit)-1]

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

	if pathParts[1] == "files" {
		if httpReq.Method == "POST" {
			path := strings.Replace(httpReq.Path, "/files/", "", 1)
			err := handlePostFile(path, directory, []byte(httpReq.Body))
			if err != nil {
				return formatResponse("Internal Server Error", status.InternalServerError), nil
			}
			return status.FormatStatus(status.Created), nil
		}

		path := strings.Replace(httpReq.Path, "/files/", "", 1)
		fileBytes, exists, err := handleGetFile(path, directory)
		if err != nil {
			return formatResponse("Internal Server Error", status.InternalServerError), nil
		}
		data := string(fileBytes)
		fmt.Println(data)
		if exists {
			return formatOcStreamResponse(data, status.OK), nil
		}
		return formatResponse("Not Found", status.NotFound), nil
	}

	return status.FormatStatus(status.NotFound), nil
}

func handleConnection(conn net.Conn, directory string) {
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
	msg, err := handleRequest(buf, directory)
	if err != nil {
		fmt.Println("Error handling request: ", err.Error())
		os.Exit(1)
	}
	conn.Write(msg)
}

func main() {
	directory := "./"
	if len(os.Args) > 2 {
		directory = os.Args[2]
	}
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
		go handleConnection(conn, directory)
	}

}
