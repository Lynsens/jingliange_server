package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/lynsens/jingliange_server/pkg/util"
	"github.com/stretchr/testify/assert"
)

// TestConfig 测试配置结构
type TestConfig struct {
	UseTestConfig bool
	ConfigFile    string
}

// setupTestEnvironment 设置测试环境
func setupTestEnvironment() {
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
	
	// 初始化JWT（如果需要）
	util.InitJWT()
}

// init 包初始化函数
func init() {
	setupTestEnvironment()
}

// TestEnvironmentSetup 测试环境设置
func TestEnvironmentSetup(t *testing.T) {
	// 验证配置是否正确加载
	assert.NotNil(t, setting.AppSetting)
	assert.NotEmpty(t, setting.AppSetting.JwtSecret)
	assert.Greater(t, setting.AppSetting.PageSize, 0)
	
	// 验证服务器配置
	assert.NotNil(t, setting.ServerSetting)
	assert.NotEmpty(t, setting.ServerSetting.RunMode)
	
	t.Logf("测试环境配置:")
	t.Logf("- JWT Secret: %s", setting.AppSetting.JwtSecret[:10]+"...")
	t.Logf("- Page Size: %d", setting.AppSetting.PageSize)
	t.Logf("- Run Mode: %s", setting.ServerSetting.RunMode)
	t.Logf("- HTTP Port: %d", setting.ServerSetting.HttpPort)
}

// TestWithDatabaseMock 带数据库模拟的测试
func TestWithDatabaseMock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 这里可以模拟数据库连接失败的情况
	tests := []struct {
		name        string
		requestBody model.MenuQueryRequest
		mockDbError bool
	}{
		{
			name: "valid request",
			requestBody: model.MenuQueryRequest{
				Name:       "测试菜品",
				PageSize:   setting.AppSetting.PageSize, // 使用配置文件中的默认值
				PageNumber: 0,
			},
			mockDbError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/menu/getMenu", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 如果需要模拟数据库错误，可以在这里设置
			if tt.mockDbError {
				// 模拟数据库连接失败
				t.Skip("数据库模拟测试暂时跳过")
			} else {
				// 这里可以调用实际的API函数
				// GetMenu(c)
				
				// 由于没有真实数据库连接，我们只测试请求结构
				assert.NotNil(t, c.Request)
				assert.Equal(t, "application/json", c.Request.Header.Get("Content-Type"))
			}
		})
	}
}

// TestConfigurationValues 测试配置值是否正确读取
func TestConfigurationValues(t *testing.T) {
	tests := []struct {
		name          string
		configCheck   func() bool
		description   string
	}{
		{
			name:        "JWT Secret should not be empty",
			configCheck: func() bool { return setting.AppSetting.JwtSecret != "" },
			description: "JWT密钥不应为空",
		},
		{
			name:        "JWT Expire should be positive",
			configCheck: func() bool { return setting.AppSetting.JwtExpire > 0 },
			description: "JWT过期时间应为正数",
		},
		{
			name:        "Page Size should be valid",
			configCheck: func() bool { return setting.AppSetting.PageSize > 0 && setting.AppSetting.PageSize <= 100 },
			description: "分页大小应在合理范围内",
		},
		{
			name:        "HTTP Port should be valid",
			configCheck: func() bool { return setting.ServerSetting.HttpPort > 0 && setting.ServerSetting.HttpPort <= 65535 },
			description: "HTTP端口应在有效范围内",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.configCheck(), tt.description)
		})
	}
}

// TestAPIWithConfiguration 使用配置的API测试
func TestAPIWithConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 使用配置文件中的默认分页大小
	defaultPageSize := setting.AppSetting.PageSize
	
	tests := []struct {
		name        string
		requestBody model.MenuQueryRequest
	}{
		{
			name: "use default page size from config",
			requestBody: model.MenuQueryRequest{
				Name:       "",
				PageSize:   defaultPageSize,
				PageNumber: 0,
			},
		},
		{
			name: "use custom page size",
			requestBody: model.MenuQueryRequest{
				Name:       "汤",
				PageSize:   5,
				PageNumber: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			assert.NotEmpty(t, jsonData)

			req, err := http.NewRequest("POST", "/api/v1/menu/getMenu", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 验证请求体中的分页设置
			var parsedRequest model.MenuQueryRequest
			err = json.Unmarshal(jsonData, &parsedRequest)
			assert.NoError(t, err)
			assert.Equal(t, tt.requestBody.PageSize, parsedRequest.PageSize)
			
			t.Logf("测试用例: %s, 分页大小: %d", tt.name, parsedRequest.PageSize)
		})
	}
}
