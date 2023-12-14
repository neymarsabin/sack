package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// create a new server to listen to port 6379
	fmt.Println("Listening on port :6379")
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error while starting listener: ", err)
		return
	}

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

		// the value object looks like this when the command is SET name sabin
		// Value{
		// 	typ: "array",
		// 	array: []Value{
		// 		Value{typ: "bulk", bulk: "SET"},
		// 		Value{typ: "bulk", bulk: "name"},
		// 		Value{typ: "bulk", bulk: "sabin"},
		// 	},
		// }

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

		result := handler(args)
		writer.Write(result)
	}
}
