package main

import (
	"os"
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DBConn struct {
	MY_DB *sql.DB
}

// 多包管理时使用, 每个package一个init()
// init() 函数必须无参数, 无返回值
//func init()(){
//
//}

func (self *DBConn) exit() {
	self.MY_DB.Close()
}

/*
	函数参数传值, 闭包传引用! tag: 闭包的函数调用并没有传递任何参数！ 不是传引用，而是 变量 a 的作用域就在这个区域，函数里的赋值就是正常变量赋值
	slice 含 values/count/capacity 等信息, 是按值传递
	按值传递的 slice 只能修改values指向的数据, 其他都不能修改
	slice 是结构体和指针的混合体
	引用类型和传引用是两个概念
*/
func (self *DBConn) init() {
	var driverName string = "mysql"
	var dataSourceName string = config.DBConfig.ToString()

	fmt.Println("driverName: ", driverName)
	fmt.Println("dataSourceName: ", dataSourceName)

	var err error
	self.MY_DB, err = sql.Open(driverName, dataSourceName)
	if checkErr(err) {
		self.MY_DB.SetMaxOpenConns(50)
		self.MY_DB.SetMaxIdleConns(20)
		self.MY_DB.Ping()
	} else {
		// 无法连接数据库, 程序直接结束. log.Fatal 会直接结束程序
		log.Fatal(err)
		os.Exit(-1)
	}
}
