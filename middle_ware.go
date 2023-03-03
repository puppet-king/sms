// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package main  中间键服务
*/
package main

import (
	"github.com/gin-gonic/gin"
	"sms/config"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if projectToken := c.MustGet("privateConfig").(*config.PrivateConfig).ProjectToken;
			projectToken != c.GetHeader("Token") {
			c.Abort()
			return
		}
	}
}
