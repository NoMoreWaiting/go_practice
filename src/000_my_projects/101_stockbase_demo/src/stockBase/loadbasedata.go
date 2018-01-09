package stockBase

import (
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"time"
	"log"
	"encoding/json"

	"sync"
)

var (
	s_TradingDayInfo TradingDayInfo	//声明全局变量
	s_InstCodeDictionaryCurr = make(map[string] InstCodeInfo)
	s_InstCodeDictionaryEver = make(map[string] InstCodeInfo)
	mutex sync.RWMutex
)

func GetTradingDayInfo() (TradingDayInfo) {
	mutex.Lock()
	defer mutex.Unlock()
	return s_TradingDayInfo
}



func GetDataFromDb() (bool){
	 var tradingDayInfo TradingDayInfo

	 if GetDataByTradingDay(&tradingDayInfo) {
		 if s_TradingDayInfo.ExchangeTradingDay.TradingDay == tradingDayInfo.ExchangeTradingDay.TradingDay &&
			 s_TradingDayInfo.StrStockVer == tradingDayInfo.StrStockVer {
				logs.Info("【交易日期或码表更新】当前交易日期:", s_TradingDayInfo.ExchangeTradingDay.TradingDay, "未更新，并且码表未更新，无需切换交易日期和更新码表")
				return false
		 }

		 if s_TradingDayInfo.ExchangeTradingDay.TradingDay == tradingDayInfo.ExchangeTradingDay.TradingDay {
			 logs.Info("【交易日期切换更新码表】当前交易日期:", s_TradingDayInfo.ExchangeTradingDay.TradingDay,
			 			"已更新为最新交易日期:", tradingDayInfo.ExchangeTradingDay.TradingDay,
			 			"码表已经更新:",s_TradingDayInfo.StrStockVer, "===>", tradingDayInfo.StrStockVer)
		 } else if s_TradingDayInfo.StrStockVer == tradingDayInfo.StrStockVer {
			 logs.Info("【交易日内更新码表】当前交易日期:",s_TradingDayInfo.ExchangeTradingDay.TradingDay,
			 			"未更新，但码表已经更新:", s_TradingDayInfo.StrStockVer, "===>", tradingDayInfo.StrStockVer)
		 }
		 mutex.Lock()
		 s_TradingDayInfo = tradingDayInfo
		 mutex.Unlock()
		 return true
	 }
	 return false
}

func GetDataByTradingDay(tradingDayInfo *TradingDayInfo) (bool) {
	var num int64
	var mapFtdcTradingDay []FtdcTradingDay
	num, mapFtdcTradingDay = GetCurrentTradingDay()
	if num != 1 {
		logs.Error("【获取当前交易日期】数据库中交易日历配置存在交集的问题，通过当前时间点，获得(%u)个交易日", num)
		return  false
	}
	var ftdcTradingDay FtdcTradingDay
	ftdcTradingDay = mapFtdcTradingDay[0]

	tradingDayInfo.ExchangeTradingDay.TradingDay  = ftdcTradingDay.TradingDay
	tradingDayInfo.ExchangeTradingDay.StartTime   = StringToUnxi(ftdcTradingDay.StartTime)
	tradingDayInfo.ExchangeTradingDay.EndTime = StringToUnxi(ftdcTradingDay.EndTime)
	tradingDayInfo.ExchangeTradingDay.TradeTimes = ftdcTradingDay.TradeTimes
	tradingDayInfo.ExchangeTradingDay.TradeTimeTs = TradeTimes2TimeTs(tradingDayInfo.ExchangeTradingDay)

	logs.Info("【获取当前交易日期:",tradingDayInfo.ExchangeTradingDay.TradingDay,
		"开始时间:",ftdcTradingDay.StartTime, "结束时间:",ftdcTradingDay.EndTime, "交易时段(分):",tradingDayInfo.ExchangeTradingDay.TradeTimes,
			"交易时段(秒):",tradingDayInfo.ExchangeTradingDay.TradeTimeTs)


	tradingDayInfo.StrStockVer = tradingDayInfo.ExchangeTradingDay.TradingDay + "," + GetStockLastestUpdateTime()
	logs.Info("【码表版本】Get StockCode version success:sql[%s]", tradingDayInfo.StrStockVer)


	mutex.Lock()
	num, s_FtdcInstrumentList := GetInstrumentByTradingDay()
	mutex.Unlock()


	//logs.Info("s_FtdcInstrumentList is", s_FtdcInstrumentList)

	tradingDayInfo.StrStockCodeInfo  = string(ProcessHttpFun(*tradingDayInfo, s_FtdcInstrumentList))
	logs.Info("【码表信息】Get StockCodeInfo success, 数据条数: %d", num);

	tradingDayInfo.StrBasicCodeInfo = string(ProcessBasicInstCode(*tradingDayInfo))
	logs.Info("【基础码表信息】Get BasicCodeInfo success")

	return true
}

