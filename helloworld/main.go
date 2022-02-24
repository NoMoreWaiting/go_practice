package main

import (
	_ "helloworld/algorithm"
	_ "helloworld/gormdemo"
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

}
