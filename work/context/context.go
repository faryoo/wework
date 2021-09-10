package context

import (
	"wechat-work/credential"
	"wechat-work/work/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
