// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sms/config"
	"strings"
	"time"
)

type MyCustomClaims struct {
	User
	jwt.RegisteredClaims
}

type User struct {
	OpenId string `json:"open_id"`
}

// GetToken 获取 token
func GetToken(user User) (string, error) {
	// Create the claims
	claims := MyCustomClaims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sms",
			Subject:   "smsLogin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Cfg.ProjectToken))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// TokenVia 校验 token 并且返回对应用户唯一标识 (当前是 openid)
func TokenVia(tokenString string) (bool, string) {
	tokenString = strings.Split(tokenString, "Bearer ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.ProjectToken), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.OpenId, claims.RegisteredClaims.Issuer)
		return true, claims.OpenId
	} else {
		fmt.Println(err)
		return false, ""
	}
}
