package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

// support ping command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

// set command to set a value to a key as a key-pair hashmap
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

	return Value{typ: "string", str: "OK"}
}

// get command to get a value from a key name
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
