package main

import (
	"chatroom/chat"
	"fmt"
)

func mainServer(connStr string) {

	server := chat.CreateServer()
	fmt.Printf("mainServer Run on %s\n", connStr)
	server.Start(connStr)
}
