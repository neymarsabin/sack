package lib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Example of RESP representation
// the RESP array looks like this:
// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
// we can split it into lines instead of using ‘\r\n’ to understand
// *2
// $5
// hello
// $5
// world

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

// this is a writer
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
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

// to read Array
// Skip the first byte because we have already read it in the Read method
// Read the integer that represents the number of elements in the array
// iterate over the array and for each line, call the Read method to parse the type according to the character at the beginning of the line
// with each iteration, append the parsed value to the array in the Value object and return it
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

// to read Bulk
// To read the Bulk, we follow these steps:
// Skip the first byte because we have already read it in the Read method.
// Read the integer that represents the number of bytes in the bulk string.
// Read the bulk string, followed by the ‘\r\n’ that indicates the end of the bulk string.
// Return the Value object.
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

// write the marshal method which will call the specific method for each type based on the Value type
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

// marshal simple strings
// prepare something like "+neymar\r\n"
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshal Bulk String
// prepare bulk serialized data
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshal array
func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// marshal Error
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshal Null
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// a method to take Value and writes the bytes it gets from the Marshal method to the Writer
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
