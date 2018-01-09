package stockBase

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)



func GetCurrentTradingDay() (int64, []FtdcTradingDay) {
	o := orm.NewOrm()

	var ftdcTradingDay []FtdcTradingDay
	num, err := o.Raw("select TradingDay,StartTime,EndTime,TradeTimes from t_ftdctradingday WHERE StartTime<=CURRENT_TIMESTAMP AND EndTime>=CURRENT_TIMESTAMP;").QueryRows(&ftdcTradingDay)

	if err != nil {
		logs.Error("query error")
	}

	return num, ftdcTradingDay

}

func GetStockLastestUpdateTime() string {
	o := orm.NewOrm()

	var ftdcInstrument []FtdcInstrument
	_, err := o.Raw("SELECT MAX(DbInsertTime) AS DbInsertTime FROM t_stcode WHERE SecurityType='A' OR (SecurityType='I' AND SecurityName!='' AND SecurityName IS NOT NULL) OR SecurityType='GNBK' OR SecurityType='HYBK' OR SecurityType='DYBK'").QueryRows(&ftdcInstrument)

	if err != nil {
		logs.Error("query error")
	}

	var dbInsertTime string

	for _, ft := range (ftdcInstrument) {
		dbInsertTime = ft.DbInsertTime
	}
	return dbInsertTime
}

func GetInstrumentByTradingDay() (int64, []FtdcInstrument) {
	o := orm.NewOrm()

	var ftdcInstrument []FtdcInstrument

	num, err := o.Raw("SELECT CirculatedShare, DbInsertTime, DecimalNum, ID, LaunchPrice, NAV, PriceUnit, SecurityCode, SecurityName, SecurityType, ShortName, Status, TotalShare, UsedName, UsedShortName, HGT, RZRQ, SGT FROM t_stcode WHERE SecurityType='A' OR (SecurityType='I' AND SecurityName IS NOT NULL AND SecurityName!='') OR SecurityType='GNBK' OR SecurityType='HYBK' OR SecurityType='DYBK' ORDER BY SecurityType DESC, SecurityCode ASC").QueryRows(&ftdcInstrument)

	if err != nil {
		logs.Error("query error")
	}

	return num, ftdcInstrument
}



