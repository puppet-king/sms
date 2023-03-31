// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package main  中间键服务
*/
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sms/controllers"
	"sms/models"
)

var AllowList = map[string]bool{
	"/v1/login":       true,
	"/v1/token-login": true,
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := AllowList[c.Request.URL.Path]; ok {
			return
		}

		// 初始化返回参数
		controllers.DefaultResult = gin.H{
			"code": 200,
			"msg":  "success",
			"data": map[string]any{},
		}

		// 校验 token 是否传递
		result := controllers.DefaultResult
		if token := c.GetHeader("Authorization"); token != "" {
			if ok, openId := models.TokenVia(token); ok {
				c.Set("userId", openId)
			} else {
				result["code"] = http.StatusBadRequest
				result["msg"] = "鉴权失败-1"
				c.JSON(http.StatusForbidden, result)
				c.Abort()
				return
			}
		} else {
			result["code"] = http.StatusBadRequest
			result["msg"] = "鉴权失败-null"
			c.JSON(http.StatusForbidden, result)
			c.Abort()
			return
		}
	}
}
