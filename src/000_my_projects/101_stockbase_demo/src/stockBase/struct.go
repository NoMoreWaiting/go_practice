package stockBase

type TradingDayInfo struct {
	StrStockVer string
	StrStockCodeInfo string
	StrBasicCodeInfo string
	ExchangeTradingDay ExchangeTradingDay
}

type ExchangeTradingDay struct {
	TradingDay string
	StartTime int64
	EndTime int64
	TradeTimes string
	TradeTimeTs string
}



type InstCodeInfo struct {
	Id 				int64		`json:"Ei"`
	SecurityCode	string		`json:"Inst"`
	SecurityName	string		`json:"SecNm"`
	SecurityType	string		`json:"Type"`
	ShortName		string		`json:"Py"`
	ExchID			string		`json:"ExchID"`
	CurrentName		bool		`json:"Current"`
	UsedName		string		`json:"UsedName"`
	UsedShortName	string		`json:"UsedShortName"`

}







