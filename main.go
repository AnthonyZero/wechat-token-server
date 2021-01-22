package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"wechat-token-server/core"
	"wechat-token-server/public"

	"github.com/labstack/echo"
	"github.com/tidwall/buntdb"
)

// Failed 错误返回
const Failed = `{"status": "fail"}`

func main() {
	var err error

	env := &core.Env{}
	env.Wt = &core.WxToken{}
	env.Data = make(map[string]string)
	env.GetAccounts("wechat.json")         // 读取配置文件中的微信AppID和AppSecret 到env.Data中
	env.DB, err = buntdb.Open("wechat.db") // 创建一个K/V数据库用来保存access_token
	if err != nil {
		log.Fatal(err)
	}
	defer env.DB.Close()

	e := echo.New()

	e.GET("/token", func(ctx echo.Context) error {
		if !public.IsValidAuth(ctx) {
			return ctx.String(http.StatusUnauthorized, "401 Authorization Required")
		}
		//获取appid参数
		appid := ctx.QueryParam("appid")
		if appid == "" {
			log.Println("[ERROR]: 没有提供AppID参数")
			return ctx.String(http.StatusNotFound, Failed)
		}

		//获取appid对应的secret
		if secret, isExist := env.Data[appid]; isExist {
			//存在对应appid的相关数据

			var accessToken string
			var recordTime string
			var response struct {
				Status      string `json:"status"`
				AccessToken string `json:"access_token"`
				ErrMsg      string `json:"errmsg"`
			}

			// 查询数据库中是否已经存在这个AppID的access_token
			recordTime = env.GetValue(appid, public.KEY_RECORD_TIME)
			accessToken = env.GetValue(appid, public.KEY_ACCESS_TOKEN)

		GetToken:
			// 如果在数据库中不存在这个appid的token就重新获取
			if accessToken == "" {
				tjs := env.Wt.Get(appid, secret)

				// 没获得access_token就返回Failed消息
				if tjs == "" {
					log.Println("[ERROR]: 没有获取到access_token")
					response.Status = "fail"
					response.ErrMsg = env.Wt.ErrMsg
					//return ctx.String(http.StatusNotFound, Failed)
					return ctx.JSON(http.StatusOK, response)
				}

				//获取Token之后更新运行时环境，然后返回access_token
				env.UpdateTokens(appid)

				response.Status = "success"
				response.AccessToken = env.Wt.AccessToken
				return ctx.JSON(http.StatusOK, response)
			}
			goto CheckTime

		CheckTime:
			// 如果数据库中已经存在了Token，就检查过期时间，如果过期了就去GetToken获取
			curTime := time.Now().Unix()

			expireTime, _ := strconv.ParseInt(recordTime, 10, 64)
			timeout, _ := strconv.ParseInt(env.GetValue(appid, public.KEY_EXPIRES_IN), 10, 64)

			if curTime >= expireTime+timeout {
				//token已过期
				accessToken = "" //fix bug 进入GetToken
				goto GetToken
			}

			response.Status = "success"
			response.AccessToken = accessToken
			return ctx.JSON(http.StatusOK, response)
		}
		log.Println("[ERROR]: 配置文件中对应AppID不存在")
		// 如果提交的appid不在配置文件中，就返回Failed消息
		return ctx.String(http.StatusNotFound, Failed)
	})

	e.Start(":8080")
}
