package js

import (
	"fmt"
	"github.com/faryoo/wework/credential"
	"github.com/faryoo/wework/util"
	"github.com/faryoo/wework/work/context"
)

// Js wx jssdk
type Js struct {
	*context.Context
	credential.JsTicketHandle
}
type Config struct {
	Corpid    string `json:"corpid"`
	Agentid   string `json:"agentid"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	Signature string `json:"signature"`
}

// NewJs init
func NewJs(context *context.Context, appID string) *Js {
	js := new(Js)
	js.Context = context
	jsTicketHandle := credential.NewDefaultJsTicket(appID, credential.CacheKeyWorkPrefix, context.Cache)
	js.SetJsTicketHandle(jsTicketHandle)
	return js
}

// SetJsTicketHandle 自定义js ticket取值方式
func (js *Js) SetJsTicketHandle(ticketHandle credential.JsTicketHandle) {
	js.JsTicketHandle = ticketHandle
}

// GetConfig 第三方平台 - 获取jssdk需要的配置参数
// uri 为当前网页地址
func (js *Js) GetConfig(uri, agentid string) (config *Config, err error) {
	config = new(Config)
	var accessToken string
	accessToken, err = js.GetAccessToken()
	if err != nil {
		return
	}
	var ticketStr string
	ticketStr, err = js.GetTicket(accessToken)
	if err != nil {
		return
	}

	nonceStr := util.RandomStr(16)
	timestamp := util.GetCurrTS()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	sigStr := util.Signature(str)

	config.Corpid = js.CorpID
	config.Agentid = js.AgentID
	config.NonceStr = nonceStr
	config.Timestamp = timestamp
	config.Signature = sigStr
	return
}
