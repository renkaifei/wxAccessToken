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

func main() {
	for {
		accessToken := GetWxKey()
		if accessToken == "" {
			resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx423d0019a7141df4&secret=b202c14ae77ef834428b1a92385ef853")
			if err != nil {
				log.Fatal("error:" + err.Error())
				return
			}

			arrbyte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("error:" + err.Error())
				return
			}
			var result AccessToken
			err = json.Unmarshal(arrbyte, &result)
			if err != nil {
				log.Fatal("error:" + err.Error())
				return
			}
			SetWxKey(result.Token)
		}
		time.Sleep(1000 * 60 * 90 * time.Millisecond)
	}
}

func GetWxKey() string {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	accessToken, _ := redis.String(conn.Do("GET", "wx_access_token"))
	return accessToken
}

func SetWxKey(value string) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("SET", "wx_access_token", value)
	conn.Send("EXPIRE", "wx_access_token", 5400)
	conn.Do("EXEC")
}
