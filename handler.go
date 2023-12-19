package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"DEL":     del,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

// PING command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

// SET command to set a value to a key as a key-pair hashmap
func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'set' command: "}
	}

	key := args[0].bulk
	value := args[1].bulk

	// we use sync.RWMutex because our server is supposed to handle requests concurrently.
	// It is to ensure that SETs map is not modified by multiple threads at the same time
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	fmt.Println("Printing the value of all SETs: ", SETs)
	return Value{typ: "string", str: "OK"}
}

// GET command to get a value from a key name
func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'get' command: "}
	}

	key := args[0].bulk
	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// HSET command
// HSET structure -> map[string]map[string]string
func hset(args []Value) Value {
	fmt.Println("All arguments defined for hset command: ", args)
	if len(args) != 3 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'hset': "}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	fmt.Println("Printing HSET All hashmaps: ", HSETs)

	return Value{typ: "string", str: "OK"}
}

// HGET Command
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'hget': "}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// HGETALL Command
func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'hgetall': "}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "NULL"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}

// DEL commands
func del(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'del': "}
	}

	key := args[0].bulk

	SETsMu.Lock()
	delete(SETs, key)
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}
