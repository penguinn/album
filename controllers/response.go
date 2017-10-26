package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	SUCCESS_CODE      = 0
	PARAM_ERROR_CODE  = 1000
	SERVER_ERROR_CODE = 1001

	AUTH_FAILED_CODE        = 2001
	NOT_AUTH_CODE           = 2002
	USERNAME_EXIST_CODE     = 2003
	USERNAME_NOT_EXIST_CODE = 2004
	KICKED_CODE             = 2005

	IMAGE_AUTH_CODE = 3001
)

type Response struct {
	Data    interface{} `json:"data"`
	ErrMsg  string      `json:"errmsg"`
	ErrCode int         `json:"errcode"` //0:正确
}

var (
	CodeMap = map[int]string{
		SUCCESS_CODE:            "success",
		PARAM_ERROR_CODE:        "参数错误",
		SERVER_ERROR_CODE:       "系统繁忙，请稍后再试",
		AUTH_FAILED_CODE:        "登录失败",
		NOT_AUTH_CODE:           "没有登录",
		USERNAME_EXIST_CODE:     "用户名已经存在",
		USERNAME_NOT_EXIST_CODE: "用户名不存在",
		KICKED_CODE:             "账户在另一个地方登录",
	}
)

func NewResponse(data interface{}, options ...interface{}) interface{} {
	code := 0
	if len(options) > 0 {
		code, _ = options[0].(int)
	}
	msg := ""
	if len(options) > 1 {

		msg = fmt.Sprintf("%v", options[1])
	}
	if msg == "" {
		msg, _ = CodeMap[code]
	}

	return gin.H{
		"data":    data,
		"errmsg":  msg,
		"errcode": code,
	}
}
