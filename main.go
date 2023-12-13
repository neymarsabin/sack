package main

import (
	"fmt"
	"net"
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
		if err != nil {
			fmt.Println(err)
			return
		}

		// print value
		fmt.Println("The parsed value that is read: ", value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}
