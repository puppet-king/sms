// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package config
*/
package config

import (
	"gopkg.in/ini.v1"
)

type PrivateConfig struct {
	ApiKey       string
	ProjectToken string
	Database     string
	Wechat
}

type Wechat struct {
	AppId     string
	AppSecret string
}

var Cfg *PrivateConfig

func LoadConfig(name string) (cfg *PrivateConfig, err error) {
	// 读取配置文件
	content, err := ini.Load(name)
	if err != nil {
		return nil, err
	}

	cfg = &PrivateConfig{
		ApiKey:       content.Section("sms").Key("apiKey").String(),
		ProjectToken: content.Section("project").Key("token").String(),
		Database:     content.Section("database").Key("dataSourceName").String(),
		Wechat: Wechat{
			AppId:     content.Section("wechat").Key("appId").String(),
			AppSecret: content.Section("wechat").Key("appSecret").String(),
		},
	}

	Cfg = cfg
	return cfg, nil
}
