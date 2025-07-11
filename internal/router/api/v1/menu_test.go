package v1

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
	"github.com/stretchr/testify/assert"
)

func TestGetMenu(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    model.MenuQueryRequest
		expectedCode   int
		expectedStatus int
	}{
		{
			name: "valid request with empty name",
			requestBody: model.MenuQueryRequest{
				Name:       "",
				PageSize:   10,
				PageNumber: 0,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
		{
			name: "valid request with name filter",
			requestBody: model.MenuQueryRequest{
				Name:       "汤",
				PageSize:   5,
				PageNumber: 0,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
		{
			name: "valid request with pagination",
			requestBody: model.MenuQueryRequest{
				Name:       "",
				PageSize:   3,
				PageNumber: 1,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求体
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// 创建HTTP请求
			req, err := http.NewRequest("POST", "/api/v1/menu/getMenu", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 创建Gin上下文
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 调用处理函数
			GetMenu(c)

			// 验证响应状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 解析响应体
			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 验证响应状态
			assert.Equal(t, tt.expectedStatus, response.Code)
		})
	}
}

func TestGetMenu_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试无效的JSON数据
	req, err := http.NewRequest("POST", "/api/v1/menu/getMenu", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GetMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.INVALID_PARAMS, response.Code)
}

func TestLikeMenu(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.MenuLikeRequest
		userID       string
		expectedCode int
	}{
		{
			name: "valid like request",
			requestBody: model.MenuLikeRequest{
				MenuID: 1,
			},
			userID:       "test_user",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/menu/like", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 模拟JWT中间件设置的用户ID
			c.Set("user_id", tt.userID)

			LikeMenu(c)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestLikeMenu_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	requestBody := model.MenuLikeRequest{
		MenuID: 1,
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/menu/like", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 不设置user_id，模拟未认证状态

	LikeMenu(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.ERROR_AUTH, response.Code)
}

func TestLikeMenu_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试无效的JSON数据
	req, err := http.NewRequest("POST", "/api/v1/menu/like", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "test_user")

	LikeMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.INVALID_PARAMS, response.Code)
}

func TestCommentMenu(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.MenuCommentRequest
		userID       string
		expectedCode int
	}{
		{
			name: "valid comment request",
			requestBody: model.MenuCommentRequest{
				MenuID:  1,
				Comment: "这道菜很好吃！",
			},
			userID:       "test_user",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/menu/comment", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)

			CommentMenu(c)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetMenuByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.MenuByIDRequest
		expectedCode int
	}{
		{
			name: "valid menu ID request",
			requestBody: model.MenuByIDRequest{
				MenuID: 1,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/menu/getMenuByID", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			GetMenuByID(c)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetMenuByID_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试无效的menu_id
	requestBody := model.MenuByIDRequest{
		MenuID: 0, // 无效的menu_id
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/menu/getMenuByID", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GetMenuByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.INVALID_PARAMS, response.Code)
}
