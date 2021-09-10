package context

import (
	"github.com/faryoo/wework/credential"
	"github.com/faryoo/wework/work/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
