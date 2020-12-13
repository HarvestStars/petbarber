package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func main() {
	fmt.Printf("测试 \n")
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4GA9FCqwAmzpcQp5MksB", "bv7IIZ6s5S5OXi1uaDSYy3jotSJ0ZR")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = "13611688224"
	request.SignName = "闪剪帮小程序"
	request.TemplateCode = "SMS_206751651"
	randomCode := rand.Intn(10000) + 1000
	codeStr := strconv.Itoa(randomCode)
	templateRaw := fmt.Sprintf("{\"code\":\"%s\"}", codeStr)
	request.TemplateParam = templateRaw
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
