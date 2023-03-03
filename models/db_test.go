// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"sms/config"
	"testing"
)

func TestGetLastActivePhoneNumber(t *testing.T) {
	type args struct {
		projectId int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"openai", args{projectId: 42}, false},
		{"not", args{projectId: -1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDb()
			got, _ := GetLastActivePhoneNumber(tt.args.projectId)
			t.Log(got)
		})
	}
}

func TestSendPhoneNumberList_CancelSmsSend(t *testing.T) {
	type fields struct {
		RequestId string
		ProjectId string
		AreaCode  string
		Number    string
		Status    int
		CancelAt  string
		SmsCode   string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"base", fields{
			RequestId: "230303101956721098368",
			ProjectId: "42",
			AreaCode:  "",
			Number:    "",
			Status:    0,
			CancelAt:  "",
			SmsCode:   "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDb()
			s := &SendPhoneNumberList{
				RequestId: tt.fields.RequestId,
				ProjectId: tt.fields.ProjectId,
				AreaCode:  tt.fields.AreaCode,
				Number:    tt.fields.Number,
				Status:    tt.fields.Status,
				CancelAt:  tt.fields.CancelAt,
				SmsCode:   tt.fields.SmsCode,
			}
			id := s.CancelSmsSend(s.RequestId)
			if id > 0 {
				t.Log(id)
			}
		})
	}
}

func TestSendPhoneNumberList_Insert(t *testing.T) {
	type fields struct {
		RequestId string
		ProjectId string
		AreaCode  string
		Number    string
		Status    int
		CancelAt  string
		SmsCode   string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"insert", fields{
			RequestId: "230303101956721098361", // 唯一索引重复测试会冲突
			ProjectId: "42",
			AreaCode:  "1",
			Number:    "12897129788",
			Status:    0,
			CancelAt:  "",
			SmsCode:   "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDb()

			s := &SendPhoneNumberList{
				RequestId: tt.fields.RequestId,
				ProjectId: tt.fields.ProjectId,
				AreaCode:  tt.fields.AreaCode,
				Number:    tt.fields.Number,
				Status:    tt.fields.Status,
				CancelAt:  tt.fields.CancelAt,
				SmsCode:   tt.fields.SmsCode,
			}
			s.Insert()
		})
	}
}

func TestSendPhoneNumberList_UpdateSmsSendSuccessStatus(t *testing.T) {
	type fields struct {
		RequestId string
		ProjectId string
		AreaCode  string
		Number    string
		Status    int
		CancelAt  string
		SmsCode   string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"insert", fields{
			RequestId: "230303101956721098368",
			ProjectId: "42",
			AreaCode:  "1",
			Number:    "12897129788",
			Status:    0,
			CancelAt:  "",
			SmsCode:   "200",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDb()
			s := &SendPhoneNumberList{
				RequestId: tt.fields.RequestId,
				ProjectId: tt.fields.ProjectId,
				AreaCode:  tt.fields.AreaCode,
				Number:    tt.fields.Number,
				Status:    tt.fields.Status,
				CancelAt:  tt.fields.CancelAt,
				SmsCode:   tt.fields.SmsCode,
			}
			if rowsAffected := s.UpdateSmsSendSuccessStatus(); rowsAffected > 0 {
				t.Log("成功修改： ", rowsAffected)
			}
		})
	}
}

func initDb() {
	cfg, _ := config.LoadConfig("../config/private_config_test.ini")
	_, _ = InitDBConnectionPool(cfg.Database)
}
