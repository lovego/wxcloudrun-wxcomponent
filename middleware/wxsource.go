package middleware

import (
	"github.com/gin-gonic/gin"
)

// WXSourceMiddleWare 中间件 判断是否来源于微信
func WXSourceMiddleWare(c *gin.Context) {
	// if _, ok := c.Request.Header[http.CanonicalHeaderKey("x-wx-source")]; ok {
	// 	fmt.Println("[WXSourceMiddleWare]from wx")
	// 	c.Next()
	// } else {
	// 	c.Abort()
	// 	c.JSON(http.StatusUnauthorized, errno.ErrNotAuthorized)
	// }
	// 二开项目，未部署在腾讯云平台，跳过来源判断
	c.Next()
}
