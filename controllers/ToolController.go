// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package controllers
*/
package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sms/models"
)

type ToolController struct {
}

// CheckDb 检查 DB
func (t ToolController) CheckDb(c *gin.Context) {
	result := DefaultResult
	db := models.DB

	dbMsg := fmt.Sprintf("最大连接数：%d,  当前总连接数；%d,  "+
		"已使用: %d, 空闲数量：%d \n", db.Stats().MaxOpenConnections, db.Stats().OpenConnections,
		db.Stats().InUse, db.Stats().Idle) + "\n" + fmt.Sprintf("数量指标 :) \n等待连接数量；%d,  等待创建新连接时长(秒): %f, 空闲超限关闭数量：%d, 空闲超时关闭数量：%d, 连接超时关闭数量：%d \n",
		db.Stats().WaitCount,
		db.Stats().WaitDuration.Seconds(),
		db.Stats().MaxIdleClosed,
		db.Stats().MaxIdleTimeClosed,
		db.Stats().MaxLifetimeClosed,
	)

	result["data"] = dbMsg
	c.JSON(http.StatusOK, result)
}

// SetCache 设置缓存
func (t ToolController) SetCache(c *gin.Context) {
	result := DefaultResult

	cache := models.NewCache()
	cache.SetAllowList()
	result["code"] = 200
	result["msg"] = "成功"
	c.JSON(http.StatusOK, result)
}

// GetCache 获取缓存
func (t ToolController) GetCache(c *gin.Context) {
	result := DefaultResult
	cache := models.NewCache()

	if data, ok := cache.Get("sms:user:allowList"); ok {
		if allowOpenidList, ok := data.(map[string]bool); ok {
			result["data"] = allowOpenidList
			c.JSON(http.StatusOK, result)
			return
		}
	}

	result["code"] = 50001
	result["msg"] = "暂无数据"
	c.JSON(http.StatusOK, result)
}

// GetLastLoginInfo 获取登录人员 openid
func (t ToolController) GetLastLoginInfo(c *gin.Context) {
	result := DefaultResult
	cache := models.NewCache()

	if data, ok := cache.Get("sms:user:loginInfo"); ok {
		if allowOpenidList, ok := data.(map[string]string); ok {
			result["data"] = allowOpenidList
			c.JSON(http.StatusOK, result)
			return
		}
	}

	result["code"] = 50001
	result["msg"] = "暂无数据"
	c.JSON(http.StatusOK, result)
}
