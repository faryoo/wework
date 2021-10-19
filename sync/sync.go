package sync

import (
	"github.com/faryoo/wework/sync/config"
	"github.com/faryoo/wework/sync/context"
	"github.com/faryoo/wework/sync/server"
	"net/http"
)

// Work 企业微信
type Sync struct {
	ctx *context.Context
}

// NewWork init work
func NewSync(cfg *config.Config) *Sync {
	ctx := &context.Context{
		Config: cfg,
	}
	return &Sync{ctx: ctx}
}

// GetServer 消息管理：接收事件，被动回复消息管理
func (sy *Sync) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	srv := server.NewServer(sy.ctx)
	srv.Request = req
	srv.Writer = writer
	return srv
}
