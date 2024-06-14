package custom

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type wxCallback struct {
	AppId      string `xml:"AppId"`
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

type WxCallbackComponentMsg struct {
	AppId                        string `xml:"AppId"`
	CreateTime                   int64  `xml:"CreateTime"`
	InfoType                     string `xml:"InfoType"`
	ComponentVerifyTicket        string `xml:"ComponentVerifyTicket"`
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int64  `xml:"AuthorizationCodeExpiredTime"`
	PreAuthCode                  string `xml:"PreAuthCode"`
}

type WxCallbackBizMsg struct {
	ToUserName   string  `xml:"ToUserName"`
	FromUserName string  `xml:"FromUserName"`
	CreateTime   int64   `xml:"CreateTime"`
	MsgType      string  `xml:"MsgType"`
	Event        string  `xml:"Event"`
	EventKey     string  `xml:"EventKey"`
	Ticket       string  `xml:"Ticket"`
	Latitude     float64 `xml:"Latitude"`
	Longitude    float64 `xml:"Longitude"`
	Precision    float64 `xml:"Precision"`
	// 消息
	Content string `xml:"Content"`
	MsgId   int64  `xml:"MsgId"`
}

func checkSign(c *gin.Context, encryptData string) error {
	var (
		timestamp = c.Query("timestamp")
		nonce     = c.Query("nonce")
	)
	signature := Signature(config.OplatformConf.Token, timestamp, nonce, encryptData)
	if c.Query("msg_signature") != signature {
		return errors.New("消息不合法，验证签名失败")
	}
	return nil
}

func Parse(c *gin.Context, msg interface{}) ([]byte, error) {
	tmp, _ := io.ReadAll(c.Request.Body)
	var data wxCallback
	if err := binding.XML.BindBody(tmp, &data); err != nil {
		return tmp, err
	}
	if err := checkSign(c, data.Encrypt); err != nil {
		return tmp, err
	}
	_, body, err := DecryptMsg(config.OplatformConf.Appid, data.Encrypt, config.OplatformConf.AesKey)
	if err != nil {
		return tmp, err
	}
	if err = xml.Unmarshal(body, &msg); err != nil {
		return tmp, err
	}
	if marshal, _ := json.Marshal(msg); len(marshal) != 0 {
		c.Request.Body = io.NopCloser(bytes.NewReader(marshal))
	}
	return tmp, nil
}
