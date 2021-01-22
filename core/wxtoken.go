package core

import (
	"encoding/json"
	"log"
	"wechat-token-server/public"
)

//微信请求token的url 正常返回{"access_token":"ACCESS_TOKEN","expires_in":7200}
// 错误时返回 {"errcode":40013,"errmsg":"invalid appid"}
const ACCESS_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token"

//token信息
type WxToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expires_in"` //过期时间 单位：秒。目前是7200秒之内的值。
}

// 获取AppID的access_token
func (wt *WxToken) Get(appid string, secret string) string {
	var params = map[string]string{
		"appid":      appid,
		"secret":     secret,
		"grant_type": "client_credential",
	}

	token := &WxToken{} //临时变量 用于判断
	at := public.GetAccessToken(ACCESS_TOKEN_URL, params)
	if err := json.Unmarshal([]byte(at), token); err != nil {
		log.Printf("json Unmarshal error : %v \n", err)
		return ""
	}
	if token.AccessToken == "" {
		//获取失败
		log.Printf("[ERROR] get wxtoken fail errmsg is : %s\n", token.ErrMsg)
		wt.ErrCode = token.ErrCode
		wt.ErrMsg = token.ErrMsg
		wt.AccessToken = ""
		wt.Expire = 0
		return ""
	}
	//获取成功
	wt.ErrCode = 0
	wt.ErrMsg = ""
	wt.AccessToken = token.AccessToken
	wt.Expire = token.Expire
	return wt.AccessToken
}
