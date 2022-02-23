package main

import (
	"bytes"
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"time"
	"runtime"
	"sync"
	"reflect"
)

// 推送函数
func PushService(client *httpClient, conn *DBConn, lastID *int, TableType interface{}) {
	var jsonBody string = query(conn, lastID, TableType)
	if jsonBody == ""{
		fmt.Println("There is no new data")
		return
	}
	var tType string = reflect.TypeOf(TableType).Name()
	fmt.Println("PushService: ", jsonBody)
	client.Post(pathMap[tType], jsonBody)
}

type httpClient struct {
	httpUrl     string
	contentType string
}

func (self *httpClient) show() {
	fmt.Println("httpUrl: ", self.httpUrl)
	fmt.Println("contentType", self.contentType)
}

// 这里的函数不是每个包单独的init()函数
func (self *httpClient) init() {
	fmt.Println("httpClient, init()")
	self.httpUrl = config.SendConf.ToString()
	self.contentType = "application/json;charset=utf-8"
	self.show()
}

func (self *httpClient) Post(path string, strBody string) {
	body := bytes.NewBuffer([]byte(strBody))
	res, err := http.Post(self.httpUrl+path, self.contentType, body)

	if nil != err {
		log.Println(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if nil != err {
		log.Println(err)
		return
	}

	fmt.Println("httpResp: ", result)
}


func main_test1() {
	var conn1, conn2 DBConn
	conn1.init()
	conn2.init()
	defer conn1.exit()
	defer conn2.exit()
	var client httpClient
	client.init()
	var market_transfer t_market_transfer;
	var market_trade t_market_trade

	ticker := time.NewTicker(time.Duration(time_interval)) // time.Second * 2
	go func() {
		for value := range ticker.C {
			fmt.Println("ticked at %v", time.Now())
			fmt.Println("value =", value)
			fmt.Printf("lastID1: %d, lastID2: %d \n", lastID1, lastID2)
			PushService(&client, &conn1, &lastID1, market_transfer)
			PushService(&client, &conn2, &lastID2, market_trade)
		}
	}()
	ch := make(chan int) // 主线程调用chan阻塞, 要不然程序结束
	value := <-ch
	fmt.Println("xxxxx value =", value)
}

func main_test2() {
	var conn1, conn2 DBConn
	conn1.init()
	conn2.init()
	defer conn1.exit()
	defer conn2.exit()
	var client httpClient
	client.init()
	var market_transfer t_market_transfer;
	var market_trade t_market_trade

	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		// 在函数退出时调用 Done 来通知 main 函数工作已经完成
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(time_interval))
		for value := range ticker.C {
			//fmt.Printf("ticked at %v\n", time.Now())
			fmt.Println("value =", value)
			fmt.Printf("lastID1: %d\n", lastID1)
			PushService(&client, &conn1, &lastID1, market_transfer)
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(time_interval))
		for value := range ticker.C {
			//fmt.Printf("ticked at %v \n", time.Now())
			fmt.Println("value =", value)
			fmt.Printf("lastID2: %d \n", lastID2)
			PushService(&client, &conn2, &lastID2, market_trade)
		}
	}()

	fmt.Println("Waiting To Finish")
	// 等待 goroutine 结束 一旦两个匿名函数创建 goroutine 来执行，main 中的代码会继续运行
	wg.Wait()
	fmt.Println("\nTerminating Program")
}
