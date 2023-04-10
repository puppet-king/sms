// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package sms
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"io"
	"os"
	"sms/config"
	"sms/console"
	api_v1 "sms/controllers"
	"sms/models"
	"time"
)

// SmsEnvMode 设置环境变量
const SmsEnvMode = "SMS_ENV_MODE"

func init() {
	// 判断环节变量
	if gin.ReleaseMode == os.Getenv(SmsEnvMode) {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig("private_config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(20)
	}

	// 加载 wechat
	models.InitWechat()

	db, err := models.InitDBConnectionPool(cfg.Database)
	// 连接数据库
	defer db.Close()

	// 开启 job 协程  并且是正式环境
	if gin.ReleaseMode == os.Getenv(SmsEnvMode) {
		go cronJob()
	}

	// 加载缓存
	cache := models.NewCache()
	cache.SetAllowList()

	// web 服务
	web()

}

// web WEB 服务
func web() {
	// 创建路由
	r := gin.New()
	// 使用中间件来将 配置文件、 db 对象保存到 context.Context 对象中
	r.Use(func(c *gin.Context) {
		c.Set("privateConfig", config.Cfg)
		c.Set("db", models.DB)
		c.Next()
	})

	// Logging to a file.
	f, _ := os.OpenFile("./log/info.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	var conf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,请求状态码:%d,响应时间:%s,错误信息:%s, \n",
				param.ClientIP,
				param.TimeStamp.Format("2006年01月02日 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				//param.Request.UserAgent(),
				param.ErrorMessage,
				//param.Request.Response.Body.Read,
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
		v1.POST("/login", api_v1.Login)
		v1.POST("/token-login", api_v1.TokenLogin)
		v1.GET("/get-balance", api_v1.GetBalance)
		v1.GET("/get-all-project", api_v1.GetAllProject)
		v1.GET("/get-sms", api_v1.GetSms)
		v1.GET("/cancel-request", api_v1.CancelRequest)
		v1.GET("/get-all-countries", api_v1.GetAllCountries)
		v1.GET("/get-phone-number", api_v1.GetPhoneNumber)
		v1.GET("/get-available-numbers", api_v1.GetAvailableNumbers)
		v1.GET("/get-available-numbers-by-uid", api_v1.GetAvailableNumbersByUid)

		// tool
		toolGroup := v1.Group("/tool")
		{
			t := new(api_v1.ToolController)
			toolGroup.GET("/check-db", t.CheckDb)
			toolGroup.GET("/set-cache", t.SetCache)
			toolGroup.GET("/get-cache", t.GetCache)
			toolGroup.GET("/get-last-login-info", t.GetLastLoginInfo)
		}
	}

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":9090")
}

// cron 定时服务
func cronJob() {
	c := cron.New()
	defer c.Stop()
	// 每秒检查一次
	spec := "*/5 * * * * ?"
	err := c.AddFunc(spec, func() {
		s := new(console.Sms)
		s.ExecTime = time.Now()
		s.AutoCancel()
		s.AutoSmsCode()
		fmt.Println("执行结束", "耗时:", time.Since(s.ExecTime).Seconds())
	})
	if err != nil {
		fmt.Errorf("AddFunc error : %v", err)
		return
	}
	c.Start()
	select {}
}
