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
	tests := []struct {
		name string
	}{
		{"base"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.NewCache()
			got.Set("wang", "value", 5*time.Second)
			time.Sleep(1 * time.Second)
			got.Set("wang", "value2", 100*time.Second)

			fmt.Println(got.Get("wang"))

			time.Sleep(6 * time.Second)
			fmt.Println(got.Get("wang"))
			//got.Set("wang", "value2", 0)

		})
	}
}
