package main

import (
	"reflect"
	"fmt"
	"database/sql"
	"encoding/json"
	"strings"
	"strconv"
	"time"
)

// http包body 头部结构体
type json_header struct {
	Ver       string `json:"Ver"`       // 版本
	PointerID string `json:"PointerId"` // 探针ID
	OperTime  int64  `json:"OperTime"`  // 采样时间（时间戳)
}

// 根据结构体字段, 查询数据库, 并返回http消息的body
func query(conn *DBConn, lastID *int, tableStruct interface{}) (resultString string) {
	var queryStr string = "select ";
	var strField string = makeField(tableStruct)
	var tableName string = reflect.TypeOf(tableStruct).Name()
	queryStr += strField + " from " + tableName + " where id > "
	queryStr += strconv.Itoa(*lastID)
	queryStr += " limit " + strconv.Itoa(config.SendConf.MaxPost)
	fmt.Println("queryStr: ", queryStr)

	rows, err := conn.MY_DB.Query(queryStr)
	checkErr(err)
	m := parseRows(rows)
	results := convTableFactory(tableStruct, m, lastID)
	if len(results) == 0 {
		return
	}
	rowsJson, err := json.Marshal(results)
	checkErr(err)
	resultString = "\"Data\":" + string(rowsJson)
	resultString = appendHead(resultString)
	return
}

// 将数据库所有字段都转化为string类型的结果集, 可以避开数据库查询字段为空的问题, 不过也消耗了性能
func parseRows(rows *sql.Rows) (map[int]map[string]string) {
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	results := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		if err := rows.Scan(scanArgs...); nil != err {
			fmt.Printf("%d scan error: %s\n", i, err)
			return nil
		}
		row := make(map[string]string)
		for k, v := range values {
			key := columns[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}
	return results
}

// 使用模板编程的实验, go语言并不支持泛型编程
func makeResultTemplete(rows *sql.Rows, tableStruct interface{}) (string) {
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var strBody string = "\"Data\":["
	tempTable := tableStruct
	tt := reflect.ValueOf(&tempTable).Elem()

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		checkErr(err)
		fmt.Println(values)

		newTable := tt.Interface()
		fmt.Println("newTable: type", newTable)

		for i, col := range values {
			// 1. 根据这里的结构体变量的类型设置相应的转换函数
			s := reflect.ValueOf(&newTable).Elem()
			fmt.Println("columns[i]: ", columns[i], ",  type:", reflect.TypeOf(columns[i]))
			//fmt.Println(s.Field(i))

			newObject := NewObject(tableStruct)
			fmt.Println("newObject", newObject)
			ttt := reflect.Indirect(reflect.ValueOf(newTable))
			fmt.Println("tttt", ttt)

			switch reflect.Value(s.FieldByName(columns[i])).Kind() {
			case reflect.Int64:
				temp, _ := strconv.ParseInt(string(col), 10, 64)
				v := reflect.ValueOf(temp)
				s.FieldByName(columns[i]).Set(v)
			case reflect.Float64:
				temp, _ := strconv.ParseInt(string(col), 10, 64)
				v := reflect.ValueOf(temp)
				s.FieldByName(columns[i]).Set(v)
			case reflect.String:
				v := reflect.ValueOf(string(col))
				s.FieldByName(columns[i]).Set(v)
			}

			// 2. 根据从数据库中读取出来的值的类型来设置转换函数. 工作量较大, 但是可以复用
		}
		body, err := json.Marshal(newTable)
		checkErr(err)
		strBody += string(body)
		strBody += ","
	}
	strBody = strings.TrimRight(strBody, ",")
	strBody += "]"

	return appendHead(strBody)
}

// 添加http包头
func appendHead(strBody string) (string) {
	var tHeader json_header
	tHeader.Ver = "1.0.0.0"
	tHeader.PointerID = "go test"
	tHeader.OperTime = int64(time.Now().Unix())

	jsonHead, err := json.Marshal(tHeader)
	checkErr(err)
	strJson := string(jsonHead)
	strJson = strings.TrimRight(strJson, "}")
	strJson += "," + strBody + "}"
	return strJson
}

// 根据结构体字段反射组成mysql要查询的string字段
func makeField(tableStruct interface{}) (string) {
	tType := reflect.TypeOf(tableStruct)
	//tValue := reflect.ValueOf(tableStruct)
	var str string
	for i := 0; i < tType.NumField(); i++ {
		str += tType.Field(i).Name
		str += ", "
	}

	str = strings.TrimRight(str, ", ") // 这里修剪的字符串, 要和添加的相同
	//fmt.Println("makeField: ", str)
	return str
}
