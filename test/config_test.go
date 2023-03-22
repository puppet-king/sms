// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package config
*/
package test

import (
	"reflect"
	config2 "sms/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		name string
	}

	config := &config2.PrivateConfig{
		"1",
		"2",
		"root:1@/db?charset=utf8",
		config2.Wechat{
			AppId:     "",
			AppSecret: "",
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *config2.PrivateConfig
		wantErr bool
	}{
		{"error_file", args{"1234"}, config, true},
		{"base", args{"private_config_test.ini"}, config, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config2.LoadConfig(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
