package main

import (
	"chatroom/chat"
	"fmt"
)

func mainClient(connStr string) {

	fmt.Printf("mainClient Run on %s\n", connStr)

	chat.StartClient(connStr)
}
