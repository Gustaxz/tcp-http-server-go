package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 1 {
		fmt.Println("Usage: hexdump <file>")
		os.Exit(1)
	}

	file := args[1]

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file: ", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		fmt.Println("Error getting file info: ", err.Error())
		os.Exit(1)
	}

	size := st.Size()
	buf := make([]byte, size)
	_, err = f.Read(buf)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(buf))
	fmt.Println(hex.Dump(buf))
}
