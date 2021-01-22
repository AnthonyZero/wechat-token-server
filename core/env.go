package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	"wechat-token-server/public"

	"github.com/tidwall/buntdb"
)

type Env struct {
	Data map[string]string //appid secret的键值对
	Conf string            //配置文件
	DB   *buntdb.DB
	Wt   *WxToken
}

type Account struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret"`
}

// GetAccounts 读取配置文件中的appid和secret值到一个map中
func (e *Env) GetAccounts(file string) {
	accounts := make([]Account, 1)
	if file == "" {
		e.Conf = "wechat.json"
	} else {
		e.Conf = file
	}

	if _, err := os.Stat(e.Conf); err != nil {
		log.Fatalln(err)
		//os.Exit(1)
	}

	raw, err := ioutil.ReadFile(e.Conf)
	if err != nil {
		log.Fatalln(err)
		//os.Exit(1)
	}

	if err := json.Unmarshal(raw, &accounts); err != nil {
		log.Fatalln(err)
		//os.Exit(1)
	}

	for _, a := range accounts {
		e.Data[a.AppID] = a.Secret
	}
}

// GetValue 通过db中以appid_key的前缀 来获取value
func (e *Env) GetValue(appid string, key string) string {
	var value string

	err := e.DB.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(appid + "_" + key)
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		value = ""
	}

	return value
}

// UpdateTokens 更新AppID上下文环境中的Access Token和到期时间
func (e *Env) UpdateTokens(appid string) {
	timestamp := time.Now().Unix()

	e.DB.Update(func(tx *buntdb.Tx) error {
		tx.Delete(appid + "_" + public.KEY_RECORD_TIME)
		tx.Delete(appid + "_" + public.KEY_ACCESS_TOKEN)
		tx.Delete(appid + "_" + public.KEY_EXPIRES_IN)

		tx.Set(appid+"_"+public.KEY_RECORD_TIME, strconv.FormatInt(timestamp, 10), nil)
		tx.Set(appid+"_"+public.KEY_ACCESS_TOKEN, e.Wt.AccessToken, nil)
		tx.Set(appid+"_"+public.KEY_EXPIRES_IN, strconv.Itoa(e.Wt.Expire), nil)
		return nil
	})
}
