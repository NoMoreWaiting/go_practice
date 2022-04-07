package main

import (
	"bytes"
	"fmt"
	_ "helloworld/algorithm"
	_ "helloworld/gormdemo"
	"helloworld/gotemplate"
	"helloworld/hello"
	_interface "helloworld/interface"
	_reflect "helloworld/reflect"
	"helloworld/textprotocol"
)

// main 函数, 测试用例
func main() {

	/*连接数据库, 需要实际的数据库， 或者改造为 sql mock*/
	// db.TestGoMySQL()

	// hello
	hello.TestHelloWorld()
	hello.TestCustomTag()

	// interface
	_interface.TestInterface()

	/*反射实验*/
	_reflect.TestInterface()

	/*解析xml, json*/
	textprotocol.TestConfigXML()
	textprotocol.TestParseJson()

	wr := &bytes.Buffer{}
	gotemplate.TranslateTemplate(wr, "gotemplate/")
	fmt.Println(wr.String())

	p := new([8]int)
	fmt.Printf("p: %v\n", p)
	fmt.Printf("*p: %v\n", *p)
	fmt.Printf("p == nil? %v\n", p == nil)
	// fmt.Printf("*p == nil? %v\n", *p == nil) // 语法错误， 不允许将 [8]int 和 nil 比较

	pp := new([]int)
	fmt.Printf("pp: %v\n", pp)
	fmt.Printf("*pp: %v\n", *pp)
	fmt.Printf("pp == nil? %v\n", pp == nil)
	fmt.Printf("*pp == nil? %v\n", *pp == nil)

}
