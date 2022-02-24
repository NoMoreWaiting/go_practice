package stockBase

import (
	"database/sql"
	"github.com/astaxie/beego/orm"
)

var db *sql.DB

///合约交易日 t_ftdctradingday
type FtdcTradingDay struct {
	ID         int64  `orm:"column(ID)"`
	TradingDay string `orm:"column(TradingDay)"`
	StartTime  string `orm:"column(StartTime)"`
	EndTime    string `orm:"column(EndTime)"`
	TradeTimes string `orm:"column(TradeTimes)"`
}

/*
   市净率=每股股价/每股净资产
   流通值=流通股本*每股股价*10000(流通股本的单位为万股)
   总市值=总股本*每股股价*10000(总股本的单位为万股)
   市盈(静)亏损=股价/最近年度摊薄每股收益(eps)
   ????市盈(动)亏损=股价/((最近一个报告期的净利润（亿元）)/总股本)*4

   流通市值指在某特定时间内当时可交易的流通股股数乘以当时股价得出的流通股票总价值
*/

///合约  t_stcode
type FtdcInstrument struct {
	Id              int64   `orm:"column(ID)"`
	SecurityCode    string  `orm:"column(SecurityCode)"`    //股票代码
	SecurityName    string  `orm:"column(SecurityName)"`    // 股票中文名称
	SecurityType    string  `orm:"column(SecurityType)"`    // 股票类型
	ShortName       string  `orm:"column(ShortName)"`       // 合约名称
	EPS             float64 `orm:"column(EPS)"`             // 最近年度摊薄每股收益(eps)
	NAV             float64 `orm:"column(NAV)"`             // 每股净资产
	TotalShare      float64 `orm:"column(TotalShare)"`      // 总股本 (万股)
	CirculatedShare float64 `orm:"column(CirculatedShare)"` // 流通股本(万股)
	PriceUnit       string  `orm:"column(PriceUnit)"`       // 计价单位(CNY,US,HK)
	LaunchPrice     float64 `orm:"column(LaunchPrice)"`     // 发行价
	Status          int     `orm:"column(Status)"`          // 状态【0-无该记录 1-上市正常交易 2-未上市 3-退市】
	DecimalNum      int     `orm:"column(DecimalNum)"`      // 价格小数点后有效位数
	DbInsertTime    string  `orm:"column(DbInsertTime)"`    // 数据库插入时间
	RZRQ            int     `orm:"column(RZRQ)"`            // 融资融券标志
	SGT             int     `orm:"column(SGT)"`             // 深港通标志
	HGT             int     `orm:"column(HGT)"`             // 沪港通标志
	GT              int     `orm:"column(GT)"`              // 沪深港通标志
	ExchangeID      string  `orm:"column(ExchangeID)"`      // 交易所ID
	UsedName        string  `orm:"column(UsedName)"`        // 曾用名
	UsedShortName   string  `orm:"column(UsedShortName)"`   // 曾用名拼音

}

func init() {

	SysConfig, _ = GetConfig()

	// set default database
	orm.RegisterDataBase("default", "mysql", SysConfig.MySqlSource(), 30)

	// register model
	orm.RegisterModel(new(FtdcTradingDay))
	orm.RegisterModel(new(FtdcInstrument))

	// create table
	orm.RunSyncdb("default", false, true)
}
