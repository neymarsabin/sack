package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// this is something that we might get from the RESP client as a request
	input := "$6\r\nNeymar\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, _ := reader.ReadByte()
	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	// convert "5" -> 5, as this is the length of the string
	size, _ := reader.ReadByte()
	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// consume \r\n
	reader.ReadByte()
	reader.ReadByte()

	// get strSize = 5 characters from the input string
	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println(string(name))
}
