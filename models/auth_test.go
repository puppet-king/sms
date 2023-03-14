// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package controllers
*/
package models

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	initDb()
	type args struct {
		user User
	}
	tests := []struct {
		name string
		args args
	}{
		{"base", args{User{
			OpenId: "openid",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetToken(tt.args.user)
		})
	}
}

func TestTokenVia(t *testing.T) {
	initDb()

	type args struct {
		tokenString string
	}
	tests := []struct {
		name string
		args args
	}{
		{"base",
			args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJPcGVuSWQiOiJvcGVuaWQiLCJpc3MiOiJzbXMiLCJzdWIiOiJzbXNMb2dpbiIsImV4cCI6MTY3ODc3MjQxNiwibmJmIjoxNjc4Njg2MDE2LCJpYXQiOjE2Nzg2ODYwMTZ9.pWYmg2eQUfHP24vhtYMLyIfuGo7wCgsjqec4OVdEJPI"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TokenVia(tt.args.tokenString)
		})
	}
}
