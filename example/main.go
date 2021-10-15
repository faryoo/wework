package main

import (
	"flag"
	"fmt"
	wechat "github.com/faryoo/wework"
	"github.com/faryoo/wework/cache"
	workConfig "github.com/faryoo/wework/work/config"
	"log"
	"net/http"
)

func main() {
	flag.Parsed()

	http.HandleFunc("/", serveWechat)
	fmt.Println("wechat server listener at", ":8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}

}

func serveWechat(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        "192.168.6.100:6379",
		Password:    "",
		Database:    6,
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 60,
	}
	redisCache := cache.NewRedis(redisOpts)
	workcfg := &workConfig.Config{
		CorpID:         "ww16dcca8975dc595e",
		AgentID:        "1000014",
		CorpSecret:     "VQd7y88ECwA6gA_BENK9ab7eb2fOfHTd2DQH38sTRa8",
		Token:          "olGwFFZeGdtULjRYs0",
		EncodingAESKey: "QacEZIJd7iM6FszkT5St6Tc9CO3Z0SGc1tpELGSGIrD",
		Cache:          redisCache,
	}
	work := wc.GetWork(workcfg)

	userinfo, _ := work.GetUser().GetUserInfo("akita")
	log.Println(userinfo)
}
