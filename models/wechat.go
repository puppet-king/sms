// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"github.com/medivhzhan/weapp/v3"
	"github.com/medivhzhan/weapp/v3/auth"
	"github.com/medivhzhan/weapp/v3/logger"
	"sms/config"
)

var WechatSdk *auth.Auth

func InitWechat() {
	sdk := weapp.NewClient(config.Cfg.Wechat.AppId, config.Cfg.Wechat.AppSecret)
	// 任意切换日志等级
	sdk.SetLogLevel(logger.Silent)
	cli := sdk.NewAuth()
	WechatSdk = cli
}

// Code2Session code 兑换 session
func Code2Session(jsCode string) (string, error) {
	Code2SessionResponse, err := WechatSdk.Code2Session(&auth.Code2SessionRequest{
		Appid:     config.Cfg.Wechat.AppId,
		Secret:    config.Cfg.Wechat.AppSecret,
		JsCode:    jsCode,
		GrantType: "authorization_code",
	})

	if err != nil {
		return "", err

	}

	return Code2SessionResponse.Openid, nil
}
