package main

import (
	"helloworld/algorithm"
	"helloworld/hello"
	_interface "helloworld/interface"
	_reflect "helloworld/reflect"
	"helloworld/textprotocol"
)

// main 函数, 测试用例
func main() {
	// algorithm
	algorithm.TestSingleNonDuplicate()

	/*连接数据库, 需要实际的数据库， 或者改造为 sql mock*/
	// db.TestGoMySQL()

	// hello
	hello.TestHelloWorld()
	hello.TestCustomTag()

	// interface
	_interface.TestInterface()

	/*解析xml, json*/
	textprotocol.TestConfigXML()
	textprotocol.TestParseJson()

	/*反射实验*/
	_reflect.TestInterface()

}