func ProcessHttpFun(tradingDayInfo TradingDayInfo, ftdc []FtdcInstrument) ([]byte) {
	var tdInfo TDInfo
	tdInfo.ErrorCode = 0
	tdInfo.ErrorMsg = "success"		//"failed"
	tdInfo.Ver = tradingDayInfo.StrStockVer

	traDate := tradingDayInfo.ExchangeTradingDay.TradingDay
	traTimes := tradingDayInfo.ExchangeTradingDay.TradeTimes
	traTimeTs := tradingDayInfo.ExchangeTradingDay.TradeTimeTs

	tdInfo.Comm.Stock = append(tdInfo.Comm.Stock, TDInfoStock{Type:"A", TraDate:traDate, TraTimes:traTimes, TraTimeTs:traTimeTs, Decm:2, PriUnit:"CNY"})
	tdInfo.Comm.Stock = append(tdInfo.Comm.Stock, TDInfoStock{Type:"I", TraDate:traDate, TraTimes:traTimes, TraTimeTs:traTimeTs, Decm:2, PriUnit:"CNY"})

	mutex.Lock()
	RemoveAll(s_InstCodeDictionaryCurr)
	RemoveAll(s_InstCodeDictionaryEver)
	for _, ft := range(ftdc) {

		if len(ft.SecurityName) > 0 && ft.Status != 3 {
			var stInstCodeInfo  InstCodeInfo
			stInstCodeInfo = InstCodeInfo{ft.Id, ft.SecurityCode, ft.SecurityName, ft.SecurityType,
				ft.ShortName, ft.ExchangeID, true, ft.UsedName, ft.UsedShortName}

			if len(ft.SecurityCode) > 0 {
				s_InstCodeDictionaryCurr[ft.SecurityCode] = stInstCodeInfo
			}

			if len(ft.SecurityName) > 0 {
				s_InstCodeDictionaryCurr[ft.SecurityName] = stInstCodeInfo
			}

			if len(ft.ShortName) > 0 {
				s_InstCodeDictionaryCurr[ft.ShortName] = stInstCodeInfo
			}

			if stInstCodeInfo.SecurityType == "A" {
				stInstCodeInfo.CurrentName = false;
				used_names := strings.Split(ft.UsedName, ",")
				used_short_names := strings.Split(ft.UsedShortName, ",")

				// 过滤掉以非中文字符ST,G,*开头的曾用名及其对应的曾用名拼音
				minsize := min(used_names, used_short_names)
				for i := 0; i < minsize -1; i++ {
					if 'G' == used_names[i][0] || 'S' == used_names[i][0] || '*' == used_names[i][0] {
						continue
					}
					stInstCodeInfo.SecurityName= used_names[i]
					stInstCodeInfo.ShortName = used_short_names[i]
					s_InstCodeDictionaryEver[used_names[i]] = stInstCodeInfo
					s_InstCodeDictionaryEver[used_short_names[i]] = stInstCodeInfo
				}
			}
		}

		if Strncasecmp(ft.SecurityCode, "sh", 2) {
			securityCode := Substr(ft.SecurityCode, 0, 2)
			ft.SecurityCode = securityCode
			ft.ExchangeID = "sh"
		} else if Strncasecmp(ft.SecurityCode, "sz", 2) {
			securityCode := Substr(ft.SecurityCode, 0, 2)
			ft.SecurityCode = securityCode
			ft.ExchangeID = "sz"
		}

		ft.DbInsertTime = string("")
		ft.NAV = float64(0)
		ft.CirculatedShare = float64(0)
		ft.TotalShare = float64(0)
		if ft.HGT > 0 || ft.SGT > 0 {
			ft.GT = 1
		} else {
			ft.GT = 0
		}
		ft.HGT = int(0)
		ft.SGT = int(0)
		ft.DecimalNum = int(0)
		ft.PriceUnit = string("")
		ft.LaunchPrice = float64(0)
		ft.UsedName = string("")
		ft.UsedShortName =string("")

		tdInfo.Comm.Data = append(tdInfo.Comm.Data, TDInfoData{Id:ft.Id, SecurityCode:ft.SecurityCode, SecurityName:ft.SecurityName, SecurityType:ft.SecurityType,
							ShortName:ft.ShortName, EPS:ft.EPS, NAV:ft.NAV, TotalShare:ft.TotalShare, CirculatedShare:ft.CirculatedShare,
							PriceUnit:ft.PriceUnit, LaunchPrice:ft.LaunchPrice, Status:ft.Status, DecimalNum:ft.DecimalNum, DbInsertTime:ft.DbInsertTime,
							RZRQ:ft.RZRQ, SGT:ft.SGT, HGT:ft.HGT, GT:ft.GT, ExchangeID:ft.ExchangeID, UsedName:ft.UsedName, UsedShortName:ft.UsedShortName})
	}
	mutex.Unlock()

	Outputstring, err := json.Marshal(tdInfo)
	if err != nil {
		logs.Error("error:", err)
	}
	return  Outputstring
	
}

