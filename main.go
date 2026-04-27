package main

import (
	"fmt"
	"net"
	"os"
	"strings"

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
		requestValue, err := r.Read()
		if err != nil {
			fmt.Printf("failed to read from connection. Error:%v", err)
			os.Exit(1)
			return
		}

		writer := resp.NewWriter(conn)
		if requestValue == nil {
			err = writer.Write(resp.Value{Typ: "error", Err: "value read was empty"})
			if err != nil {
				fmt.Printf("failed to write response to user. Err:%v", err)
				os.Exit(1)
			}
			continue
		}

		if len(requestValue.Array) == 0 {
			fmt.Println("no array value detected")
			continue
		}

		// Get the respective handler
		handler := resp.Handler[strings.ToUpper(requestValue.Array[0].Bulk)]
		if handler == nil {
			err = writer.Write(resp.Value{Typ: "error", Err: "no handler for the given command"})
			if err != nil {
				fmt.Printf("failed to write response to user. Err:%v", err)
				os.Exit(1)
			}
			continue
		}

		response := handler(requestValue.Array)
		err = writer.Write(response)
		if err != nil {
			fmt.Printf("failed to write response to user. Err:%v", err)
			continue
		}
	}
}
