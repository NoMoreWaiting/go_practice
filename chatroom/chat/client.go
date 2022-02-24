echopackage chat

import (
	"net"
	"bufio"
	"log"
	"os"
	"strings"
)

type Client struct {
	conn     net.Conn
	incoming Message // 无缓冲的channel
	outgoing Message
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan net.Conn // 传递socket连接, net.Conn 引用类型
	name     string
}

func (c *Client) GetName() string {
	return c.name
}

func (c *Client) SetName(name string) () {
	c.name = name
}

func (c *Client) GetIncoming() string {
	return <-c.incoming
}

func (c *Client) PutOutgoing(message string) () {
	c.outgoing <- message
}

func StartClient(address string) {
	conn, err := net.Dial("tcp", address)
	if nil != err {
		log.Fatal(err)
	}
	defer conn.Close()

	client := CreateClient(conn)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	// 终端显示
	go func() {
		for {
			out.WriteString(client.GetIncoming() + "\n")
			out.Flush()
		}
	}()

	// 客户端主线程不要退出
	// 终端读取
	for {
		line, _, _ := in.ReadLine()
		//fmt.Println("---", line)
		client.PutOutgoing(string(line))
		if strings.HasPrefix(string(line), ":quit") {
			return // 可以使用channel通知out退出(select方式), 而不是直接退出
		}
	}
}

func CreateClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	client := &Client{
		conn:     conn,
		incoming: make(Message),
		outgoing: make(Message),
		reader:   reader,
		writer:   writer,
		quiting:  make(chan net.Conn),
		//name: "xxx",  // 由server赋值
	}
	client.Listen() // 调用Listen, 对socket连接进行读写. server和client都需要读写
	return client
}

func (c *Client) Listen() {
	go c.Read()
	go c.Write()

}

func (c *Client) quit() {
	c.quiting <- c.conn
}

func (c *Client) Close() {
	c.conn.Close()
}

// 读取socket连接
func (c *Client) Read() {
	for {
		if line, _, err := c.reader.ReadLine(); nil != err {
			log.Printf("Read error: %s\n", err)
			c.quit()
			return
		} else {
			c.incoming <- string(line)
		}
	}
}

// 写入socket连接
func (c *Client) Write() {
	for data := range c.outgoing {
		if _, err := c.writer.WriteString(data + "\n"); nil != err {
			c.quit()
			return
		}
		if err := c.writer.Flush(); nil != err {
			log.Printf("Write error: %s\n", err)
			c.quit()
			return
		}
	}
}
