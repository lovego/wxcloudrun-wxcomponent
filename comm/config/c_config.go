package config

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/go-ini/ini"
)

type Oplatform struct {
	Appid  string
	Token  string
	AesKey string
}

type Whitelist struct {
	Ip []string
}

var OplatformConf = &Oplatform{}
var WhitelistConf = &Whitelist{}

func init() {
	var err error
	cfg, err = ini.Load("comm/config/c_server.conf")
	if err != nil {
		log.Errorf("load server.conf': %v", err)
		return
	}
	mapTo("oplatform", OplatformConf)
	mapTo("whitelist", WhitelistConf)
	log.Debug("load server oplatform conf")
}
