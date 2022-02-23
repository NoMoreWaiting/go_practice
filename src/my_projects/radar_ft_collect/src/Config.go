package main

import (
	"encoding/xml"
	"io/ioutil"
	"fmt"
)

var pathMap map[string]string
var lastID1, lastID2 = 0, 0                      // go语言没有静态变量, 使用全局变量代替
var time_interval int64 = 5 * 1000 * 1000 * 1000 // 纳秒为单位
var config Config

// 初始化, 自动调用, 没有参数和返回值
func init() {
	pathMap = make(map[string]string)
	pathMap["t_market_transfer"] = "/ysradar/biz/ftmoneyio/"
	pathMap["t_market_trade"] = "/ysradar/biz/fttrade/"
	readConfig()
}

// 读取配置文件
func readConfig() {
	content, err := ioutil.ReadFile("conf/cfg.xml")
	checkErr(err)
	err = xml.Unmarshal(content, &config)
	checkErr(err)
	fmt.Println(config)
	time_interval = int64(config.SendConf.PostInter) * 1000 * 1000 * 1000

}

// 配置文件
type Config struct {
	XMLName  xml.Name `xml:"config"`
	DBConfig DBConfig `xml:"trade_mysql"`
	SendConf SendConf `xml:"send_config"`
}

// 数据库配置
type DBConfig struct {
	XMLName xml.Name `xml:"trade_mysql"`
	IP      string   `xml:"ip,attr"`
	Port    string   `xml:"port,attr"`
	User    string   `xml:"user,attr"`
	Pwd     string   `xml:"pwd,attr"`
	DBName  string   `xml:"db_name,attr"`
	CharSet string   `xml:"character_set,attr"`
}

// 转换为数据库连接字串
func (self *DBConfig) ToString() (string) {
	//"tradepro:trade123456ing@tcp(192.168.19.192:3306)/riskcontrol?charset=utf8"
	return self.User + ":" + self.Pwd + "@tcp(" + self.IP + ":" + self.Port + ")/" + self.DBName + "?charset=" + self.CharSet
}

// http 发送地址
type SendConf struct {
	XMLName   xml.Name `xml:"send_config"`
	IP        string   `xml:"ip,attr"`
	Port      string   `xml:"port,attr"`
	PostInter int      `xml:"post_interval,attr"`
	MaxPost   int      `xml:"max_post,attr"`
}

// 转换为http连接字串
func (self *SendConf) ToString() (string) {
	return "http://" + self.IP + ":" + self.Port
}
