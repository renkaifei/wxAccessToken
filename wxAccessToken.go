package main

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

type JsapiTicket struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

func main() {
	for {
		resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=&secret=")
		if err != nil {
			log.Fatal(err)
		}
		arrbyte, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		var result AccessToken
		log.Println(string(arrbyte))
		err = json.Unmarshal(arrbyte, &result)
		if err != nil {
			log.Fatal(err)
		}
		SetWxKey("wx_access_token", result.Token)
		err = FetchJsapiticket(result.Token)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1000 * time.Millisecond * 60 * 90)
	}
}

func FetchJsapiticket(accessToken string) error {
	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi")
	if err != nil {
		return err
	}
	var result JsapiTicket
	arrbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(arrbyte, &result)
	if err != nil {
		return err
	}
	SetWxKey("jsapi_ticket", result.Ticket)
	return nil
}

func GetWxKey(key string) (string, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	accessToken, _ := redis.String(conn.Do("GET", key))
	return accessToken, nil
}

func SetWxKey(key string, value string) error {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("SET", key, value)
	conn.Send("EXPIRE", key, 7200)
	conn.Do("EXEC")
	return nil
}
