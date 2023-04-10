// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package test

import (
	"fmt"
	"sms/models"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	initDb()
	tests := []struct {
		name string
	}{
		{"base"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.Cache{}
			//got.Set("wang", "value", 5*time.Second)
			//time.Sleep(1 * time.Second)
			//got.Set("wang", "value2", 100*time.Second)
			//
			//fmt.Println(got.Get("wang"))
			//
			//time.Sleep(6 * time.Second)
			//fmt.Println(got.Get("wang"))
			////got.Set("wang", "value2", 0)
			cache := models.Cache{}
			cache.SetAllowList()
			loginList := cache.GetLastLoginInfo()
			fmt.Println(loginList)
			loginList["opz1q5VY9-g3NbEGCaverijyU_TU"] = time.Now().Format("2006-01-02 15:04:05")

			got.Set("sms:user:loginInfo", loginList, 0)
			fmt.Println(got.GetLastLoginInfo())
		})
	}
}
