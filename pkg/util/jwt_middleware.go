package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
)

// JWT 中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("Authorization")

		if token == "" {
			logging.Error("JWT中间件 - Authorization头为空")
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			// 去掉 Bearer 前缀
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
				logging.Info("JWT中间件 - 提取到token:", token[:20]+"...")
			} else {
				logging.Error("JWT中间件 - Authorization头格式错误，缺少Bearer前缀")
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}

			if code == e.SUCCESS {
				claims, err := ParseToken(token)
				if err != nil {
					logging.Error("JWT中间件 - token解析失败:", err)
					code = e.ERROR_AUTH_TOKEN
				} else {
					// 将用户ID存入上下文
					c.Set("user_id", claims.UserID)
					if claims.Role != "" {
						c.Set("role", claims.Role)
					}
					logging.Info("JWT中间件 - token验证成功, 用户ID:", claims.UserID)
				}
			}
		}

		if code != e.SUCCESS {
			logging.Error("JWT中间件 - 认证失败, 错误码:", code)
			appG := app.Gin{C: c}
			appG.Response(http.StatusUnauthorized, code, data)
			c.Abort()
			return
		}

		c.Next()
	}
}
