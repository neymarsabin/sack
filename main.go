package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

	// start infinite loop
	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client:: ", err.Error())
			os.Exit(1)
		}

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}
