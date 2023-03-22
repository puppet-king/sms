// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package console
*/
package test

import (
	"sms/console"
	"testing"
	"time"
)

func TestSms_autoCancel(t *testing.T) {
	initDb()

	type fields struct {
		ExecTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"base", fields{ExecTime: time.Now()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := console.Sms{
				ExecTime: tt.fields.ExecTime,
			}
			s.AutoCancel()
		})
	}
}
