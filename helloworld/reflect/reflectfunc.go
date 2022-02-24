package _reflect

import (
	"fmt"
	"reflect"
)

//---- 以下是用reflect实现一些类型无关的泛型编程示例
//new object same the type as sample
func New(sample interface{}) interface{} {
	t := reflect.ValueOf(sample).Type()
	v := reflect.New(t).Interface()
	return v
}

// check type of aninterface
func CheckType(val interface{}, kind reflect.Kind) bool {
	v := reflect.ValueOf(val)
	return kind == v.Kind()
}

// if _func is not a functionor para num and type not match,it will cause panic
func Call(_func interface{}, params ...interface{}) (result []interface{}, err error) {
	f := reflect.ValueOf(_func)
	if len(params) != f.Type().NumIn() {
		ss := fmt.Sprintf("The number of params is not adapted.%s", f.String())
		panic(ss)
	}
	var in []reflect.Value
	if len(params) > 0 { //prepare in paras
		in = make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
	}
	out := f.Call(in)
	if len(out) > 0 { //prepare out paras
		result = make([]interface{}, len(out))
		for i, v := range out {
			result[i] = v.Interface()
		}
	}
	return
}

// if ch is not channel, it will panic
func ChanRecv(ch interface{}) (r interface{}) {
	v := reflect.ValueOf(ch)
	if x, ok := v.Recv(); ok {
		r = x.Interface()
	}
	return
}

// reflect fields of a struct
func ReflectStructInfo(it interface{}) {
	t := reflect.TypeOf(it)
	fmt.Printf("interface info:%s %s %s %s\n", t.Kind(), t.PkgPath(), t.Name(), t)
	if t.Kind() == reflect.Ptr { //if it is pointer, get it element type
		tt := t.Elem()
		if t.Kind() == reflect.Interface {
			fmt.Println(t.PkgPath(), t.Name())
			for i := 0; i < tt.NumMethod(); i++ {
				f := tt.Method(i)
				fmt.Println(i, f)
			}
		}
	}
	v := reflect.ValueOf(it)
	k := t.Kind()
	if k == reflect.Ptr {
		v = v.Elem() //指针转换为对应的结构
		t = v.Type()
		k = t.Kind()
	}
	fmt.Printf("value type info:%s %s %s\n", t.Kind(), t.PkgPath(), t.Name())
	if k == reflect.Struct { //反射结构体成员信息
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Printf("%d %v\n", i, f)
		}
		for i := 0; i < t.NumMethod(); i++ {
			f := t.Method(i)
			fmt.Println(i, f)
		}
		fmt.Printf("Fileds:\n")
		f := v.MethodByName("func_name")
		if f.IsValid() { //执行某个成员函数
			arg := []reflect.Value{reflect.ValueOf(int(2))}
			f.Call(arg)
		}
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !f.CanInterface() {
				fmt.Printf("%d:[%v] %v\n", i, t.Field(i), f.Type())
				continue
			}
			val := f.Interface()
			fmt.Printf("%d:[%v] %v %v\n", i, t.Field(i), f.Type(), val)
		}
		fmt.Printf("Methods:\n")
		for i := 0; i < v.NumMethod(); i++ {
			m := v.Method(i)
			fmt.Printf("%d:[%v] %v\n", i, t.Method(i), m)
		}
	}
}
