// Copyright 2022 The wangkai. ALL rights reserved.
/*
	思考：随着接口继续写下去, 发现存在 controller 文件下的方法不能相同命令的问题, 导致写 api 命名很冗余,
这个用 method 来实现; 紧接着又遇到需要声明各种结构的问题, 如果放置在 function 里面就不具备可复用,可如果不放,
看起来就很混乱(还是我奇怪？), 这个时候需要想下复用的可能性, API 往往应该做成一个 service 层然后把对应 API 放置进去 ? 不太行因为
万一要用里面的结构体就麻烦了，所以我觉得最好定义一个标准解决看起来很混乱的问题，例如 入参 出参标准化 [API | ][req|res]functionName

*/

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
	"time"
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
	"openid":                            true,
	"opz1q5VY9-g3NbEGCaverijyU_TU":      true,
	"opz1q5Vn8dwFsh2zrxU6s8bQwfwY":      true,
	"opz1q5esfV55VYlolpooTk-sNYjw":      true,
	"tokenOpz1q5esfV55VYlolpooTk-sNYjw": true,
	"opz1q5QDbXHxc6UKOA30bzBGymX8":      true,
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
	if err != nil || openId == "" {
		result["code"] = http.StatusInternalServerError
		result["msg"] = "无效数据"
		c.JSON(http.StatusForbidden, result)
		return
	}

	//writeFile("openid.txt", openId+"\r\n")
	cache := models.NewCache()
	loginList := cache.GetLastLoginInfo()
	loginList[openId] = time.Now().Format("2006-01-02 15:04:05")
	cache.Set("sms:user:loginInfo", loginList, 0)

	user := models.User{OpenId: openId}

	// 过滤无效用户列表
	allowOpenidList := cache.GetAllowList()
	if _, ok := allowOpenidList[user.OpenId]; !ok {
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

// TokenRequest 入参结构体解析
type TokenRequest struct {
	Token string `json:"token"`
}

// TokenLogin Token 登录
func TokenLogin(c *gin.Context) {
	result := DefaultResult
	// 获取的是 code
	loginRequest := TokenRequest{}
	err := c.BindJSON(&loginRequest)
	if err != nil || loginRequest.Token == "" {
		result["code"] = http.StatusInternalServerError
		result["msg"] = "服务器异常"
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	user := models.User{OpenId: loginRequest.Token}

	// 过滤无效用户列表
	cache := models.NewCache()
	allowOpenidList := cache.GetAllowList()
	if _, ok := allowOpenidList[user.OpenId]; !ok {
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
	params.Set("project_id", c.DefaultQuery("project_id", strconv.Itoa(ProjectId)))

	// 国家在接口没有传递时读取 DB 配置 （大多数情况）
	countryId := c.Query("country_id")
	if countryId == "" {
		projectId, _ := strconv.Atoi(c.DefaultQuery("project_id", strconv.Itoa(ProjectId)))
		if defaultCountryId, ok := models.GetDefaultCountryId(projectId); ok {
			countryId = strconv.Itoa(defaultCountryId)
		} else {
			countryId = strconv.Itoa(CountryId)
		}
	}
	params.Set("country_id", countryId)

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/get/number",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	//resp := "{\"code\":200,\"message\":\"Operation Success\",\"data\":{\"request_id\":\"230303101956721098368\",\"number\":\"12897129788\",\"area_code\":\"1\"}}"
	s := ResultGetPhoneNumber{}
	_ = json.Unmarshal([]byte(resp), &s)
	if s.Code != 200 {
		c.String(http.StatusInternalServerError, s.Message)
		return
	}

	SendPhoneNumberList := models.SendPhoneNumberList{
		RequestId: s.Data.RequestId,
		UserId:    c.MustGet("userId").(string),
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
			c.String(http.StatusBadRequest, "requestId 参数缺失")
			return
		}
		requestId = row.RequestId
	}

	params.Set("request_id", requestId)

	// 先从数据库读取是否已存在
	table := models.SendPhoneNumberList{
		RequestId: requestId,
	}

	info, err := table.GetInfoByRequestId()
	if err == nil && info.Status > 0 {
		if info.Status == 1 {
			//{"code":200,"message":"Operation Success","data":"077866"}
			success := ResultGetSms{
				Code:    200,
				Message: "Operation Success",
				Data:    info.SmsCode,
			}
			c.JSON(http.StatusOK, success)
			return
		} else {
			c.String(http.StatusBadRequest, "Number has been released or timeout, please reacquire")
			return
		}
	}

	curl := models.BaseCurl{
		Host:   HOST,
		Path:   "/api/control/get/sms",
		Params: params,
	}

	resp, err := curl.GET()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var s ResultGetSms
	_ = json.Unmarshal([]byte(resp), &s)
	if s.Code != 200 {
		c.String(http.StatusInternalServerError, s.Message)
		return
	}

	// 修改短信信息 TODO 没有做异常处理
	table.SmsCode = s.Data
	table.UpdateSmsSendSuccessStatus()

	c.String(http.StatusOK, resp)
}

// GetCacheSms 获取缓存手机短信验证码
func GetCacheSms(c *gin.Context) {
	params := setToken(c.MustGet("privateConfig").(*config.PrivateConfig).ApiKey)
	requestId := c.DefaultQuery("request_id", "0")
	if requestId == "0" {
		// 获取默认配置下的最新一条数据
		row, err := models.GetLastActivePhoneNumber(ProjectId)
		if err != nil || row.RequestId == "" {
			c.String(http.StatusBadRequest, "requestId 参数缺失")
			return
		}
		requestId = row.RequestId
	}

	params.Set("request_id", requestId)

	// 先从数据库读取是否已存在
	table := models.SendPhoneNumberList{
		RequestId: requestId,
	}

	info, err := table.GetInfoByRequestId()
	if err == nil && info.Status > 0 {
		if info.Status == 1 {
			//{"code":200,"message":"Operation Success","data":"077866"}
			success := ResultGetSms{
				Code:    200,
				Message: "Operation Success",
				Data:    info.SmsCode,
			}
			c.JSON(http.StatusOK, success)
			return
		} else {
			c.String(http.StatusBadRequest, "Number has been released or timeout, please reacquire")
			return
		}
	}

	c.String(http.StatusBadRequest, "not cache")
}

// ResultCancelRequest 返回 CancelRequest 请求 API 结构体
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
			c.String(http.StatusBadRequest, "requestId 参数缺失")
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var s ResultCancelRequest
	_ = json.Unmarshal([]byte(resp), &s)
	// 这个 返回有待商榷
	if s.Code == 50103 {
		c.String(http.StatusBadRequest, requestId+"已经过期无法使用")
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
		c.String(http.StatusInternalServerError, err.Error())
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
		c.String(http.StatusInternalServerError, err.Error())
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusBadRequest, resp)
}

// AvailableCountries 可用国家返回结构子集
type AvailableCountries struct {
	CountryId  int     `json:"country_id"`
	ProjectId  int     `json:"project_id"`
	Cost       float64 `json:"cost"`
	TotalCount int     `json:"total_count"`
}

// AvailableCountriesData 可用国家返回结构
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var s AvailableCountriesData
	_ = json.Unmarshal([]byte(resp), &s)
	data := AvailableCountries{}

	projectId, _ := strconv.Atoi(c.DefaultQuery("project_id", strconv.Itoa(ProjectId)))
	for k, v := range s.Data {
		if v.ProjectId == projectId {
			data = s.Data[k]
		}
	}

	c.JSON(http.StatusBadRequest, data)
}

type ReqGetAvailableNumbersByUid struct {
	models.SendPhoneNumberList
	RefreshTime string
}

// GetAvailableNumbersByUid 根据 uid 获取可用号码
func GetAvailableNumbersByUid(c *gin.Context) {
	result := DefaultResult

	// 获取数据
	sms := models.SendPhoneNumberList{}
	list, err := sms.GetListByUid(c.MustGet("userId").(string))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 处理数据
	// 更改时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("时区设置不正确：", err)
		return
	}

	var res []ReqGetAvailableNumbersByUid
	for _, v := range list {
		createTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.CreateAt, loc)
		if err != nil {
			continue
		}

		var refreshTime = ""
		seconds := time.Since(createTime).Seconds()
		if seconds < 60 {
			refreshTime = strconv.Itoa(int(seconds)) + "秒前"
		} else {
			refreshTime = strconv.Itoa(int(seconds/60)) + "分钟前"
		}

		res = append(res, ReqGetAvailableNumbersByUid{
			v,
			refreshTime,
		})
	}

	result["code"] = 200
	result["msg"] = "成功"
	result["data"] = res
	c.JSON(http.StatusOK, result)
	return
}

// setToken 设置 bus token, 解决请求 API 鉴权的问题
func setToken(token string) *url.Values {
	params := url.Values{}
	params.Set("token", token)
	return &params
}

// writeFile 写入文件, 解决记录 openid 做白名单的问题
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
