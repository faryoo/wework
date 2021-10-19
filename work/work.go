package work

import (
	"github.com/faryoo/wework/credential"
	"github.com/faryoo/wework/work/appmessage"
	"github.com/faryoo/wework/work/config"
	"github.com/faryoo/wework/work/context"
	"github.com/faryoo/wework/work/member/department"
	"github.com/faryoo/wework/work/member/user"
	"github.com/faryoo/wework/work/menu"
	"github.com/faryoo/wework/work/msgaudit"
	"github.com/faryoo/wework/work/oauth"
	"github.com/faryoo/wework/work/server"
	"net/http"
)

// Work 企业微信
type Work struct {
	ctx *context.Context
}

// NewWork init work
func NewWork(cfg *config.Config) *Work {
	defaultAkHandle := credential.NewWorkAccessToken(cfg.CorpID, cfg.CorpSecret, cfg.AgentID, credential.CacheKeyWorkPrefix, cfg.Cache)
	ctx := &context.Context{
		Config:            cfg,
		AccessTokenHandle: defaultAkHandle,
	}
	return &Work{ctx: ctx}
}

// SetAccessTokenHandle 自定义access_token获取方式
func (wk *Work) SetAccessTokenHandle(accessTokenHandle credential.AccessTokenHandle) {
	wk.ctx.AccessTokenHandle = accessTokenHandle
}

// GetContext get Context
func (wk *Work) GetContext() *context.Context {
	return wk.ctx
}

// GetMenu 菜单管理接口
func (wk *Work) GetMenu() *menu.Menu {
	return menu.NewMenu(wk.ctx)
}

// GetServer 消息管理：接收事件，被动回复消息管理
func (wk *Work) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	srv := server.NewServer(wk.ctx)
	srv.Request = req
	srv.Writer = writer
	return srv
}

// GetAccessToken 获取access_token
func (wk *Work) GetAccessToken() (string, error) {
	return wk.ctx.GetAccessToken()
}

// GetOauth get oauth
func (wk *Work) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wk.ctx)
}

// GetMsgAudit get msgAudit
func (wk *Work) GetMsgAudit() (*msgaudit.Client, error) {
	return msgaudit.NewClient(wk.ctx.Config)
}

func (wk *Work) GetUser() *user.User {
	return user.NewUser(wk.ctx)
}
func (wk *Work) GetDepartMent() *department.DepartMent {
	return department.NewDepartMent(wk.ctx)
}

func (wk *Work) GetAppMsg() *appmessage.AppMsg {
	return appmessage.NewAppMsg(wk.ctx)
}
