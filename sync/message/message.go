package message

import (
	"encoding/xml"
)

// EventType 事件类型
type EventType string

const (
	// EventSubscribe 订阅
	EventSubscribe EventType = "subscribe"
	// EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	// EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	// EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	// EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	// EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	// EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	// EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	// EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	// EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	// EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	// EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
	// EventTemplateSendJobFinish 发送模板消息推送通知
	EventTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
	// EventMassSendJobFinish 群发消息推送通知
	EventMassSendJobFinish = "MASSSENDJOBFINISH"
	// EventWxaMediaCheck 异步校验图片/音频是否含有违法违规内容推送事件
	EventWxaMediaCheck = "wxa_media_check"
)

// MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	// 事件相关
	Event          EventType `xml:"Event"`
	ChangeType     string    `xml:"ChangeType"`
	UserID         string    `xml:"UserID"`
	Position       string    `xml:"Position"`
	Name           string    `xml:"Name"`
	Department     []int     `xml:"Department"`
	MainDepartment int       `xml:"MainDepartment"`
	IsLeaderInDept []int     `xml:"IsLeaderInDept"`
	Id             string    `xml:"Id"`
	ParentId       string    `xml:"ParentId"`
	Status         int       `xml:"Status"`
}

// EncryptedXMLMsg 安全模式下的消息体
type EncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"    json:"Encrypt"`
}

// ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptedMsg string   `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        string   `xml:"Nonce"        json:"Nonce"`
}

// CDATA  使用该类型,在序列化为 xml 文本时文本会被解析器忽略
type CDATA string

// MarshalXML 实现自己的序列化方法
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

// GetOpenID get the FromUserName value
func (msg *MixMessage) GetUserID() string {
	return msg.UserID
}
