package stockBase

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

const (
	ERROR_CODE_SUCCESS               = iota //成功，正确 	value --> 0
	ERROR_CODE_VER_SAME                     // 版本号相同 		value --> 1
	ERROR_CODE_MISS_VER                     //value --> 2
	ERROR_CODE_JSON_FORMAT_INVALID          //value --> 3
	ERROR_CODE_UNKNOWN_REQ                  //value --> 4
	ERROR_CODE_INVALID_HTTP_REQ             //value --> 5
	ERROR_CODE_UNSUPPORT_PERIOD             //value --> 6
	ERROR_CODE_MISS_DATAS                   //value --> 7
	ERROR_CODE_DATAS_NOT_OBJECT             //value --> 8
	ERROR_CODE_MISS_EI                      //value --> 9
	ERROR_CODE_MISS_PERIOD                  //value --> 10
	ERROR_CODE_MISS_QRYDAY                  //value --> 11
	ERROR_CODE_MISS_DAYNUM                  //value --> 12
	ERROR_CODE_DAYNUM_BEYOND                //value --> 13
	ERROR_CODE_NOT_FIND_DATA                //value --> 14
	ERROR_CODE_NOT_INNER_CODE               //value --> 15
	ERROR_CODE_JSON_PARAM_INVALID           //value --> 16
	ERROR_CODE_JSON_PARAM_ERROR             //value --> 17
	ERROR_CODE_CACHEDB_ERROR                //value --> 18
	ERROR_CODE_DBEXCETION                   //value --> 19
	ERROR_CODE_RELOADKLINE_MAXLIMITS        //value --> 20
)

func ResultMsg(errorCode int, errorMsg string) (outJson string) {
	error := ErrorCode{
		errorCode,
		errorMsg,
	}

	b, err := json.Marshal(error)
	if err != nil {
		logs.Error("error:", err)
	}
	logs.Info("err is ", b)
	return outJson
}
