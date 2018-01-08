package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
)

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}

func testReflectMethod() {
	myType := &MyType{22, "wowzai"}
	//fmt.Println(myType)     //就是检查一下myType对象内容
	//println("---------------")
	/*
		mtV := reflect.ValueOf(&myType).Elem()
		fmt.Println("Before:",mtV.MethodByName("String").Call(nil)[0])
		params := make([]reflect.Value,1)
		params[0] = reflect.ValueOf(18)
		mtV.MethodByName("SetI").Call(params)
		params[0] = reflect.ValueOf("reflection test")
		mtV.MethodByName("SetName").Call(params)
		fmt.Println("After:",mtV.MethodByName("String").Call(nil)[0])
	*/

	mtV2 := reflect.ValueOf(&myType).Elem()
	fmt.Println("Before:", mtV2.Method(2).Call(nil)[0])
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(18)
	mtV2.Method(0).Call(params)
	params[0] = reflect.ValueOf("reflection test")
	mtV2.Method(1).Call(params)
	fmt.Println("After:", mtV2.Method(2).Call(nil)[0])

}

// ---------------------------------
// 发送的表消息
type TradeTable struct {
	InvestorID string  `投资者ID`
	CurStorage int     `当前存量`
	DealAmount float64 `当日成交金额`
	DealCount  int     `当日成交量`
	SerCharge  float64 `当日手续费`
}

type User struct {
	// go语言中的变量或者结构体变量通过首字母大小写来决定是否被外界引用
	// 要想外界解析所有字段, 那么首字母必须大写. 然后通过 tag 来使得json小写
	Id   int    `json:"id"`   // json id
	Name string `json:"name"` // json name
	Age  int    `json:"age"`  // json age
	Sex  string `json:"-"`    // 直接忽略此字段
}

func (u *User) SayHello() {
	fmt.Println("I'm "+u.Name+", Id is ", u.Id, ". Nice to meet you.")
}

//
func testReflectInter(tableStruct interface{}) {
	newStruct := tableStruct.(User)
	object := reflect.ValueOf(&newStruct)
	myref := object.Elem()
	typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)
		fmt.Println("---", myref.Field(0), myref.FieldByName("Name"), myref.FieldByName("Name").Type(), myref.FieldByName("Name").CanSet())
		fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
	}
}

// 通过反射调用方法
func testReflectDemo() {
	p := 1
	q := p
	// 0xc04215a130, 0xc04215a138 自动赋值为一个新的值
	fmt.Printf("%p, %p \n", &p, &q)

	tongydon := &User{1, "TangXiaoDong", 25, "男"}

	testReflectInter(User{2, "Tang", 26, "男--"})

	object := reflect.ValueOf(tongydon)
	myref := object.Elem()
	typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)
		fmt.Println("---", myref.Field(0), myref.FieldByName("Name"), myref.FieldByName("Name").Type(), myref.FieldByName("Name").CanSet())
		fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
	}
	tongydon.SayHello()
	v := object.MethodByName("SayHello")
	v.Call([]reflect.Value{})
}

type ReflectDemo struct {
}

func (this *ReflectDemo) testReflect() {

	testAny()

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // "*os.File"
	fmt.Printf("%T\n", 3)          // "int"

	va := reflect.ValueOf(3) // a reflect.Value
	fmt.Println(va)          // "3"
	fmt.Printf("%v\n", va)   // "3"
	fmt.Println(va.String()) // NOTE: "<int Value>"
	// 逆操作
	x := va.Interface() // an interface{}
	// 一个空的接口隐藏了值对应 的表示方式和所有的公开的方法, 因此只有我们知道具体的动态类型才能使用类型断言来访问内部的值
	i := x.(int)          // an int
	fmt.Printf("%d\n", i) // "3"

	// 反射, 以后可以深入研究
	var value interface{} = &User{1, "Tom", 12, "nan"}
	v := reflect.ValueOf(value)
	vv := reflect.TypeOf(value)
	fmt.Println(v, vv)

	if true {
		tValue := reflect.ValueOf(TradeTable{})
		tType := reflect.TypeOf(TradeTable{})
		for i := 0; i < tType.NumField(); i++ {
			fmt.Println(tType.Field(i).Name, tType.Field(i).Tag, tValue.Field(i).Type(), tValue.Field(i).Interface())
		}
	}

	v2 := 122
	ty := reflect.ValueOf(&v2).Elem()
	s := ty.Interface()
	fmt.Println("v2's value is : ", v2, ", type is : ", reflect.TypeOf(v2), &v2)
	fmt.Println("ty: ", ty)
	fmt.Println("s's value is : ", s, ", type is : ", reflect.TypeOf(s), &s)

}

