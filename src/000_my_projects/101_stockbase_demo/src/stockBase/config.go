// Config.go
package stockBase

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	SysConfig Config
	//第一个参数，为参数名称，第二个参数为默认值，第三个参数是说明
	configFilePath = flag.String("c", "../conf/cfg.xml", "Config.xml File Path")
)

type Basic struct {
	XmlName      xml.Name `xml:"basic"`
	Logdir       string   `xml:"logdir"`
	Logprefix    string   `xml:"logprefix"`
	Runlog       int      `xml:"runlog"`
	ListenIpPort string   `xml:"listenIpPort"`
}

type DBConfig struct {
	XmlName  xml.Name `xml:"dbconfig"`
	Host     string   `xml:"host"`
	Port     int      `xml:"port"`
	UserName string   `xml:"username"`
	PassWord string   `xml:"password"`
	DataBase string   `xml:"database"`
}

type Config struct {
	XmlName  xml.Name `xml:"config"`
	Basic    Basic    `xml:"basic"`
	DBConfig DBConfig `xml:"dbconfig"`
}

func (config Config) String() string {
	return fmt.Sprintf("【Logdir】: %s\n 【Host】 : %s\n - 【Port】 : %d\n 【MySqlSource】 : %s\n",
		config.Basic.Logdir, config.DBConfig.Host, config.DBConfig.Port, config.MySqlSource())
}

func (config Config) MySqlSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.DBConfig.UserName, config.DBConfig.PassWord, config.DBConfig.Host, config.DBConfig.Port, config.DBConfig.DataBase)
}

func (basic Basic) String() string {
	return fmt.Sprintf(" logdir : %s,logprefix : %s, runlog : %d, ListenIpPort:%s",
		basic.Logdir, basic.Logprefix, basic.Runlog, basic.ListenIpPort)
}

/*
	读取配置文件
*/
func GetConfig() (Config, bool) {
	flag.Parse()

	fmt.Println("config opening file: ", *configFilePath) //配置文件路径
	var config Config

	xmlFile, err := os.Open(*configFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return config, false
	}
	defer xmlFile.Close()

	XmlData, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(XmlData, &config)

	fmt.Println(config)

	return config, true
}
