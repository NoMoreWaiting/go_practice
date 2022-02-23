package stockBase

type TDInfo struct {
	ErrorCode int     `json:"errorCode"`
	ErrorMsg  string  `json:"errorMsg"`
	Ver       string  `json:"Ver"`
	Comm      Comment `json:"Comm"`
}

type Comment struct {
	Stock []TDInfoStock `json:"stock"`
	Data  []TDInfoData  `json:"data"`
}

type TDInfoStock struct {
	Type      string
	TraDate   string
	TraTimes  string
	TraTimeTs string
	Decm      int
	PriUnit   string
}

type TDInfoData struct {
	Id              int64
	SecurityCode    string  //股票代码
	SecurityName    string  // 股票中文名称
	SecurityType    string  // 股票类型
	ShortName       string  // 合约名称
	EPS             float64 // 最近年度摊薄每股收益(eps)
	NAV             float64 // 每股净资产
	TotalShare      float64 // 总股本 (万股)
	CirculatedShare float64 // 流通股本(万股)
	PriceUnit       string  // 计价单位(CNY,US,HK)
	LaunchPrice     float64 // 发行价
	Status          int     // 状态【0-无该记录 1-上市正常交易 2-未上市 3-退市】
	DecimalNum      int     // 价格小数点后有效位数
	DbInsertTime    string  // 数据库插入时间
	RZRQ            int     // 融资融券标志
	SGT             int     // 深港通标志
	HGT             int     // 沪港通标志
	GT              int     // 沪深港通标志
	ExchangeID      string  // 交易所ID
	UsedName        string  // 曾用名
	UsedShortName   string  // 曾用名拼音
}

type InCodeInfo struct {
	ErrorCode int                `json:"errorCode"`
	ErrorMsg  string             `json:"errorMsg"`
	Data      []InstCodeInfoData `json:"data"`
}

type InstCodeInfoData struct {
	Type      string
	TraDate   string
	TraTimes  string
	TraTimeTs string
	Decm      int
	PriUnit   string
}
type ErrorCode struct {
	ErrorCode int
	ErrorMsg  string
}

type IcCodeInfo struct {
	ErrorCode int        `json:"errorCode"`
	ErrorMsg  string     `json:"errorMsg"`
	Comm      string     `json:"Comm"`
	Data      []InstInfo `json:"data"`
}

type InstInfo struct {
	Ei      uint64
	Inst    string
	SecNm   string
	Type    string
	Py      string
	ExchID  string
	Current bool
}

type ReqInstCode struct {
	Ver string `json:"Ver"`
	Num int    `json:"num"`
	Key string `json:"key"`
}
