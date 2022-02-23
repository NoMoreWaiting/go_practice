package chat

import "net"

type Message chan string // 无缓冲的channel 消息string
type Token chan int // 无缓冲的channel 限制登录人数
type ClientTable map[net.Conn]*Client // client登录表

