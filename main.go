package main

import (
	"fmt"
	"inquiry-transfer/conf"
	"inquiry-transfer/http"
	"strconv"
	"time"
)

var flag string

var cycle time.Duration

var myConfig *conf.Config

var myRequest *http.Request

func init() {
	// 初始化配置信息
	myConfig = new(conf.Config)
	myConfig.InitConfig("./app.conf")
	flag = myConfig.Read("default", "flag")
	fmt.Printf("app version %s", myConfig.Read("default", "version"))
	fmt.Println()
	// 初始化请求结构体
	myRequest = new(http.Request)
	// 初始化循环时间
	_cycle, _ := strconv.Atoi(myConfig.Read(flag, "cycle"))
	cycle = time.Duration(_cycle)
}

func main() {

	inquiryTimer := time.NewTimer(time.Second * cycle)
	for {
		<-inquiryTimer.C
		fmt.Println("start job**********************************************")
		inquiryTransfer() // 询价转单主函数
		fmt.Println("job done**********************************************")
		inquiryTimer.Reset(time.Second * cycle)
	}
}

func inquiryTransfer() {
	// 1 LOGIN
	myRequest.Method = "POST"
	myRequest.Url = myConfig.Read(flag, "url") + "/login"
	myRequest.Params = map[string]string{
		"tel":      myConfig.Read(flag, "phone"),
		"password": myConfig.Read(flag, "password"),
	}
	_, err := myRequest.Request()
	if err != nil {
		fmt.Printf("error:%#v\n", err)
	}
	// 2 TRANSFER
	// 2.1 search
	myRequest.Method = "GET"
	myRequest.Url = myConfig.Read(flag, "url") + "/inquiry/search?command=2&page_index=1&page_count=10"
	data, err := myRequest.Request()
	if err != nil {
		fmt.Printf("error:%#v\n", err)
	}
	//2.2 search
	myRequest.Method = "GET"
	myRequest.Url = myConfig.Read(flag, "url") + "/app/one-category/sub-list?category_id=1"
	data, err = myRequest.Request()
	if err != nil {
		fmt.Printf("error:%#v\n", err)
	}
	fmt.Println(data)
}
