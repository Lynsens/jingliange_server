package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/lynsens/jingliange_server/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

// @Summary 管理员登录
// @Description 使用配置文件中的管理员账号登录，返回管理员JWT token
// @Tags Admin
// @Accept json
// @Param auth body model.AdminLoginRequest true "管理员登录参数" schemaexample({"username":"admin","password":"<admin-password>"})
// @Produce json
// @Success 200 {object} app.Response{data=map[string]interface{}} "{"code":200,"msg":"ok","data":{"token":"...","username":"admin","role":"admin"}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"请求参数错误","data":"Username and password are required"}"
// @Failure 401 {object} app.Response "{"code":20004,"msg":"Token错误","data":"Invalid username or password"}"
// @Router /api/admin/login [post]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}

	var loginReq model.AdminLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Username and password are required")
		return
	}

	if setting.AdminSetting.Username == "" || setting.AdminSetting.PasswordHash == "" {
		logging.Error("Admin login - admin credentials are not configured")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if loginReq.Username != setting.AdminSetting.Username {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(setting.AdminSetting.PasswordHash), []byte(loginReq.Password)); err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "Invalid username or password")
		return
	}

	token, err := util.GenerateAdminToken(loginReq.Username)
	if err != nil {
		logging.Error("Admin login - token generation failed:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"token":    token,
		"username": loginReq.Username,
		"role":     util.AdminRole,
	})
}
