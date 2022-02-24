package textprotocol

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
)

const CONFIG_PATH = "config/"

type Result struct {
	Person []Person `xml:"person"`
}

type Person struct {
	Name      string    `xml:"name,attr"`
	Age       int       `xml:"age,attr"`
	Career    string    `xml:"career"`
	Interests Interests `xml:"interests"`
}

type Interests struct {
	Interest []string `xml:"interest"`
}

func (person *Person) Chkis18() (flag bool) {
	if person.Age > 18 {
		flag = true
	}
	return flag
}

type Checker interface {
	Chkis18() (flag bool)
}

func test1() {
	content, err := ioutil.ReadFile(CONFIG_PATH + "config1.xml")
	if err != nil {
		log.Fatal(err)
	}
	var result Result
	err = xml.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result.Person[0].Name) // 先不要使用反射, 先使用基础的来解决
	log.Println(result.Person)
	log.Println(result)
}

type StringResources struct {
	XMLName        xml.Name         `xml:"resources"`
	ResourceString []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"name,attr"`
	InnerText  string   `xml:",innerxml"` // 注意这里的属性中不要有空格, 严格按照,逗号来区分的
}

func test2() {
	content, err := ioutil.ReadFile(CONFIG_PATH + "config2.xml")
	if nil != err {
		log.Fatal(err)
	}
	var result StringResources
	err = xml.Unmarshal(content, &result)
	if nil != err {
		log.Fatal(err)
	}

	log.Println(result)
	log.Println(result.ResourceString)
	for t, o := range result.ResourceString {
		fmt.Println("t: ", t, " o: ", o)
		log.Println(o.StringName + "===" + o.InnerText)
	}
}

type Config struct {
	XMLName     xml.Name `xml:"config"`
	BasicConfig []Basic  `xml:"basic"`
	CoMysql     CoMysql  `xml:"co_mysql"`
	CoRedis     CoRedis  `xml:"co_redis"`
	QuoteSvr    QuoteSvr `xml:"quote_svr"`
	TimeOut     TimeOut  `xml:"time_out"`
	RmqGw       RmqGw    `xml:"rmq_gw"`
}
type Basic struct {
	XMLName          xml.Name `xml:"basic"` // 可以省略, 但是加上会更好, 不容易出错
	Name             string   `xml:"name,attr"`
	Group_id         int      `xml:"group_id,attr"`
	Local_ip         string   `xml:"local_ip,attr"`
	Run_log          int      `xml:"run_log,attr"`
	Async_queue_size int      `xml:"async_queue_size,attr"`
}

// type public struct {
// 	IP   string `xml:"ip,attr"`
// 	Port int    `xml:"port,attr"`
// 	User string `xml:"user,attr"`
// }

type CoMysql struct {
	XMLName xml.Name `xml:"co_mysql"`
	//Public  public // 不能直接这样写, xml文件中的结构并没有嵌套
	IP             string `xml:"ip,attr"`
	Port           int    `xml:"port,attr"`
	User           string `xml:"user,attr"`
	Pwd            string `xml:"pwd,attr"`
	DbName         string `xml:"db_name,attr"`
	ConnectTimeOut int    `xml:"connect_time_out,attr"`
	CharacterSet   string `xml:"character_set,attr"`
}
type CoRedis struct {
	XMLName        xml.Name `xml:"co_redis"`
	IP             string   `xml:"ip,attr"`
	Port           int      `xml:"port,attr"`
	User           string   `xml:"user,attr"`
	Pwd            string   `xml:"pwd,attr"`
	DbName         string   `xml:"db_name,attr"`
	ConnectTimeOut int      `xml:"connect_time_out,attr"`
	CharacterSet   string   `xml:"character_set,attr"`
}
type QuoteSvr struct {
	XMLName     xml.Name `xml:"quote_svr"`
	IP          string   `xml:"ip,attr"`
	Port        string   `xml:"port,attr"`
	ReadTimeOut int      `xml:"read_time_out,attr"`
}

type TimeOut struct {
	XMLName       xml.Name `xml:"time_out"`
	CheckInterval int      `xml:"check_interval,attr"`
	LoginTimeOut  int      `xml:"login_time_out,attr"`
	OrderTimeOut  int      `xml:"order_time_out,attr"`
}

type RmqGw struct {
	XMLName xml.Name `xml:"rmq_gw"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName  xml.Name `xml:"item"`
	IP       string   `xml:"ip,attr"`
	Port     string   `xml:"port,attr"`
	BranchID int      `xml:"branch_id,attr"`
	SvrID    int      `xml:"svr_id,attr"`
	GwType   int      `xml:"gw_type,attr"`
	BrokerID string   `xml:"broker_id,attr"`
}

func test3() {
	content, err := ioutil.ReadFile(CONFIG_PATH + "config3.xml")
	if nil != err {
		log.Fatal(err) // 会直接结束程序
	}
	var Config Config
	err = xml.Unmarshal(content, &Config)
	if nil != err {
		log.Fatal(err)
	}

	//fmt.Println(Config.BasicConfig) // 单个元素直接引用

	for t, o := range Config.BasicConfig {
		fmt.Println("t: ", t, "    o: ", o)
		log.Println(o.Local_ip, o.Async_queue_size)
	}
	fmt.Println(Config.CoRedis)
	fmt.Println(Config.CoMysql)
	fmt.Println(Config.QuoteSvr)
	fmt.Println(Config.TimeOut)
	fmt.Println(Config.RmqGw)
	for t, o := range Config.RmqGw.Items {
		fmt.Println("t: ", t, "   o: ", o)
	}
}

func TestConfigXML() {
	test1()
	test2()
	test3()
}
