// the value object looks like this when the command is SET name sabin
// Value{
// 	typ: "array",
// 	array: []Value{
// 		Value{typ: "bulk", bulk: "SET"},
// 		Value{typ: "bulk", bulk: "name"},
// 		Value{typ: "bulk", bulk: "sabin"},
// 	},
// }

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/cristalhq/acmd"
)

func main() {
	// start with initializing a few commands
	config := &Configuration{}
	cmds := []acmd.Command{
		{
			Name:        "ping",
			Description: "Ping the application",
			ExecFunc: func(ctx context.Context, args []string) error {
				fmt.Println("Pong!")
				return nil
			},
		},
		{
			Name:        "port",
			Description: "set port of the application",
			ExecFunc: func(ctx context.Context, args []string) error {
				config.port = 6379
				fmt.Println("The port is set for configuration: ", config.port)
				return nil
			},
		},
	}

	app := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "sack",
		AppDescription: "a simple ping command library",
	})

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// create a new server to listen to port 6379
	fmt.Println("Listening on port :6379")
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error while starting listener: ", err)
		return
	}

	// start with persistance file
	persistance, err := NewPersistance("db.sack")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer persistance.Close()

	// read logs from file and save it in memory
	// does not have the Read method without using the ~file.Read~ option
	persistance.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid Command: ", command)
			return
		}

		handler(args)
	})

	// start listening on the specified port
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error while listening: ", err)
		return
	}

	// starting the server
	fmt.Println("Starting the server::: ")

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)
		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			persistance.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
