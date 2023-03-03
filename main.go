// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package sms
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"sms/config"
	api_v1 "sms/controllers"
	"sms/models"
)

type PrivateConfig struct {
	apiKey       string
	projectToken string
	database     string
}

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("private_config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(20)
	}

	db, err := models.InitDBConnectionPool(config.Database)
	// 连接数据库
	defer db.Close()

	// 创建路由
	r := gin.New()

	// 使用中间件来将 配置文件、 db 对象保存到 context.Context 对象中
	r.Use(func(c *gin.Context) {
		c.Set("privateConfig", config)
		c.Set("db", db)
		c.Next()
	})

	// Logging to a file.
	f, _ := os.OpenFile("./log/info.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	var conf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,请求状态码:%d,响应时间:%s,客户端:%s，错误信息:%s\n",
				param.ClientIP,
				param.TimeStamp.Format("2006年01月02日 15:03:04"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		//Output: io.MultiWriter(f),
		Output: io.MultiWriter(os.Stdout, f),
	}
	r.Use(gin.LoggerWithConfig(conf))

	// 调用中间件, 用于解决授权
	r.Use(MiddleWare())

	// 创建路由组
	v1 := r.Group("/v1")
	{
		// 获取手机号
		v1.GET("/get-balance", api_v1.GetBalance)
		v1.GET("/get-all-project", api_v1.GetAllProject)
		v1.GET("/get-sms", api_v1.GetSms)
		v1.GET("/cancel-request", api_v1.CancelRequest)
		v1.GET("/get-all-countries", api_v1.GetAllCountries)
		v1.GET("/get-phone-number", api_v1.GetPhoneNumber)
		v1.GET("/get-available-numbers", api_v1.GetAvailableNumbers)
	}

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":9090")
}
