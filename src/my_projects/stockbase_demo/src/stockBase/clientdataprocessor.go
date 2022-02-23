package stockBase

import (
	"encoding/json"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
)

func QueryBasicInstCode(w http.ResponseWriter, r *http.Request) {
	tradingDayInfo := GetTradingDayInfo();
	w.Write([]byte(tradingDayInfo.StrBasicCodeInfo))
	logs.Info("recieve json is:", tradingDayInfo.StrBasicCodeInfo)

}

func QueryInstCode(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // 解析参数
	var param ReqInstCode
	parseKlineQueryParam(r, &param)

	tradingDayInfo := GetTradingDayInfo()
	nIndex := Count(tradingDayInfo.StrStockVer, ",")

	var outJson string
	if tradingDayInfo.StrStockVer == param.Ver {
		logs.Info("版本(%s)相同，不响应码表查询结果集!", tradingDayInfo.StrStockVer)
		outJson = ResultMsg(ERROR_CODE_VER_SAME, "版本号相同")
	} else if Substr(tradingDayInfo.StrStockVer, 0, nIndex) != Substr(param.Ver, 0, nIndex) || tradingDayInfo.StrStockVer > param.Ver {
		outJson = tradingDayInfo.StrStockCodeInfo;
	} else {
		logs.Info("当前系统码表版本(%s)小于客户端码表(%s)版本，不响应码表查询结果集!", tradingDayInfo.StrStockVer, param.Ver)
		outJson = ResultMsg(ERROR_CODE_VER_SAME, "当前交易日期客户端码表版本更高")
	}

	w.Write([]byte(outJson))
	logs.Info("recieve json is:", outJson)
}

func QueryAdditionalInstCode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数
	var param ReqInstCode
	parseKlineQueryParam(r, &param)
	var outJson string
	var version string
	logs.Info("Ver is ", param.Ver)
	if param.Ver == "" {
		outJson = ResultMsg(ERROR_CODE_JSON_PARAM_ERROR, "json参数错误")
	} else {
		ver := strings.Split(param.Ver, ",")
		for index, v := range ver {
			if index == 1 {
				version = v
			}
		}
		_, s_FtdcInstrumentList := GetInstrumentByTradingDay()
		outJson = string(ProcessHttpFunc(GetTradingDayInfo(), version, s_FtdcInstrumentList))
	}

	w.Write([]byte(outJson))
	logs.Info("recieve json is:", outJson)

}

func SearchInstCodeByKey(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // 解析参数
	var param ReqInstCode
	parseKlineQueryParam(r, &param)

	key := strings.ToLower(param.Key)

	var vecSearcResultCurr []InstCodeInfo
	var vecSearcResultEver []InstCodeInfo

	mutex.Lock()
	retCurr, vecSearcResultCurr := SearchByPrefix(key, s_InstCodeDictionaryCurr, int64(param.Num))
	if retCurr && len(vecSearcResultCurr) < param.Num {
		LastNum := param.Num - len(vecSearcResultCurr)
		_, vecSearcResultEver = SearchByPrefix(key, s_InstCodeDictionaryEver, int64(LastNum))
	}
	mutex.Unlock()

	//量小，直接遍历去重
	for i, pIterEver := range vecSearcResultEver {
		for _, pIterCurr := range vecSearcResultCurr {
			if pIterEver.SecurityCode == pIterCurr.SecurityCode {
				vecSearcResultEver = append(vecSearcResultEver[:i], vecSearcResultEver[i+1:]...) //删除这个元素
				break
			}
		}
	}

	for i, itsec := range vecSearcResultEver {
		temp := vecSearcResultEver[i+1]
		if temp.SecurityCode == itsec.SecurityCode {
			vecSearcResultEver = append(vecSearcResultEver[:i+1], vecSearcResultEver[i+2:]...) //删除temp这个元素
			break
		}
	}

	var icCodeInfo IcCodeInfo
	icCodeInfo.ErrorCode = 0
	icCodeInfo.ErrorMsg = "success"
	icCodeInfo.Comm = ""
	for _, pIter := range vecSearcResultCurr {
		icCodeInfo.Data = append(icCodeInfo.Data, InstInfo{Ei: uint64(pIter.Id), Inst: pIter.SecurityCode, SecNm: pIter.SecurityName,
			Type: pIter.SecurityType, Py: pIter.ShortName, ExchID: pIter.ExchID, Current: pIter.CurrentName})
	}

	for _, pIter := range vecSearcResultEver {
		icCodeInfo.Data = append(icCodeInfo.Data, InstInfo{Ei: uint64(pIter.Id), Inst: pIter.SecurityCode, SecNm: pIter.SecurityName,
			Type: pIter.SecurityType, Py: pIter.ShortName, ExchID: pIter.ExchID, Current: pIter.CurrentName})
	}

	response, _ := json.Marshal(icCodeInfo)
	w.Write([]byte(response))
	logs.Info("recieve json is:", icCodeInfo)

}

func Count(str, pos string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		ch := str[i]
		if ch == ',' {
			break
		}
		count++
	}

	return count
}

func parseKlineQueryParam(r *http.Request, param *ReqInstCode) {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	logs.Info("receive request para is %s", string(body))
	error := json.Unmarshal(body, param)
	if nil != error {
		logs.Error("parse json param error:%s", error.Error())
	}

}

//去除空格
func remove_blanks(string string) (string) {
	string = strings.Replace(string, " ", "", -1)
	return string
}

func SearchByPrefix(prefix string, map_result map[string]InstCodeInfo, max_result_size int64) (ret bool, ICInfo []InstCodeInfo) {
	var count int64
	var str_prefix string
	count = 0

	str_prefix = remove_blanks(prefix)
	if strings.Index(str_prefix, ",") != -1 {
		logs.Error("str_prefix has ',' ")
		return false, nil
	}

	for key, _ := range map_result {
		if !strings.Contains(key, str_prefix) {
			delete(map_result, key)
		} else {
			value, ok := map_result[key]
			if ok {
				ICInfo = append(ICInfo, value)
				count++
				if count == max_result_size && max_result_size != -1 {
					break
				}
			}
		}
	}

	return true, ICInfo
}
