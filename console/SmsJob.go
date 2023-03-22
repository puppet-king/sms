// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package console
*/
package console

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sms/config"
	"sms/controllers"
	"sms/models"
	"time"
)

const HOST = "https://sms-bus.com"

type Sms struct {
	ExecTime time.Time // 执行时间
}

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
		if time.Since(createTime).Minutes() <= 4 {
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

// setToken 设置 bus token, 解决请求 API 鉴权的问题
func setToken(token string) *url.Values {
	params := url.Values{}
	params.Set("token", token)
	return &params
}
