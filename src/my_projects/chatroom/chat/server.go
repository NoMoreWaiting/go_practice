package chat

import (
	"net"
	"log"
	"strings"
	"fmt"
	"os"
)

const (
	MAX_CLIENTS = 50
)

type Server struct {
	listener net.Listener
	clients  ClientTable
	tokens   Token
	pending  chan net.Conn
	quiting  chan net.Conn
	incoming Message
	outgoing Message
}

func (s *Server) genToken() {
	s.tokens <- 0
}

func (s *Server) takeToken() {
	<-s.tokens
}

func CreateServer() *Server {
	server := &Server{
		clients:  make(ClientTable, MAX_CLIENTS),
		tokens:   make(Token, MAX_CLIENTS), // 带缓存的Channel
		pending:  make(chan net.Conn),
		quiting:  make(chan net.Conn),
		incoming: make(Message),
		outgoing: make(Message),
	}
	//server.listen()
	return server
}

func (s *Server) Start(connString string) {

	s.listen() // select监听Channel通道消息

	var err error
	s.listener, err = net.Listen("tcp", connString)
	if nil != err {
		log.Println(err)
		os.Exit(-2)
	}

	log.Printf("Server %p starts\n", s)

	// 生成最大数目的tokens
	// filling the tokens
	for i := 0; i < MAX_CLIENTS; i++ {
		s.genToken()
	}

	for {
		conn, err := s.listener.Accept()
		if nil != err {
			log.Println(err)
			return
		}
		log.Printf("A new connectiong %v kicks\n", conn)
		s.takeToken() // 取走token
		s.pending <- conn // 新客户端接入
	}
}

// fixme: need to figure out if this is the correct approach to gracefully terminate a server
func (s *Server) Stop() {
	s.listener.Close()
}

// Server开始监听Channel动态
func (s *Server) listen() {
	go func() {
		for {
			select {
			case conn := <-s.pending:
				s.join(conn)
			case message := <-s.incoming:
				s.broadcast(message)
			case conn := <-s.quiting:
				s.leave(conn)
			}
		}
	}()
}

// 新User加入, Server启动2个goroutinue服务: 消息解析和发送, 退出请求
func (s *Server) join(conn net.Conn) {
	client := CreateClient(conn)
	name := getUniqName()
	client.SetName(name)
	s.clients[conn] = client
	log.Printf("Auto assigned name for conn %p: %s\n", conn, name)

	// 获取客户端发送的消息, 解析处理
	go func() {
		for {
			msg := <-client.incoming
			log.Printf("Got message: %s from client %s\n", msg, client.GetName())

			if strings.HasPrefix(msg, ":") {
				if cmd, err := parseCommand(msg); nil == err {
					if err = executeCommand(s, client, cmd); err == nil {
						continue
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			}
			// fallthrough to normal message if it is not parsable or executable
			s.incoming <- fmt.Sprintf("%s says: %s", client.GetName(), msg)
		}
	}()

	// 监听客户端退出
	go func() {
		for {
			conn := <-client.quiting
			log.Printf("Client %s is quiting", client.GetName())
			s.quiting <- conn
		}
	}()
}

// 广播客户端发送的消息
func (s *Server) broadcast(message string) {
	log.Printf("Broadcasting message: %s\n", message)
	for _, client := range s.clients { // fixme: close要和range做同步. 这里是单线程的(如果不指定goroutinue线程数, 那么就是单线程的)
		client.outgoing <- message
	}
}

// 客户端断开连接
func (s *Server) leave(conn net.Conn) {
	if nil != conn {
		conn.Close()
		delete(s.clients, conn)
	}
	s.genToken() // 生成一个新token, 产生新空位
}


