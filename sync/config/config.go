// Package config 企业微信config配置
package config

// Config for 企业微信
type Config struct {
	CorpID         string `json:"corp_id"`
	Secret         string `json:"secret"`           // corp_secret,如果需要获取会话存档实例，当前参数请填写聊天内容存档的Secret，可以在企业微信管理端--管理工具--聊天内容存档查看
	Token          string `json:"token"`            // 微信客服回调配置，用于生成签名校验回调请求的合法性
	EncodingAESKey string `json:"encoding_aes_key"` // 微信客服回调p配置，用于解密回调消息内容对应的密文
}
