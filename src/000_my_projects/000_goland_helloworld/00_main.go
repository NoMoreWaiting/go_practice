package main

import (
	"fmt"
	"reflect"
)

type Foo struct {
	name string
}
type Bar struct {
	name string
}

//用于保存实例化的结构体对象
var regStruct map[string]interface{}

func testStructReflect() {
	str := "Bar"
	if regStruct[str] != nil {
		t := reflect.ValueOf(regStruct[str]).Type()
		v := reflect.New(t).Elem()
		tt := reflect.ValueOf("songyunxuan")
		fmt.Println(tt)
		//v.FieldByName("name").Set(tt)
		fmt.Println(v)
	}
}

func init() {
	regStruct = make(map[string]interface{})
	regStruct["Foo"] = Foo{}
	regStruct["Bar"] = Bar{}
}

// main 函数, 测试用例
func main() {
	//testHelloWorld()

	/*algorithm*/
	//testSingleNonDuplicate()

	/*go programming language*/
	//testGoProLan()

	/*http server and client*/
	//testHttpClient()
	//testHttpServer()
	//testHttpSayHello()
	//testHttpReadJson()
	//testHttpTime()

	/*解析xml, json*/
	//testConfigXML()
	//testParseJson()

	/*连接数据库*/
	//testGoMySQL()

	/*反射实验*/
	testReflectMethod()
	testReflectDemo()
	testReflectInterface()
	//(&ReflectDemo{}).testReflect()
	//testStructReflect()
	//testReflect()

	/*interface*/
	//testCusTag()

	/*and so on*/
	//testTimeSwitch()
	//testPeriodTask()
}
