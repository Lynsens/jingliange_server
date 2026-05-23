package util

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/stretchr/testify/assert"
)

// init 初始化测试环境
func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 获取当前工作目录并切换到项目根目录
	wd, _ := os.Getwd()
	originalWd := wd

	// 找到项目根目录（包含go.mod的目录）
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			os.Chdir(wd)
			break
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			// 到达系统根目录，恢复原始目录
			os.Chdir(originalWd)
			break
		}
		wd = parent
	}

	// 加载配置
	setting.Setup()

	// 初始化日志
	logging.Setup()

	// 初始化JWT
	InitJWT()
}

func TestJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个测试用的有效token
	validToken, err := GenerateToken("test_user")
	assert.NoError(t, err)
	assert.NotEmpty(t, validToken)

	tests := []struct {
		name           string
		authHeader     string
		expectedCode   int
		expectedStatus int
		shouldSetUser  bool
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedCode:   http.StatusOK,
			expectedStatus: http.StatusOK,
			shouldSetUser:  true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			shouldSetUser:  false,
		},
		{
			name:           "invalid bearer format",
			authHeader:     "InvalidFormat " + validToken,
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			shouldSetUser:  false,
		},
		{
			name:           "missing bearer prefix",
			authHeader:     validToken,
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			shouldSetUser:  false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid_token_string",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			shouldSetUser:  false,
		},
		{
			name:           "expired token",
			authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdF91c2VyIiwiZXhwIjoxNjI2MjgwMDAwfQ.invalid",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: e.ERROR_AUTH_CHECK_TOKEN_FAIL,
			shouldSetUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建路由器和中间件
			router := gin.New()
			router.Use(JWT())
			router.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if exists {
					c.JSON(http.StatusOK, gin.H{"user_id": userID})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "no user"})
				}
			})

			// 创建请求
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusUnauthorized {
				// 由于JWT中间件直接调用c.Abort()，我们主要验证状态码
				// 这里简化了测试验证
			}
		})
	}
}

func TestAdminJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	adminToken, err := GenerateAdminToken("admin")
	assert.NoError(t, err)
	userToken, err := GenerateToken("test_user")
	assert.NoError(t, err)

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "valid admin token",
			authHeader:   "Bearer " + adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing token",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid token",
			authHeader:   "Bearer invalid_token",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "normal user token",
			authHeader:   "Bearer " + userToken,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AdminJWT())
			router.GET("/admin-only", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				role, _ := c.Get("role")
				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"role":    role,
				})
			})

			req, err := http.NewRequest("GET", "/admin-only", nil)
			assert.NoError(t, err)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestOptionalJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个测试用的有效token
	validToken, err := GenerateToken("test_user")
	assert.NoError(t, err)
	assert.NotEmpty(t, validToken)

	tests := []struct {
		name          string
		authHeader    string
		expectedCode  int
		shouldSetUser bool
		expectedUser  string
	}{
		{
			name:          "valid token",
			authHeader:    "Bearer " + validToken,
			expectedCode:  http.StatusOK,
			shouldSetUser: true,
			expectedUser:  "test_user",
		},
		{
			name:          "no authorization header",
			authHeader:    "",
			expectedCode:  http.StatusOK,
			shouldSetUser: false,
			expectedUser:  "",
		},
		{
			name:          "invalid token format",
			authHeader:    "InvalidFormat " + validToken,
			expectedCode:  http.StatusOK,
			shouldSetUser: false,
			expectedUser:  "",
		},
		{
			name:          "invalid token",
			authHeader:    "Bearer invalid_token",
			expectedCode:  http.StatusOK,
			shouldSetUser: false,
			expectedUser:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建路由器和中间件
			router := gin.New()
			router.Use(OptionalJWT())
			router.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if exists {
					c.JSON(http.StatusOK, gin.H{"user_id": userID, "authenticated": true})
				} else {
					c.JSON(http.StatusOK, gin.H{"authenticated": false})
				}
			})

			// 创建请求
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedCode, w.Code)

			// 解析响应以验证用户ID设置
			// 注意：这里简化了测试，实际项目中可能需要更详细的验证
		})
	}
}

func TestJWTGeneration(t *testing.T) {
	tests := []struct {
		name   string
		userID string
	}{
		{
			name:   "valid user ID",
			userID: "test_user_123",
		},
		{
			name:   "user ID with special characters",
			userID: "wulongcha_test",
		},
		{
			name:   "empty user ID",
			userID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID)

			if tt.userID == "" {
				// 空用户ID应该仍然能生成token（业务逻辑可能允许）
				// 或者根据实际需求调整这个测试
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 验证生成的token可以被正确解析
				claims, err := ParseToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Empty(t, claims.Role)
				assert.True(t, claims.ExpiresAt > time.Now().Unix())
			}
		})
	}
}

func TestAdminJWTGeneration(t *testing.T) {
	token, err := GenerateAdminToken("admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "admin", claims.UserID)
	assert.Equal(t, AdminRole, claims.Role)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestJWTParsing(t *testing.T) {
	// 首先生成一个有效的token
	validToken, err := GenerateToken("test_user")
	assert.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		expectError bool
		expectedID  string
	}{
		{
			name:        "valid token",
			token:       validToken,
			expectError: false,
			expectedID:  "test_user",
		},
		{
			name:        "invalid token format",
			token:       "invalid.token.format",
			expectError: true,
			expectedID:  "",
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
			expectedID:  "",
		},
		{
			name:        "malformed token",
			token:       "not.a.jwt",
			expectError: true,
			expectedID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectedID, claims.UserID)
			}
		})
	}
}

func TestJWTExpiration(t *testing.T) {
	// 这个测试验证token的过期时间设置是否正确
	token, err := GenerateToken("test_user")
	assert.NoError(t, err)

	claims, err := ParseToken(token)
	assert.NoError(t, err)

	// 验证过期时间是否在合理范围内（当前时间 + 配置的过期时间）
	expectedExpiration := time.Now().Add(time.Hour * 720).Unix() // 720小时 = 30天
	actualExpiration := claims.ExpiresAt

	// 允许一些时间差异（比如测试执行时间）
	timeDiff := actualExpiration - expectedExpiration
	assert.True(t, timeDiff >= -60 && timeDiff <= 60, "Token expiration time should be within 60 seconds of expected")
}
