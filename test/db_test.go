// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package test

import (
	"fmt"
	"reflect"
	"sms/config"
	"sms/models"
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
			got, _ := models.GetLastActivePhoneNumber(tt.args.projectId)
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
			s := &models.SendPhoneNumberList{
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
			RequestId: "2303031019567210983611", // 唯一索引重复测试会冲突
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

			s := &models.SendPhoneNumberList{
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
			s := &models.SendPhoneNumberList{
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
	_, _ = models.InitDBConnectionPool(cfg.Database)
	models.InitWechat()
}

func TestSendPhoneNumberList_GetListByStatus(t *testing.T) {
	initDb()
	type fields struct {
		RequestId string
		ProjectId string
		AreaCode  string
		Number    string
		Status    int
		CancelAt  string
		SmsCode   string
	}
	type args struct {
		status int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.SendPhoneNumberList
		wantErr bool
	}{
		{"base", fields{
			RequestId: "",
			ProjectId: "",
			AreaCode:  "",
			Number:    "",
			Status:    0,
			CancelAt:  "",
			SmsCode:   "",
		}, args{status: 2}, []models.SendPhoneNumberList{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.SendPhoneNumberList{
				RequestId: tt.fields.RequestId,
				ProjectId: tt.fields.ProjectId,
				AreaCode:  tt.fields.AreaCode,
				Number:    tt.fields.Number,
				Status:    tt.fields.Status,
				CancelAt:  tt.fields.CancelAt,
				SmsCode:   tt.fields.SmsCode,
			}
			got, err := s.GetListByStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println(got)
			fmt.Println(reflect.TypeOf(got))

			//if !reflect {
			//	t.Errorf("GetListByStatus() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
