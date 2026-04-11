package main

import (
	"fmt"
	"net"
	"os"

	"github.com/DevitoDbug/redis_go_v1/resp"
)

func main() {
	port := ":6379"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("listening to port %v\n", port)

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
		r := resp.NewResp(conn)
		val, err := r.Read()
		if err != nil {
			fmt.Printf("failed to read from connection. Error:%v", err)
			os.Exit(1)
			return
		}

		fmt.Printf("%v: %v\n", port, *val)

		writer := resp.NewWriter(conn)
		err = writer.Write(resp.Value{Typ: "string", Str: "OK"})
		if err != nil {
			fmt.Printf("failed to write response to user. Err:%v", err)
			os.Exit(1)
		}
	}
}
