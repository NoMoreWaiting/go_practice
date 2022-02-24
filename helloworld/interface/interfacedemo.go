package _interface

import (
	"fmt"
)

// ---------------------------------
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

// ---------------------------------
type fruit interface {
	catName() string
	getName() string
	setName(name string)
}

type apple struct {
	name string
}

func (a apple) catName() string {
	fmt.Printf("%p\n", &a)
	fmt.Printf("%T\n", a)
	return a.name
}

func (a *apple) getName() string {
	return a.name
}

// 如果想要通过接口方法修改属性, 需要在传入指针的结构体才行
func (a *apple) setName(name string) {
	a.name = name
}

func testMoreInterface() {
	a := apple{"红富士"}
	fmt.Println(a)
	fmt.Printf("%+v\n", a)
	fmt.Printf("%p\n", &a)
	fmt.Printf("%T\n", a)
	fmt.Println(a.getName())
	fmt.Println(a.catName())

	a.setName("树顶红")
	fmt.Println(a.getName())
	fmt.Println(a.catName())

	b := new(apple)
	*b = a
	fmt.Printf("%T\n", b)
	fmt.Println(b.getName())
	fmt.Println(b.catName())
	a.setName("树顶红---")
	fmt.Println(b.getName())
	fmt.Println(b.catName())

	var fruitVar fruit
	fruitVar = &a
	fruitVar = b
	fmt.Println(fruitVar.getName())
}

func TestInterface() {
	testInterface()
	testMoreInterface()
}
