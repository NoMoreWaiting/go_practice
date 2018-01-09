package main

import (
	"net/http"
	"log"
	. "000_my_projects/101_stockbase_demo/src/stockBase"
	"github.com/astaxie/beego/logs"
)

func main() {

	logs.NewLogger(10000)          // 创建一个日志记录器，参数为缓冲区的大小
	logs.SetLogger("console", "")  // 设置日志记录方式：控制台记录
	logs.SetLevel(logs.LevelDebug) // 设置日志写入缓冲区的等级：Debug级别（最低级别，所以所有log都会输入到缓冲区）
	logs.EnableFuncCallDepth(true) // 输出log时能显示输出文件名和行号（非必须）*/

	SysConfig, _ = GetConfig()

	if !GetDataFromDb() {
		log.Fatal("GetDataFromDb() Failure!")
	}

	http.HandleFunc("/stockbase/queryBasicInstCode", QueryBasicInstCode)
	http.HandleFunc("/stockbase/queryInstCode", QueryInstCode)
	http.HandleFunc("/stockbase/queryAdditionalInstCode", QueryAdditionalInstCode)
	http.HandleFunc("/stockbase/searchInstCodeByKey", SearchInstCodeByKey)

	err1 := http.ListenAndServe(SysConfig.Basic.ListenIpPort, nil)
	if err1 != nil {
		log.Fatal("Listen:", err1)
	}

}
