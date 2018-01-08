package chat

import (
	"errors"
	"fmt"
	"github.com/tyrchen/goutil/regex"
	"regexp"
)

type Command struct {
	cmd string
	arg string
}

// 函数指针类型
type Run func(server *Server, client *Client, arg string)

const (
	CMD_REGEX = `:(?P<cmd>\w+)\s*(?P<arg>.*)`
)

var (
	commands map[string]Run
)

// init()函数, 无参数无返回值. 包导入时自动调用
func init() {
	commands = map[string]Run{
		"name": changeName,
		"quit": quit,
	}
}

// 使用正则解析命令行命令
func parseCommand(msg string) (cmd Command, err error) {
	r := regexp.MustCompile(CMD_REGEX)
	if values, ok := regex.MatchAll(r, msg); ok {
		cmd.cmd, _ = values[0]["cmd"]
		cmd.arg, _ = values[0]["arg"]
		return
	}
	err = errors.New("Unparsed message: " + msg)
	return
}

// 执行内置支持的功能函数
func executeCommand(s *Server, c *Client, cmd Command) (err error) {
	if f, ok := commands[cmd.cmd]; ok {
		f(s, c, cmd.arg)
		return
	}
	err = errors.New("Unsupported command: " + cmd.cmd)
	return
}

// commands
// 修改名字
func changeName(s *Server, c *Client, arg string) {
	oldName := c.GetName()
	c.SetName(arg)
	s.broadcast(fmt.Sprintf("Notification: %s changed its name to %s", oldName, c.GetName()))
}

// 退出
func quit(s *Server, c *Client, arg string) {
	c.quit()
	s.broadcast(fmt.Sprintf("Notification: %s quit the chat room.", c.GetName()))
}
