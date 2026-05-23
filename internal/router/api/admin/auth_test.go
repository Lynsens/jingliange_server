package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/lynsens/jingliange_server/pkg/util"
	"github.com/stretchr/testify/assert"
)

const testAdminPasswordHash = "$2a$10$DNblWL.KBv/DALuWq7Ac5OmNv2yWEAQm0s9kdlZz1SnHixipyw3CW"

func setupAdminLoginTest(t *testing.T) {
	t.Helper()

	originalApp := *setting.AppSetting
	originalAdmin := *setting.AdminSetting

	setting.AppSetting.JwtSecret = "test_secret_key_for_admin_login"
	setting.AppSetting.JwtExpire = 720
	setting.AdminSetting.Username = "admin"
	setting.AdminSetting.PasswordHash = testAdminPasswordHash
	util.InitJWT()

	t.Cleanup(func() {
		*setting.AppSetting = originalApp
		*setting.AdminSetting = originalAdmin
		util.InitJWT()
	})
}

func performAdminLogin(t *testing.T, body interface{}) (int, app.Response) {
	t.Helper()

	jsonData, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	Login(c)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	return w.Code, response
}

func TestAdminLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAdminLoginTest(t)

	statusCode, response := performAdminLogin(t, model.AdminLoginRequest{
		Username: "admin",
		Password: "jingliange-admin",
	})

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, e.SUCCESS, response.Code)

	dataMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "admin", dataMap["username"])
	assert.Equal(t, util.AdminRole, dataMap["role"])

	token, ok := dataMap["token"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, token)

	claims, err := util.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "admin", claims.UserID)
	assert.Equal(t, util.AdminRole, claims.Role)
}

func TestAdminLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAdminLoginTest(t)

	tests := []struct {
		name string
		body model.AdminLoginRequest
	}{
		{
			name: "wrong username",
			body: model.AdminLoginRequest{
				Username: "other",
				Password: "jingliange-admin",
			},
		},
		{
			name: "wrong password",
			body: model.AdminLoginRequest{
				Username: "admin",
				Password: "wrong-password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, response := performAdminLogin(t, tt.body)

			assert.Equal(t, http.StatusUnauthorized, statusCode)
			assert.Equal(t, e.ERROR_AUTH, response.Code)
		})
	}
}

func TestAdminLoginMissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAdminLoginTest(t)

	tests := []struct {
		name string
		body model.AdminLoginRequest
	}{
		{
			name: "missing username",
			body: model.AdminLoginRequest{
				Password: "jingliange-admin",
			},
		},
		{
			name: "missing password",
			body: model.AdminLoginRequest{
				Username: "admin",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, response := performAdminLogin(t, tt.body)

			assert.Equal(t, http.StatusBadRequest, statusCode)
			assert.Equal(t, e.INVALID_PARAMS, response.Code)
		})
	}
}

func TestAdminLoginInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAdminLoginTest(t)

	req, err := http.NewRequest("POST", "/api/admin/login", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.INVALID_PARAMS, response.Code)
}
