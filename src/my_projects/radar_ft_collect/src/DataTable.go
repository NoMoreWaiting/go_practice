package main


// 出入金表监控
type t_market_transfer struct{
	// 字段小写, 查询数据库之后解析为json就会隐藏
	Id 				int64 	`json:"DbId"` 			// 源表自增ID
	Direction   	string	`json:"Direction"`		// 出入金方向
	BrokerID 		string 	`json:"BrokerId"`			// 经纪公司代码
	BankID			string	`json:"BankId"`			// 银行ID
	BankBranchID 	string	`json:"BankBranchId"`		// 银行分支ID
	AccountID   	string 	`json:"UserId"`   		// 投资者Id
	BankAccount		string	`json:"BankAccountId"`	// 银行账号
	TradeAmount		float64	`json:"Money"` 			// 出入金金额
	RequestID		int64	`json:"RequestId"`		//
	ErrorCode  		int64	`json:"RC"`				// 出入金异常类型
	ErrorMsg		string	`json:"RM"`				// 异常原因
	ReqInsertTime	int64	`json:"ReqTime"`			// 请求时间
	ResInsertTime	int64	`json:"RspTime"`			// 响应时间
	NetAddr 		string	`json:"NetAddr"`			// 网络物理地址
	ServiceName		string	`json:"ServiceName"`		// 服务名称
	ServiceNode		string	`json:"ServiceNode"`		// 服务节点标识
}

// 交易表监控
type t_market_trade struct{
	Id 				int64 	`json:"DbId"`				// 源表ID
	OrderRef		string	`json:"OrderRef"`			// 交易所返回字段
	BrokerID 		string 	`json:"BrokerId"`			// 经纪公司代码
	ExchangeID 		string 	`json:"ExchangeId"`		// 交易所ID
	InvestorID  	string 	`json:"InvestorId"`		// 投资商ID
	InstrumentID	string	`json:"InstrumentId"`		// 合约ID
	UserID 			string	`json:"UserId"`			// 用户ID
	TradeType		string	`json:"TradeType"`		// 交易类型
	TradeID			string	`json:"TradeId"`			// 成交编号ID
	Direction		string	`json:"Direction"`		// 买卖方向
	OrderSysID 		string	`json:"OrderSysId"`		// 报单系统编号
	OffsetFlag 		string	`json:"OffsetFlag"`		// 开平标志
	HedgeFlag 		string	`json:"HedgeFlag"`		// 投机套保标志
	Price			float64	`json:"Price"`			// 价格
	Volume 			int64	`json:"Volume"`			// 成交数量
	TradeDate 		string	`json:"TradeDate"`		// 交易日期 yyyyMMdd
	TradeTime 		string	`json:"TradeTime"`		// 成交时间 HH:mm:SS
	TraderID		string	`json:"TraderId"`			// 交易员代码
	OrderLocalID    string	`json:"OrderLocalId"`		// 报单本地编号
	ClearingPartID  string	`json:"ClearingPartId"`	// 结算会员编号
	SequenceNo		int64	`json:"SequenceNo"`		// 序号
	SettlementID 	int64	`json:"SettlementId"`		// 结算会员编号
	BrokerOrderSeq  int64	`json:"BrokerOrderSeq"`	// 经纪公司报单编号
	TradingDay      string	`json:"TradingDay"`		// 交易日 yyyyMMdd
	Uname 			string	`json:"Uname"`			//
	InsertTime 		string	`json:"InsertTime"`		//
	NetAddr 		string	`json:"NetAddr"`			// 网络物理地址
	ServiceName		string	`json:"ServiceName"`		// 服务名称
	ServiceNode		string	`json:"ServiceNode"`		// 服务节点标识
}
