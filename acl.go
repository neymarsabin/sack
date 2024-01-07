package main

import (
	"sync"
)

type ACL struct {
	user *User
}

type User struct {
	uuid     string
	userName string
	password string
}

var ACLsConfig map[string]ACL
var ACLsMu = sync.RWMutex{}

// define user based commands
func user(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Error wrong number of arguments for 'user': "}
	}

	subCommand := args[0].bulk
	subArg := args[1].bulk

	if subCommand == "create" {
		ACLsMu.Lock()
		newUser := User{userName: subArg}
		newAcl := ACL{user: &newUser}
		ACLsConfig[subArg] = newAcl
		ACLsMu.Unlock()
	}

	return Value{typ: "string", str: "OK"}
}
