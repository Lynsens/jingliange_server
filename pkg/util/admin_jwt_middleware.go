package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
)

// AdminJWT 管理员JWT中间件
func AdminJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			abortAdminAuth(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			abortAdminAuth(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL)
			return
		}

		claims, err := ParseToken(token[7:])
		if err != nil {
			logging.Error("AdminJWT中间件 - token解析失败:", err)
			abortAdminAuth(c, e.ERROR_AUTH_TOKEN)
			return
		}

		if claims.Role != AdminRole {
			logging.Error("AdminJWT中间件 - 非管理员token, 用户ID:", claims.UserID)
			abortAdminAuth(c, e.ERROR_AUTH)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func abortAdminAuth(c *gin.Context, code int) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusUnauthorized, code, nil)
	c.Abort()
}