func ProcessHttpFunc(tradingDayInfo TradingDayInfo, strStockVer string, ftdc []FtdcInstrument) ([]byte)  {
	var tdInfo TDInfo
	tdInfo.ErrorCode = 0
	tdInfo.ErrorMsg = "success"		//"failed"
	tdInfo.Ver = tradingDayInfo.StrStockVer

	traDate := tradingDayInfo.ExchangeTradingDay.TradingDay
	traTimes := tradingDayInfo.ExchangeTradingDay.TradeTimes
	traTimeTs := tradingDayInfo.ExchangeTradingDay.TradeTimeTs

	tdInfo.Comm.Stock = append(tdInfo.Comm.Stock, TDInfoStock{Type:"A", TraDate:traDate, TraTimes:traTimes, TraTimeTs:traTimeTs, Decm:2, PriUnit:"CNY"})
	tdInfo.Comm.Stock = append(tdInfo.Comm.Stock, TDInfoStock{Type:"I", TraDate:traDate, TraTimes:traTimes, TraTimeTs:traTimeTs, Decm:2, PriUnit:"CNY"})

	mutex.Lock()
	for _, ft := range ftdc {
		ft.DbInsertTime = string("")
		ft.NAV = float64(0)
		ft.CirculatedShare = float64(0)
		ft.TotalShare = float64(0)
		if ft.HGT > 0 || ft.SGT > 0 {
			ft.GT = 1
		} else {
			ft.GT = 0
		}
		ft.HGT = int(0)
		ft.SGT = int(0)
		ft.DecimalNum = int(0)
		ft.PriceUnit = string("")
		ft.LaunchPrice = float64(0)
		ft.UsedName = string("")
		ft.UsedShortName =string("")

		tdInfo.Comm.Data = append(tdInfo.Comm.Data, TDInfoData{Id:ft.Id, SecurityCode:ft.SecurityCode, SecurityName:ft.SecurityName, SecurityType:ft.SecurityType,
			ShortName:ft.ShortName, EPS:ft.EPS, NAV:ft.NAV, TotalShare:ft.TotalShare, CirculatedShare:ft.CirculatedShare,
			PriceUnit:ft.PriceUnit, LaunchPrice:ft.LaunchPrice, Status:ft.Status, DecimalNum:ft.DecimalNum, DbInsertTime:ft.DbInsertTime,
			RZRQ:ft.RZRQ, SGT:ft.SGT, HGT:ft.HGT, GT:ft.GT, ExchangeID:ft.ExchangeID, UsedName:ft.UsedName, UsedShortName:ft.UsedShortName})
	}
	mutex.Unlock()

	Outputstring, err := json.Marshal(tdInfo)
	if err != nil {
		log.Println("error:", err)
	}
	return  Outputstring
}



