// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package console
*/
package console

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sms/config"
	"sms/controllers"
	"sms/models"
	"time"
)

const HOST = "https://sms-bus.com"
const TimeoutMinutes = 4

type Sms struct {
	ExecTime time.Time // 执行时间
}

// AutoCancel 自动取消短信
func (s Sms) AutoCancel() {
	// 获取参数 暂时不实现

	// 获取数据
	sms := models.SendPhoneNumberList{}
	list, err := sms.GetListByStatus(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 处理数据
	// 更改时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("时区设置不正确：", err)
		return
	}

	// 请求 API 应该在封装下的
	params := setToken(config.Cfg.ApiKey)
	for _, j := range list {
		if j.RequestId == "" {
			continue
		}

		createTime, err := time.ParseInLocation("2006-01-02 15:04:05", j.CreateAt, loc)
		if err != nil {
			continue
		}

		// 小于 4分钟无需取消
		if time.Since(createTime).Minutes() <= TimeoutMinutes {
			continue
		}

		params.Set("request_id", j.RequestId)
		curl := models.BaseCurl{
			Host:   HOST,
			Path:   "/api/control/cancel",
			Params: params,
		}

		resp, err := curl.GET()
		if err != nil {
			fmt.Println("脚本异常")
			continue
		}

		// 回写数据库
		var s controllers.ResultCancelRequest
		_ = json.Unmarshal([]byte(resp), &s)
		if s.Code == 200 {
			sms.SetSmsStatus(j.RequestId, 2, "")
		} else {
			sms.SetSmsStatus(j.RequestId, 3, s.Message)
		}
	}
}

// AutoSmsCode 自动获取短信验证码
func (s Sms) AutoSmsCode() {
	// 获取数据
	sms := models.SendPhoneNumberList{}
	list, err := sms.GetListByStatus(3)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 处理数据
	// 更改时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("时区设置不正确：", err)
		return
	}

	// 请求 API 应该在封装下的
	params := setToken(config.Cfg.ApiKey)
	for _, j := range list {
		if j.RequestId == "" {
			continue
		}

		createTime, err := time.ParseInLocation("2006-01-02 15:04:05", j.CreateAt, loc)
		if err != nil {
			continue
		}

		// 大于 4分钟无需取消
		if time.Since(createTime).Minutes() > TimeoutMinutes {
			//continue
		}

		// 小于 10s 无需处理（大概率没有抓取到）
		if time.Since(createTime).Seconds() <= 10 {
			//continue
		}

		// 定时延迟机制 0 ~ 10 s 分别获取
		//now := time.Now()
		rand.Seed(time.Now().UnixNano())
		millisecond := rand.Intn(10)
		time.Sleep(time.Duration(millisecond) * time.Second)
		//fmt.Println("触发延迟逻辑:", time.Since(now).Seconds())

		params.Set("request_id", j.RequestId)
		curl := models.BaseCurl{
			Host:   HOST,
			Path:   "/api/control/get/sms",
			Params: params,
		}

		resp, err := curl.GET()
		if err != nil {
			fmt.Println("auto sms error", err.Error())
			continue
		}

		var s controllers.ResultGetSms
		_ = json.Unmarshal([]byte(resp), &s)
		if s.Code != 200 {
			fmt.Println("auto sms error")
			continue
		}

		// 修改短信信息 TODO 没有做异常处理
		//table := models.SendPhoneNumberList{
		//	RequestId: j.RequestId,
		//	SmsCode:   s.Data,
		//}
		//table.UpdateSmsSendSuccessStatus()
	}
}

// setToken 设置 bus token, 解决请求 API 鉴权的问题
func setToken(token string) *url.Values {
	params := url.Values{}
	params.Set("token", token)
	return &params
}
