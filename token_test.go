package main

import (
	"log"
	"testing"
	"wechat-token/core"
)

//一个测试函数是以Test为函数名前缀的函数
func TestWxToken(t *testing.T) {
	token := &core.WxToken{}
	accessToken := token.Get("appid", "appsecret")
	log.Println(accessToken)
}
