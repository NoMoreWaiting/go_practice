package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
type Server struct {
	ServerName string
	ServerIP   string
}

type ServerSlice struct {
	Servers []Server
	ServersID  string
}
*/

func testHttpServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9001", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Fprintf(w, "Hi, I love you %s", html.EscapeString(r.URL.Path[1:]))

	// url 路径 /xxx/xxx/bbb  需要自己分割字符串
	fmt.Println(r.URL.Path[1:])

	for t, v := range r.URL.Path {
		fmt.Printf("%c  ", v)
		fmt.Println(t, v)
	}

	if r.Method == "GET" {
		fmt.Println("method:", r.Method) //获取请求的方法

		fmt.Println("username", r.Form["username"])
		fmt.Println("password", r.Form["password"])

		for k, v := range r.Form {
			fmt.Print("key:", k, "; ")
			fmt.Println("val:", strings.Join(v, ""))
		}
	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		fmt.Printf("%s\n", result)

		// 应该修改为单独的函数, 通过递归处理所有的嵌套json
		// 未知类型的推荐处理方法
		fmt.Println("\n未知类型的推荐处理方法")
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{}) // 接口类型的内部转换
		for k, v := range m {

			// value, ok = element.(T)，这里value就是变量的值，ok是一个bool类型，element是interface变量，T是断言的类型
			// v.(type) 只能在 switch 中使用, 外面使用上面的方式
			switch vv := v.(type) {
			case string:
				fmt.Println(k, "is string", vv)
			case int:
				fmt.Println(k, "is int", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				fmt.Println(k, vv)

				// 这里的vv是一个数组, i是下标, u是值, interface{}
				for i, u := range vv {
					fmt.Println(i, u)
					_, ok := u.(string) // 判断是否是string 类型
					fmt.Print("is string ", ok)
					_, ok = u.(interface{})
					fmt.Println("  is interface{} ", ok)
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}

		// 结构已知，解析到结构体
		fmt.Println("\n结构已知，解析到结构体")
		var s ServerSlice
		json.Unmarshal([]byte(result), &s)

		fmt.Println(s.ServersID)

		for i := 0; i < len(s.Servers); i++ {
			fmt.Println(s.Servers[i].ServerName)
			fmt.Println(s.Servers[i].ServerIP)
		}
	}
}

// ---------------------------------
func SayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}

func testHttpSayHello() {
	http.HandleFunc("/helloworld", SayHello)
	http.ListenAndServe(":8080", nil)
}

// ---------------------------------
type ReadJson struct {
	InvestorID string
	CurStorage int
	DealAmount float64
	DealCount  int
	SerCharge  float64
}

func (this *ReadJson) show() {
	fmt.Println("InvestorID: ", this.InvestorID)
	fmt.Println("CurStorage: ", this.CurStorage)
	fmt.Println("DealAmount: ", this.DealAmount)
	fmt.Println("DealCount: ", this.DealCount)
	fmt.Println("SerCharge: ", this.SerCharge)
}

func (this *ReadJson) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	//path := r.URL.String()
	bodyStr := make([]byte, 1024)
	bodyLen, readErr := r.Body.Read(bodyStr)
	defer r.Body.Close() // 延迟函数, 等调用结束时再执行此操作
	tempStrMsg := string(bodyStr[:bodyLen])

	// 如果错误是EOF(读到结尾), 或者错误为空, 那么就是正确
	if io.EOF == readErr || nil == readErr {
		//io.WriteString(w, path)
		io.WriteString(w, tempStrMsg)
		fmt.Println(bodyLen, tempStrMsg)

		// 读取get参数使用
		//for k, v := range r.Form {
		//	fmt.Print("key:", k, "; ")
		//	fmt.Println("val:", v)
		//}

		// 解析json参数
		//b, err := json.Marshal(tempStrMsg)
		//if nil != err {
		//	fmt.Println("encoding failed", err)
		//} else {
		//	fmt.Println("encoding data : ")
		//	fmt.Println(b)
		//	fmt.Println(string(b))
		//}

		// go 语言通道, 比较高级的用法, 精髓
		ch := make(chan string, 1)
		go func(c chan string, str string) {
			c <- str
		}(ch, tempStrMsg)
		strData := <-ch

		fmt.Println("--------------------------------")
		tempReadJson := &ReadJson{}

		err := json.Unmarshal([]byte(strData), &tempReadJson)
		if nil != err {
			fmt.Println("Unmarshal failed", err)
		} else {
			fmt.Println("Unmarshal success")
			tempReadJson.show()
		}
	} else {
		io.WriteString(w, "read err")
	}
}

func testHttpReadJson() {
	http.ListenAndServe(":8080", &ReadJson{})
}

// ---------------------------------
type timeHandler struct {
	format string
}

func (this *timeHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	tm := time.Now().Format(this.format)
	res.Write([]byte("The time is: " + tm))
}

func testHttpTime() {
	mux := http.NewServeMux()
	th11233 := &timeHandler{format: time.RFC1123}
	mux.Handle("/time/rfc1123", th11233)
	th3339 := &timeHandler{format: time.RFC3339}
	mux.Handle("/time/rfc3339", th3339)

	log.Println("Listening...")
	http.ListenAndServe(":8080", mux)
}