// test func Any
func testAny() {
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(Any(x))               // "1"
	fmt.Println(Any(d))               // "1"
	fmt.Println(Any([]int64{x}))      // "[]int64 0x8202b87b0"
	fmt.Println(([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
}

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

type GenericSlice struct {
	elemType   reflect.Type
	sliceValue reflect.Value
}

func (self *GenericSlice) Init(sample interface{}) {
	value := reflect.ValueOf(sample)
	self.sliceValue = reflect.MakeSlice(value.Type(), 0, 0)
	self.elemType = reflect.TypeOf(sample).Elem()
}

func (self *GenericSlice) Append(e interface{}) bool {
	if reflect.TypeOf(e) != self.elemType {
		return false
	}
	self.sliceValue = reflect.Append(self.sliceValue, reflect.ValueOf(e))
	return true
}

func (self *GenericSlice) ElemType() reflect.Type {
	return self.elemType
}

func (self *GenericSlice) Interface() interface{} {
	return self.sliceValue.Interface()
}

func testReflect() {
	gs := GenericSlice{}
	gs.Init(make([]int, 0))
	fmt.Printf("Element Type: %s\n", gs.ElemType().Kind()) // => Element Type:int
	result := gs.Append(2)
	fmt.Printf("Result: %v\n", result)             // => Result: true
	fmt.Printf("sliceValue: %v\n", gs.Interface()) // => sliceValue: [2]
}

// 给空接口类型的结构体参数赋值
type fullName struct {
	FName string `json:"fname"`
	MName string `json:"mName"`
	LName string `json:"lname"`
}

type people struct {
	Name   fullName
	Sex    string `json:"sex"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

func (p people) String() string {
	return "my name is " + p.Name.FName + " " + p.Name.MName + " " + p.Name.LName + "," +
		"sex is " + p.Sex + ",Height is " + strconv.Itoa(p.Height) + ",Weight is " + strconv.Itoa(p.Weight)
}

type cat struct {
	Name  fullName
	Color string `json:"color"`
}

type dog struct {
	Name  fullName
	Color string `json:"color"`
	Breed string `json:"breed"`
}

func testReflectInterface() {
	var p, pp people
	p.Name = fullName{FName: "Fdog", MName: "Sdog", LName: "Ldog"}
	p.Sex = "diaosi"
	p.Height = 170
	p.Weight = 60

	fmt.Println("p:", p)
	fmt.Println("pp:", pp)
	fmt.Printf("%p\n", &p)
	setName(&p, &pp)
	fmt.Println("p:", p)
	fmt.Println("pp:", pp)

}

func setName(param interface{}, resp interface{}) {

	newPeople := *(param.(*people))
	fmt.Printf("%p\n", &newPeople)

	fmt.Printf("%T,  %T\n", param, resp)

	fmt.Printf("%T,  %T\n", reflect.ValueOf(param), reflect.ValueOf(param).Elem())

	fmt.Printf("%T,  %T, %v\n", reflect.ValueOf(param).Elem().FieldByName("Name"), reflect.ValueOf(param).Elem().Field(0), reflect.ValueOf(param).Elem().Field(0).CanSet())

	// param 不同的参数输入 具有同一个Name结构体字段
	full := reflect.ValueOf(param).Elem().FieldByName("Name") // Elem 获取接口包含的原始值

	// 修改值
	full.Set(reflect.ValueOf(fullName{FName: "FFdog", MName: "SSdog", LName: "LLdog"}))

	b, _ := json.Marshal(full.Interface()) // 编码结构体的一部分
	fmt.Println(string(b))

	// 解码到结构体的指定字段
	r := reflect.ValueOf(resp).Elem().FieldByName("Name")
	err := json.Unmarshal(b, r.Addr().Interface()) // r是结构体Name对应的reflect.Value类型值
	fmt.Println(err)

	fmt.Println(r)
}