func ProcessBasicInstCode(tradingDayInfo TradingDayInfo) ([]byte) {
	var icInfo InCodeInfo
	icInfo.ErrorCode = 0
	icInfo.ErrorMsg = "success"		//"failed"
	icInfo.Data = append(icInfo.Data, InstCodeInfoData{Type:"A", TraDate:tradingDayInfo.ExchangeTradingDay.TradingDay, TraTimes:tradingDayInfo.ExchangeTradingDay.TradeTimes, TraTimeTs:tradingDayInfo.ExchangeTradingDay.TradeTimeTs, Decm:2, PriUnit:"CNY"})
	icInfo.Data = append(icInfo.Data, InstCodeInfoData{Type:"I", TraDate:tradingDayInfo.ExchangeTradingDay.TradingDay, TraTimes:tradingDayInfo.ExchangeTradingDay.TradeTimes, TraTimeTs:tradingDayInfo.ExchangeTradingDay.TradeTimeTs, Decm:2, PriUnit:"CNY"})

	Outputstring, err := json.Marshal(icInfo)
	if err != nil {
		log.Println("error:", err)
	}
	return  Outputstring
}

//strncasecmp忽略大小写比较字符串
func Strncasecmp(source string, compare string, n int) bool {
	var ret  = false
	var i = 0
	var j = 0
	if len(compare) < n || len(source) < n {
		logs.Error("char length error")
		return ret
	}
	lowerSource := strings.ToLower(source)
	lowerCompare := strings.ToLower(compare)
	for lowerSource[i] == lowerCompare[j] {
		i++
		j++
		if i == n || j == n {
			ret = true
			return ret
		}
	}
	return ret
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//比较两个数组切片的大小
func min(str1 []string, str2 []string) int  {
	if len(str1) <= len(str2) {
		minLen := len(str1)
		return minLen
	} else  {
		minLen := len(str2)
		return minLen
	}
}


func Split(s []string) string {
	var v string
	for _, v = range s {
		strings.Split(v, ";")
	}

	return v
}

//string转unix
func StringToUnxi(str string) int64 {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, str, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()	//将time_t时间格式转为unix时间戳格式
	return sr
}

func TradeTimes2TimeTs(exchangeTradingDay ExchangeTradingDay) (string) {
	var start int
	var end int
	var outSlice []string
	var result string

	year := Substr(exchangeTradingDay.TradingDay, 0, 4)
	month := Substr(exchangeTradingDay.TradingDay, 4, 2)
	day := Substr(exchangeTradingDay.TradingDay, 6, 2)
	unixTime := StringToUnxi(strings.Join([]string{year, month, day}, "-") + " 00:00:00")

	slice := strings.Split(exchangeTradingDay.TradeTimes, ";")
	for _, v := range(slice) {
		slice2 := strings.Split(v, "-")
		for index, v2 := range(slice2) {
			if index == 0 {
				start, _ = strconv.Atoi(v2)
			} else if index == 1 {
				end, _ = strconv.Atoi(v2)
			} else {
				logs.Error("index error")
				break
			}
		}

		intStart := int(unixTime) + start*60
		intEnd := int(unixTime) + end*60

		outString := (strconv.Itoa(intStart)) + "-" + (strconv.Itoa(intEnd))
		outSlice = append(outSlice, outString)
		result= strings.Join(outSlice, ";")
	}

	return result
}

//清空map里的所有数据
func RemoveAll(m map[string] InstCodeInfo)   {
	for key, _ := range m {
		delete(m, key)
	}
}

