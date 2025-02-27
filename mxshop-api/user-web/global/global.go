package global

import (
	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/config"
	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig=&config.ServerConfig{}
)
