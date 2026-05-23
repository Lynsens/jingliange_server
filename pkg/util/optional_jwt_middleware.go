package util

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/logging"
)

// OptionalJWT 可选的JWT中间件 - 如果有token则验证，没有token也允许通过
func OptionalJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token != "" {
			// 如果提供了token，则验证
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
				logging.Info("OptionalJWT中间件 - 检测到token，开始验证")

				claims, err := ParseToken(token)
				if err != nil {
					logging.Warn("OptionalJWT中间件 - token验证失败，继续处理为未认证用户:", err)
					// token无效，但继续处理为未认证用户
				} else {
					// token有效，设置用户ID
					c.Set("user_id", claims.UserID)
					if claims.Role != "" {
						c.Set("role", claims.Role)
					}
					logging.Info("OptionalJWT中间件 - token验证成功, 用户ID:", claims.UserID)
				}
			}
		}

		c.Next()
	}
}
