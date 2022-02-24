package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <port>, %s <type>: server or client \n", os.Args[1], os.Args[2])
		os.Exit(-1)
	}

	if os.Args[2] == "server" {
		mainServer(os.Args[1])
	} else if os.Args[2] == "client" {
		mainClient(os.Args[1])
	} else {
		fmt.Println("type must be server or client")
		os.Exit(-1)
	}

}
