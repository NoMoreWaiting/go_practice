package textprotocol

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStruct struct {
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	AnalyticsFile     string   `json:"analytics_file"`
	StaticFileVersion int      `json:"static_file_version"`
	StaticDir         string   `json:"static_dir"`
	TemplatesDir      string   `json:"templates_dir"`
	SerTcpSocketHost  string   `json:"serTcpSocketHost"`
	SerTcpSocketPort  int      `json:"serTcpSocketPort"`
	Fruits            []string `json:"fruits"`
}

type Other struct {
	SerTcpSocketHost string   `json:"serTcpSocketHost"`
	SerTcpSocketPort int      `json:"serTcpSocketPort"`
	Fruits           []string `json:"fruits"`
}

type Object []interface{}

func TestParseJson() {

	// 有问题, 下面的解析未适配
	/*
		jsonStr := `{"host": "http://localhost:9090",
					"port": 9090,
					"analytics_file": "",
					"static_file_version": 1,
					"static_dir": "E:/Project/goTest/src/",
					"templates_dir": "E:/Project/goTest/src/templates/",
					"serTcpSocketHost": ":12340",
					"serTcpSocketPort": 12340,
					"fruits": ["apple", "peach"]}`
	*/

	// 注意这里的字符串的限界符是 ` , 1左边的那个键, 不是'"引号
	jsonStr := `{"accessToken":"507b5e08ee444dck887b66bd08672905",
				"clientToken":"64e3a5415bfe405d9485f1jf2ea5c68e",
				"selectedProfile":{"id":"selID","name":"Bluek404"},
				"availableProfiles":[{"id":"测试ava","name":"Bluek404"}]}`

	//json str 转map
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
		fmt.Println("\n================= json str 转map ====================")
		fmt.Println(data)

		for x, y := range data {
			fmt.Println(x, y)
		}

		mapTmp := data["selectedProfile"].(map[string]interface{})
		fmt.Println("selectedProfile.id: ", mapTmp["id"])

		var data2 map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &data2); err == nil {
			fmt.Println("data2---> ")
			for x, y := range data2 {
				fmt.Println(x, " : ", y)
			}
			fmt.Println("<--- data2")
		}

		mapTmp2 := (data["availableProfiles"].([]interface{}))[0].(map[string]interface{})
		//mapTmp3 := mapTmp2[0].(map[string]interface {})
		fmt.Println("mapTmp2[\"id\"]: ", mapTmp2["id"])
	}

	//json str 转struct
	var config ConfigStruct
	if err := json.Unmarshal([]byte(jsonStr), &config); err == nil {
		fmt.Println("\n================ json str 转struct =================")
		fmt.Println("ConfigStruct: ", config)
		fmt.Println("ConfigStruct.Host: ", config.Host)
	}

	//json str 转struct(部份字段)
	var part Other
	if err := json.Unmarshal([]byte(jsonStr), &part); err == nil {
		fmt.Println("\n================ json str 转struct =================")
		fmt.Println("Other: ", part)
		fmt.Println("Other.SerTcpSocketPort: ", part.SerTcpSocketPort)
	}

	//struct 到json str
	if b, err := json.Marshal(config); err == nil {
		fmt.Println("\n================ struct 到json str==================")
		fmt.Println(string(b))
	}

	//map 到json str
	fmt.Println("\n================ map 到json str=====================")
	enc := json.NewEncoder(os.Stdout) // 这里的 os.Stdout 如何使用的?
	enc.Encode(data)

	//array 到 json str
	fmt.Println("\n================ array 到 json str==================")
	arr := []string{"hello", "apple", "python", "golang", "base", "peach", "pear"}
	lang, err := json.Marshal(arr)
	if err == nil {
		fmt.Println(string(lang))
	}

	//json 到 []string
	fmt.Println("\n================ json 到 []string===================")
	var str []string
	if err := json.Unmarshal(lang, &str); err == nil {
		fmt.Println(str)
	}
}
