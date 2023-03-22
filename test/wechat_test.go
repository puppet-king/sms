// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package test

import (
	"sms/models"
	"testing"
)

func TestCode2Session(t *testing.T) {
	initDb()

	type args struct {
		jsCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"base",
			args{jsCode: "031JkPGa1lxmVE0HElFa1Jd1ne4JkPGc"},
			"opz1q5VY9-g3NbEGCaverijyU_TU",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := models.Code2Session(tt.args.jsCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Code2Session() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Code2Session() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestInitWechat(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.InitWechat()
		})
	}
}
