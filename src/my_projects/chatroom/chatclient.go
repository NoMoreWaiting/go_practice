package main

import (
	"fmt"
	"os"
	"000_my_projects/100_chatroom/chat" // 当前目录相对路径引用, . "./chat" 可以省略包名引用函数 alias "./chat" 别名引用
)

func mainClient() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(-1)
	}

	chat.StartClient(os.Args[1])
}
