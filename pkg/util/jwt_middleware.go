package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
)

// JWT 中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("Authorization")

		if token == "" {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			// 去掉 Bearer 前缀
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
			}

			claims, err := ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				// 将用户ID存入上下文
				c.Set("user_id", claims.UserID)
			}
		}

		if code != e.SUCCESS {
			appG := app.Gin{C: c}
			appG.Response(http.StatusUnauthorized, code, data)
			c.Abort()
			return
		}

		c.Next()
	}
}
