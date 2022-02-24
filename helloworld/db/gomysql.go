package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 匿名导入的方式(在包路径前添加 _) 此驱动会自行初始化并注册自己到Golang的database/sql上下文中, 可以通过 database/sql 包提供的方法访问数据库
)

//插入demo
func insert() {
	// open 打开的是连接池, 等到执行(Exec, 或者query)的时候才会调用线程执行
	db, err := sql.Open("mysql", "tradepro:trade123456ing@tcp(192.168.19.192:3306)/cpp_quantdo_hxw?charset=utf8")
	defer db.Close()
	checkErr(err)
	stmt, err := db.Prepare(`INSERT user (user_name,user_age,user_sex) values (?,?,?)`)
	defer stmt.Close()
	checkErr(err)
	res, err := stmt.Exec("tony", 20, 1)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
}

//查询demo
func query() {
	db, err := sql.Open("mysql", "tradepro:trade123456ing@tcp(192.168.19.192:3306)/cpp_quantdo_hxw?charset=utf8")
	defer db.Close()
	checkErr(err)
	rows, err := db.Query("SELECT * FROM user") // query 的字符串直接在外层使用 string 自己拼接
	checkErr(err)

	//普通demo
	/*
		for rows.Next() {
			var userId int; var userName string; var userAge int; var userSex int

			rows.Columns()
			err = rows.Scan(&userId, &userName, &userAge, &userSex)
			checkErr(err)

			fmt.Println(userId); fmt.Println(userName); fmt.Println(userAge); fmt.Println(userSex)
		}
	*/

	// 字典类型
	// 构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns)) // 存储values中值的地址, 供Scan扫描使用
	values := make([]interface{}, len(columns))
	for i := range values {
		fmt.Println("len(values): ", len(values))
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// 将行数据保存到record字典
		//err = rows.Scan(scanArgs[0],scanArgs[1],scanArgs[2],scanArgs[3] )
		err = rows.Scan(scanArgs...)
		checkErr(err)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				// 对于任意类型的 interface{} 对象, 要转换为string(或者其他类型时), 需要使用 interface{}.([]byte) 转换为字节数组, 然后再转化为string string(interface{}.([]byte) )
				record[columns[i]] = string(col.([]byte))
			}
		}

		fmt.Println(record)
	}
}

//更新数据
func update() {
	db, err := sql.Open("mysql", "tradepro:trade123456ing@tcp(192.168.19.192:3306)/cpp_quantdo_hxw?charset=utf8")
	defer db.Close()
	checkErr(err)
	stmt, err := db.Prepare(`UPDATE user SET user_age=?,user_sex=? WHERE user_id=?`)
	defer stmt.Close()
	checkErr(err)
	res, err := stmt.Exec(21, 2, 1)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

//删除数据
func remove() {
	db, err := sql.Open("mysql", "tradepro:trade123456ing@tcp(192.168.19.192:3306)/cpp_quantdo_hxw?charset=utf8")
	defer db.Close()

	checkErr(err)
	stmt, err := db.Prepare(`DELETE FROM user WHERE user_id=?`)
	defer stmt.Close()
	checkErr(err)
	res, err := stmt.Exec(2)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestGoMySQL() {
	insert()
	update()
	remove()
	query()
}
