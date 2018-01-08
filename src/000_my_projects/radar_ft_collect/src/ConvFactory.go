package main

import (
	"reflect"
	"fmt"
)

var _ = fmt.Println

// 表格结构体转换工厂, 将查询结果集转换为对应结构体的数组, 使用 []interface{} 返回
func convTableFactory(tableStruct interface{}, stringRow map[int]map[string]string, lastId *int) (results []interface{}) {
	switch tableStruct.(type) {
	case t_market_transfer:
		results = convMarketTransfer(tableStruct, stringRow, lastId)
	case t_market_trade:
		results = convMarketTrade(tableStruct, stringRow, lastId)
	}
	return results
}

// 转换 t_market_transfer 结构体
func convMarketTransfer(tableStruct interface{}, stringRow map[int]map[string]string, lastId *int) (results []interface{}) {
	for i := 0; i < len(stringRow); i++ {
		row := stringRow[i]
		// todo 这里有一步断言的操作, 需要解决
		// 1. 不使用断言
		// 2. incorrect 断言的类型, 使用一个函数进行替代, 做成反射的目的
		// reflect.TypeOf() 是reflect描述类型的接口类型, 本身并不是定义新值的时候使用的类型(如 int ,string等等)
		//fmt.Printf("%v, ++++ %T, %T\n", reflect.TypeOf(tableStruct), reflect.TypeOf(tableStruct), tableStruct)
		newStruct := tableStruct.(t_market_transfer)
		//fmt.Println(reflect.TypeOf(&newStruct))
		object := reflect.ValueOf(&newStruct)
		myRef := object.Elem()

		for k, v := range row {
			field := myRef.FieldByName(k)
			value, _ := convField(field.Type(), v)
			reflectValue := reflect.ValueOf(value)
			field.Set(reflectValue)
		}
		*lastId = int(myRef.FieldByName("Id").Int())
		results = append(results, newStruct)
	}
	return
}

// 转换 t_market_trade 结构体
func convMarketTrade(tableStruct interface{}, stringRow map[int]map[string]string, lastId *int) (results []interface{}) {
	for i := 0; i < len(stringRow); i++ {
		row := stringRow[i]
		newStruct := tableStruct.(t_market_trade)
		object := reflect.ValueOf(&newStruct)
		myRef := object.Elem()

		for k, v := range row {
			field := myRef.FieldByName(k)
			value, _ := convField(field.Type(), v)
			reflectValue := reflect.ValueOf(value)
			field.Set(reflectValue)
		}
		*lastId = int(myRef.FieldByName("Id").Int())
		results = append(results, newStruct)
	}
	return
}