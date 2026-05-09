package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/DevitoDbug/redis_go_v1/aof"
	"github.com/DevitoDbug/redis_go_v1/resp"
	"github.com/DevitoDbug/redis_go_v1/storage"
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

	storage := storage.NewStorage()
	persistance, err := aof.NewAof()
	if err != nil {
		log.Fatal("could not initialize aof")
	}

	// Initialize in memory data for permanent persistence
	err = persistance.Read(func(f io.Reader) error {
		r := resp.NewResp(f, storage)

		for {
			requestValue, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Printf("failed to read content of persistence storage. Error; %v", err)
					return err
				}
			}

			if requestValue == nil {
				fmt.Println("no value detected")
				continue
			}

			if len(requestValue.Array) == 0 {
				fmt.Println("no array value detected")
				continue
			}

			command := strings.ToUpper(requestValue.Array[0].Bulk)
			args := requestValue.Array[1:]

			handler := r.Handlers[command]
			if handler == nil {
				continue
			}

			_ = handler(args)
		}
		return nil
	})
	if err != nil {
		log.Fatal("could not load persisted data")
	}

	for {
		r := resp.NewResp(conn, storage)
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

		// TODO: Handle other types that are not arrays, in particular strings. Basically convert them into
		// arrays so that the code blow this point remains the same

		if len(requestValue.Array) == 0 {
			fmt.Println("no array value detected")
			continue
		}

		command := strings.ToUpper(requestValue.Array[0].Bulk)
		args := requestValue.Array[1:]

		// Get the respective handler
		handler := r.Handlers[command]
		if handler == nil {
			err = writer.Write(resp.Value{Typ: "error", Err: "no handler for the given command"})
			if err != nil {
				fmt.Printf("failed to write response to user. Err:%v", err)
				os.Exit(1)
			}
			continue
		}

		if command == "SET" || command == "HSET" {
			err := persistance.WriteFile(requestValue.Marshal())
			if err != nil {
				fmt.Printf("failed to write to aof. Err:%v", err)
				os.Exit(1)
			}
		}

		response := handler(args)
		err = writer.Write(response)
		if err != nil {
			fmt.Printf("failed to write response to user. \nErr:%v", err)
			continue
		}
	}
}
