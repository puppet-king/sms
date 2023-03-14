// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package config
*/
package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		name string
	}

	config := &PrivateConfig{
		"1",
		"2",
		"root:1@/db?charset=utf8",
		Wechat{
			AppId:     "",
			AppSecret: "",
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *PrivateConfig
		wantErr bool
	}{
		{"error_file", args{"1234"}, config, true},
		{"base", args{"private_config_test.ini"}, config, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.name)
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
