package appmessage

import (
	"fmt"
	"github.com/faryoo/wework/util"
	"github.com/faryoo/wework/work/context"
	"strings"
)

const (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// MsgType 发送消息类型
type MsgType string

const (
	// MsgTypeNews 图文消息
	MsgTypeNews MsgType = "mpnews"
	// MsgTypeText 文本
	MsgTypeText MsgType = "text"
)

// AppMsg 群发消息
type AppMsg struct {
	*context.Context
	preview bool
}

// NewAppMsg new
func NewAppMsg(ctx *context.Context) *AppMsg {
	appmsg := new(AppMsg)
	appmsg.Context = ctx
	appmsg.preview = false
	return appmsg
}

// User 发送的用户
type User struct {
	PartyID []string
	TagID   []string
	UserID  []string
}

// Result 群发返回结果
type Result struct {
	util.CommonError
	MsgID     int64  `json:"msg_id"`
	MsgDataID int64  `json:"msg_data_id"`
	MsgStatus string `json:"msg_status"`
}

// sendRequest 发送请求的数据
type sendRequest struct {
	// 根据OpenID发送
	ToUser  interface{} `json:"touser,omitempty"`
	ToTag   interface{} `json:"totag,omitempty"`
	ToParty interface{} `json:"toparty,omitempty"`
	// 应用ID
	AgentID string `json:"agentid,omitempty"`
	// 发送文本
	Text map[string]interface{} `json:"text,omitempty"`
	// 发送图文消息
	Mpnews            map[string]interface{} `json:"mpnews,omitempty"`
	MsgType           MsgType                `json:"msgtype"`
	SendIgnoreReprint int32                  `json:"send_ignore_reprint,omitempty"`
}

// SendText 群发文本
// user 为nil，表示全员发送
// &User{TagID:2} 根据tag发送
// &User{OpenID:[]string("xxx","xxx")} 根据openid发送
func (appmsg *AppMsg) SendText(user *User, content string) (*Result, error) {
	ak, err := appmsg.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req := &sendRequest{
		ToUser:  nil,
		AgentID: appmsg.AgentID,
		MsgType: MsgTypeText,
	}
	req.Text = map[string]interface{}{
		"content": content,
	}
	req = appmsg.chooseTagOrUserOrParty(user, req)
	url := fmt.Sprintf("%s?access_token=%s", sendURL, ak)
	data, err := util.PostJSON(url, req)
	if err != nil {
		return nil, err
	}
	res := &Result{}
	err = util.DecodeWithError(data, res, "SendText")
	return res, err
}

// SendNews 发送图文
func (appmsg *AppMsg) SendNews(user *User, mediaID string, ignoreReprint bool) (*Result, error) {
	ak, err := appmsg.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req := &sendRequest{
		ToUser:  nil,
		MsgType: MsgTypeNews,
	}
	if ignoreReprint {
		req.SendIgnoreReprint = 1
	}
	req.Mpnews = map[string]interface{}{
		"media_id": mediaID,
	}
	req = appmsg.chooseTagOrUserOrParty(user, req)
	url := fmt.Sprintf("%s?access_token=%s", sendURL, ak)
	data, err := util.PostJSON(url, req)
	if err != nil {
		return nil, err
	}
	res := &Result{}
	err = util.DecodeWithError(data, res, "SendNews")
	return res, err
}

func (appmsg *AppMsg) chooseTagOrUserOrParty(user *User, req *sendRequest) (ret *sendRequest) {

	if user == nil {
		req.ToUser = "@all"

	} else {
		if appmsg.preview {
			// 预览 默认发给第一个用户
			if len(user.UserID) != 0 {
				req.ToUser = user.UserID[0]
			}
		} else {
			if len(user.TagID) != 0 {
				req.ToTag = strings.Join(user.TagID, " | ")
			}
			if len(user.PartyID) != 0 {
				req.ToParty = strings.Join(user.PartyID, " | ")
			}
			if len(user.UserID) != 0 {
				req.ToUser = strings.Join(user.UserID, " | ")
			}
		}
	}
	return req
}
