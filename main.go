package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("failed to accept connection. Error:%v\n", err)
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("failed to close tcp connection. Error:%v\n", err)
		}
	}()

	for {
		buf := make([]byte, 1024)

		// Read message from client
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("failed to accept connection. Error:%v\n", err)
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}
