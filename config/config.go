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
}

func LoadConfig(name string) (Cfg *PrivateConfig, err error) {
	// 读取配置文件
	content, err := ini.Load(name)
	if err != nil {
		return nil, err
	}

	cfg := &PrivateConfig{
		ApiKey:       content.Section("sms").Key("apiKey").String(),
		ProjectToken: content.Section("project").Key("token").String(),
		Database:     content.Section("database").Key("dataSourceName").String(),
	}

	return cfg, nil
}
