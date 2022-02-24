package main

import (
	"log"
	"strings"
	"strconv"
	"time"
	"reflect"
	"errors"
)

// 根据已有类型对象, 新建一个同类型对象, 以 interface{} 返回
func NewObject(sample interface{}) interface{} {
	t := reflect.ValueOf(sample).Type()
	v := reflect.New(t).Interface()
	return v
}

// 将date时间转换为 linux时间戳 int64类型返回
func timeStr2Unix(toBeCharge string) (int64) {
	toBeCharge = strings.TrimSpace(toBeCharge)
	//timeLayout := "2006-01-02 15:04:05"  							//这个模板很头疼, 有严格的格式要求, 对应 1 2 3 4 5 6 7
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")                     //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	//var strTime string = strconv.FormatInt(sr,10)
	//fmt.Println("string(sr) : ",  strTime)
	return sr
}

// 将 date 时间转换为 linux时间戳, 字符串返回
func timeStr2UnixStr(toBeCharge string) (string) {
	toBeCharge = strings.TrimSpace(toBeCharge)
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	sr := theTime.Unix()
	var strTime string = strconv.FormatInt(sr, 10)
	return strTime
}

// 检查错误类型, 错误程序结束
func checkErr(err error) (bool) {
	if nil != err {
		log.Panic(err)
		return false
	}
	return true
}


// 根据传入的字段类型值, 将string转换为对应类型
func convField(t reflect.Type, v string) (interface{}, error) {
	switch t.Kind() {
	case reflect.String:
		return v, nil
	case reflect.Int:
		return strconv.ParseInt(v, 10, 0)
	case reflect.Int32:
		return strconv.ParseInt(v, 10, 0)
	case reflect.Int64:
		return strconv.ParseInt(v, 10, 0)
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.Float64:
		return strconv.ParseFloat(v, 10)
	default:
		err := errors.New("no support field type")
		return nil, err
	}
}