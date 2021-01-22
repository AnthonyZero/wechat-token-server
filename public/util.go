package public

import (
	"log"

	"github.com/labstack/echo"
	"github.com/levigross/grequests"
)

// 定义Basic Auth的用户名和密码用来防止接口被恶意访问
var Auth = map[string]string{
	"user": "token",
	"pass": "anthonyzero",
}

// 封装grequests的Get方法获取字符串内容
// GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func GetAccessToken(url string, params map[string]string) string {
	ro := &grequests.RequestOptions{
		Params: params,
	}
	response, err := grequests.Get(url, ro)
	if err != nil {
		log.Panicf("[ERROR] grequests.Get err:%v\n", err)
	}
	log.Printf("[SUCCESS] get wxtoken info: %s\n", response.String())
	return response.String()
}

// 检查request的Basic Auth用户名和密码
func IsValidAuth(ctx echo.Context) bool {
	ctx.Response().Header().Set("WWW-Authenticate", `Basic realm="unixs.org"`)
	if uname, upass, ok := ctx.Request().BasicAuth(); ok {
		if Auth["user"] == uname && Auth["pass"] == upass {
			return ok
		}
	}
	return false
}
