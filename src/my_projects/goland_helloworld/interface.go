package main

import (
	"fmt"
)

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
	b = &a
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
