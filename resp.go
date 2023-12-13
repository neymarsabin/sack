package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

// define struct to use for serialization/deserialization process
// holds all the commands and arguments we receive from the client
type Value struct {
	// typ is to determine the data type carried by the https://redis.io/docs/reference/protocol-spec/#resp-protocol-description
	typ string

	// str holds the value of the string received from the https://redis.io/docs/reference/protocol-spec/#resp-simple-strings
	str string

	// num holds the value of the integers received
	num int

	// bulk holds store the string received from the bulk strings
	bulk string

	// array holds all the values received from the arrays
	array []Value
}

// Reader to contain all the methods to help us read from the buffer and store it in the Value struct
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// we need to implement two methods, as they are not available in bufio
// readLine -> reads the line from the buffer
// read one byte at a time until we reach \r which indicates the end of the line
// then return the line without the last 2 bytes, \r\n and also return number of bytes in the line
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger -> reads the integer from the buffer
// parse integer from the line and return the integer, number of bytes and error if exists
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// parsing or deserialization process
// create a method that will read from the buffer recursively. read the Value again for each step of the input we receive, so that we can parse it according to the character at the beginning of the line
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}
