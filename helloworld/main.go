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

type People interface {
	Speak(string) string
}

type NewStudent struct{}

// 类型A的指针实现了接口就得传递A的指针, (接口类型值可以由A的指针赋值)
func (n *NewStudent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

type Stduent struct{}

// 如果改成 stu *Stduent, 编译将无法通过. 类型B的值类型实现了接口, 就得将类型B的值赋值给接口值
func (stu Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func testInterface() {
	var peo People = Stduent{}
	think := "bitch"
	fmt.Println(peo.Speak(think))

	var n *NewStudent
	peo = n
	fmt.Println(peo.Speak(think))

	peo = &NewStudent{}
	fmt.Println(peo.Speak(think))
}

func testDefer() {
	println(DeferFunc1(1))
	println(DeferFunc2(1))
	println(DeferFunc3(1))
}

func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc2(i int) int {
	t := i // t的左右与仅为内部函数, 不是整个 DeferFunc2, 不会影响外界的返回值
	defer func() {
		t += 3
	}()
	return t
}

func DeferFunc3(i int) (t int) {
	defer func() {
		t += i
	}()
	return 2
}

const (
	x  = iota // 0
	y         // 1
	z  = "zz" // zz
	kk        // zz
	p  = iota // 4
)

// main 函数, 测试用例
func main() {
	fmt.Println(x, y, z, kk, p)
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s1 = append(s1, s2...) // 如果是 append(s1, s2) 报错. 不能直接合并
	fmt.Println(s1)

	testInterface()
	testMoreInterface()
	testDefer()
	//testHelloWorld()

	/*algorithm*/
	//testSingleNonDuplicate()

	/*go programming language*/
	// testGoProLan()

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
	// testReflectMethod()
	// testReflectDemo()
	// testReflectInterface()
	//(&ReflectDemo{}).testReflect()
	//testStructReflect()
	//testReflect()

	/*interface*/
	//testCusTag()

	/*and so on*/
	//testTimeSwitch()
	//testPeriodTask()
}
