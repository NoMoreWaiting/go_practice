package main

import (
	"os"
	"fmt"
	"000_my_projects/100_chatroom/chat"
)

func mainServer(){
	if len(os.Args) != 2{
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(-1)
	}

	server := chat.CreateServer()
	fmt.Printf("Run on %s\n", os.Args[1])
	server.Start(os.Args[1])
}