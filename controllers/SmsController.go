// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package controllers
*/
package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"sms/config"
	"sms/models"
	"strconv"
)

const HOST = "https://sms-bus.com"
const ProjectId = 52 // openai
const CountryId = 5  // 美国

type ResultGetPhoneNumber struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RequestId string `json:"request_id"`
		Number    string `json:"number"`
		AreaCode  string `json:"area_code"`
	} `json:"data"`
}

var DefaultResult = gin.H{
	"code": 200,
	"msg":  "success",
	"data": map[string]any{},
}

var AllowOpenidList = map[string]bool{
	"openid":                       true,
	"opz1q5VY9-g3NbEGCaverijyU_TU": true,
	"opz1q5Vn8dwFsh2zrxU6s8bQwfwY": true,
	"opz1q5esfV55VYlolpooTk-sNYjw": true,
	//"031qsd0w3rHBh03iSN3w3kjjnF1qsd0m": true,
}

type LoginRequest struct {
	Code string `json:"code"`
}

func Login(c *gin.Context) {
	result := DefaultResult
	// 获取的是 code
	loginRequest := LoginRequest{}
	err := c.BindJSON(&loginRequest)
	if err != nil || loginRequest.Code == "" {
		result["code"] = http.StatusInternalServerError
		result["msg"] = "服务器异常"
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	// 兑换 openid
	openId, err := models.Code2Session(loginRequest.Code)
	writeFile("openid.txt", openId+"\r\n")

	if err != nil || openId == "" {
		result["code"] = http.StatusInternalServerError
		result["msg"] = "无效数据"
		c.JSON(http.StatusForbidden, result)
		return
	}

	user := models.User{OpenId: openId}
	// 过滤无效用户列表
	if _, ok := AllowOpenidList[user.OpenId]; !ok {
		result["code"] = http.StatusInternalServerError
		result["msg"] = err.Error()
		c.JSON(http.StatusForbidden, result)
		return
	}

	// 生成 token
	token, err := models.GetToken(user)
	if err != nil {
		result["code"] = http.StatusInternalServerError
		result["msg"] = err.Error()
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	result["data"] = map[string]string{
		"token":   token,
		"open_id": user.OpenId,
	}
	result["code"] = 200
	result["msg"] = "成功"
	c.JSON(http.StatusOK, result)
}

// GetPhoneNumber 获取手机号码
func GetPhoneNumber(c *gin.Context) {
	params := setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey)
	params.Set("country_id", c.DefaultQuery("country_id", strconv.Itoa(CountryId)))
	params.Set("project_id", c.DefaultQuery("project_id", strconv.Itoa(ProjectId)))

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/get/number",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//resp := "{\"code\":200,\"message\":\"Operation Success\",\"data\":{\"request_id\":\"230303101956721098368\",\"number\":\"12897129788\",\"area_code\":\"1\"}}"
	s := ResultGetPhoneNumber{}
	_ = json.Unmarshal([]byte(resp), &s)
	SendPhoneNumberList := models.SendPhoneNumberList{
		RequestId: s.Data.RequestId,
		ProjectId: c.DefaultQuery("project_id", strconv.Itoa(ProjectId)),
		AreaCode:  c.DefaultQuery("country_id", strconv.Itoa(CountryId)),
		Number:    s.Data.Number,
		Status:    0,
		CancelAt:  "",
		SmsCode:   "",
	}

	_, err = SendPhoneNumberList.Insert()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, resp)
}

// ResultGetSms GetSms 接口的返回参数
type ResultGetSms struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// GetSms 获取手机短信验证码
func GetSms(c *gin.Context) {
	params := setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey)
	requestId := c.DefaultQuery("request_id", "0")
	if requestId == "0" {
		// 获取默认配置下的最新一条数据
		row, err := models.GetLastActivePhoneNumber(ProjectId)
		if err != nil || row.RequestId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "requestId 参数缺失"})
			return
		}
		requestId = row.RequestId
	}

	params.Set("request_id", requestId)

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/get/sms",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var s ResultGetSms
	_ = json.Unmarshal([]byte(resp), &s)

	// 修改短信信息 TODO 没有做异常处理
	table := models.SendPhoneNumberList{
		RequestId: requestId,
		SmsCode:   s.Data,
	}
	table.UpdateSmsSendSuccessStatus()

	c.String(http.StatusOK, resp)
}

type ResultCancelRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CancelRequest 取消发送
func CancelRequest(c *gin.Context) {
	params := setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey)
	requestId := c.DefaultQuery("request_id", "0")
	if requestId == "0" {
		// 获取默认配置下的最新一条数据
		row, err := models.GetLastActivePhoneNumber(ProjectId)
		if err != nil || row.RequestId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "requestId 参数缺失"})
			return
		}
		requestId = row.RequestId
	}
	params.Set("request_id", requestId)

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/cancel",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s ResultCancelRequest
	_ = json.Unmarshal([]byte(resp), &s)
	if s.Code == 50103 {
		c.JSON(http.StatusBadRequest, gin.H{"error": requestId + "已经过期无法使用"})
		return
	}

	// 修改短信信息
	table := models.SendPhoneNumberList{
		RequestId: requestId,
	}
	table.CancelSmsSend(requestId)

	c.String(http.StatusOK, resp)
}

// GetBalance 获取余额
func GetBalance(c *gin.Context) {
	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/get/balance",
		Params: setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey),
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, resp)
}

// GetAllProject 获取所有项目
func GetAllProject(c *gin.Context) {
	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/list/projects",
		Params: setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey),
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, resp)
}

// GetAllCountries 获取所有国家
func GetAllCountries(c *gin.Context) {
	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/list/countries",
		Params: setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey),
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusBadRequest, resp)
}

type AvailableCountries struct {
	CountryId  int     `json:"country_id"`
	ProjectId  int     `json:"project_id"`
	Cost       float64 `json:"cost"`
	TotalCount int     `json:"total_count"`
}

type AvailableCountriesData struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    map[string]AvailableCountries `json:"data"`
}

// GetAvailableNumbers 获取所有国家根据项目的可用数量
func GetAvailableNumbers(c *gin.Context) {
	params := setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey)
	params.Set("country_id", c.DefaultQuery("country_id", strconv.Itoa(CountryId)))
	params.Set("project_id", c.DefaultQuery("project_id", strconv.Itoa(ProjectId)))

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/list/prices",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var s AvailableCountriesData
	_ = json.Unmarshal([]byte(resp), &s)
	data := AvailableCountries{}
	for k, v := range s.Data {
		if v.ProjectId == ProjectId {
			data = s.Data[k]
		}
	}

	c.JSON(http.StatusBadRequest, data)
}

func setToken(token string) *url.Values {
	params := url.Values{}
	params.Set("token", token)
	return &params
}

func writeFile(fileName string, content string) bool {
	file, err := os.OpenFile("./log/"+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件无法打开", err)
		return false
	}

	// *及时关闭 file 句柄, 会在函数执行结束后回调该方法
	defer func(file *os.File) {
		//fmt.Println("执行结束回调 defer")
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭失败, 请及时处理", err)
		}
	}(file)

	// 写入文件
	write := bufio.NewWriter(file)
	_, err = write.WriteString(content)
	err = write.Flush()
	if err != nil {
		fmt.Println("文件写入失败", err)
		return false
	}

	return true
}
