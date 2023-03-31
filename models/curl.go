// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BaseCurl struct {
	Host   string
	Path   string
	Params *url.Values
}

type curlError struct {
	msg string // 错误描述
}

func (e *curlError) Error() string { return e.msg }

func (b BaseCurl) GET() (string, error) {
	u, err := url.ParseRequestURI(b.Host + b.Path)
	if err != nil {
		fmt.Printf("failed, err:%v \n", err)
		return "", err
	}
	//fmt.Println(b.Params)
	u.RawQuery = b.Params.Encode() // URL encode
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// 获取 HTTP 状态码 和 body 内容
	respBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return string(respBody), nil
	} else {
		return "", &curlError{"状态码异常, 并非 200"}
	}
}

func (b BaseCurl) POST(requestUrl string, header string, body io.Reader) (string, error) {
	resp, err := http.Post(b.Host+b.Path, header, body)
	if err != nil {
		return "", err
	}

	// 获取 HTTP 状态码 和 body 内容
	respBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return string(respBody), nil
	} else {
		return "", &curlError{"状态码异常, 并非 200"}
	}
}
