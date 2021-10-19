package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/faryoo/wework/sync/context"
	"github.com/faryoo/wework/sync/message"
	log "github.com/sirupsen/logrus"

	"github.com/faryoo/wework/util"
)

// Server struct
type Server struct {
	*context.Context
	Writer  http.ResponseWriter
	Request *http.Request

	skipValidate bool

	openID string

	RequestRawXMLMsg  []byte
	RequestMsg        *message.MixMessage
	ResponseRawXMLMsg []byte
	ResponseMsg       interface{}

	isSafeMode bool
	random     []byte
	nonce      string
	timestamp  int64
}

// NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

// SkipValidate set skip validate
func (srv *Server) SkipValidate(skip bool) {
	srv.skipValidate = skip
}

// Validate 校验请求是否合法
func (srv *Server) Validate() bool {
	if srv.skipValidate {
		return true
	}
	timestamp := srv.Query("timestamp")
	nonce := srv.Query("nonce")
	signature := srv.Query("msg_signature")
	log.Debugf("validate signature, timestamp=%s, nonce=%s", timestamp, nonce)
	fmt.Println(srv.Token)
	return signature == util.Signature(srv.Token, timestamp, nonce)
}

// getMessage 解析微信返回的消息
func (srv *Server) GetMessage() (*message.MixMessage, error) {
	var rawXMLMsgBytes []byte
	var err error

	var encryptedXMLMsg message.EncryptedXMLMsg
	if err := xml.NewDecoder(srv.Request.Body).Decode(&encryptedXMLMsg); err != nil {
		return nil, fmt.Errorf("从body中解析xml失败,err=%v", err)
	}

	// 验证消息签名
	timestamp := srv.Query("timestamp")
	srv.timestamp, err = strconv.ParseInt(timestamp, 10, 32)
	if err != nil {
		return nil, err
	}
	nonce := srv.Query("nonce")
	srv.nonce = nonce
	msgSignature := srv.Query("msg_signature")
	msgSignatureGen := util.Signature(srv.Token, timestamp, nonce, encryptedXMLMsg.EncryptedMsg)
	if msgSignature != msgSignatureGen {
		return nil, fmt.Errorf("消息不合法，验证签名失败")
	}
	// 解密
	srv.random, rawXMLMsgBytes, err = util.DecryptMsg(srv.CorpID, encryptedXMLMsg.EncryptedMsg, srv.EncodingAESKey)
	if err != nil {
		return nil, fmt.Errorf("消息解密失败, err=%v", err)
	}

	srv.RequestRawXMLMsg = rawXMLMsgBytes

	return srv.parseRequestMessage(rawXMLMsgBytes)
}

func (srv *Server) parseRequestMessage(rawXMLMsgBytes []byte) (msg *message.MixMessage, err error) {
	msg = &message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, msg)
	return
}
