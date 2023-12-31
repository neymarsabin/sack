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
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/cristalhq/acmd"
)

func main() {
	config := &Configuration{}
	// eg: sack start --port=8080
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
			Name:        "start",
			Description: "set port of the application",
			ExecFunc: func(ctx context.Context, args []string) error {
				fs := flag.NewFlagSet("port", flag.ContinueOnError)
				portFlag := fs.Int("port", 6379, "Port number to use")

				if err := fs.Parse(args); err != nil {
					return err
				}

				config.port = *portFlag
				fmt.Println("Server is listening.....")
				portConfiguration := fmt.Sprintf(":%d", config.port)
				listener, err := net.Listen("tcp", portConfiguration)
				if err != nil {
					fmt.Println("Error while starting listener: ", err)
					return err
				}

				// start with persistance file
				persistance, err := NewPersistance("db.sack")
				if err != nil {
					fmt.Println(err)
					return err
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
					return err
				}

				// starting the server
				fmt.Println("Starting the server::: ")

				defer conn.Close()

				for {
					resp := NewResp(conn)
					value, err := resp.Read()

					if err != nil {
						fmt.Println(err)
						return err
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
					return nil
				}
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
}
