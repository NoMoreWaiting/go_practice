package _http


import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"bytes"
)


func HttpGet() {
	u, _ := url.Parse("http://localhost:9001/xiaoyue")
	q := u.Query()
	q.Set("username", "user")
	q.Set("password", "passwd")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String());
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s", result)
}


type Server struct {
	ServerName string
	ServerIP   string
}

type ServerSlice struct {
	Servers []Server // 这里定义的是数组, 所有使用 map[string]interface{} 解析出来之后也是一个数组, 数组之内的内容就是  Server 结构体, 可以转换出来的
	ServersID  string
}


func HttpPost() {

	var s ServerSlice
	var newServer Server
	s.ServersID = "team1"
	newServer.ServerName = "Guangzhou_VPN";
	newServer.ServerIP = "127.0.0.1"
	s.Servers = append(s.Servers, newServer)
	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.2"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.3"})

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}

	body := bytes.NewBuffer([]byte(b))
	res,err := http.Post("http://localhost:9001/xiaoyue", "application/json;charset=utf-8", body)
	if err != nil {
		log.Fatal(err)
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s", result)
}


func testHttpClient(){
	//HttpGet()
	HttpPost()
}

