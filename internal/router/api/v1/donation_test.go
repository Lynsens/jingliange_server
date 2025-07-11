package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/stretchr/testify/assert"
)

func TestGetDonationList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    model.DonationQueryRequest
		expectedCode   int
		expectedStatus int
	}{
		{
			name: "valid request with default values",
			requestBody: model.DonationQueryRequest{
				Year:       2025,
				Period:     "all",
				DonorName:  "",
				SortBy:     "time",
				SortOrder:  "desc",
				PageSize:   10,
				PageNumber: 0,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
		{
			name: "valid request with donor name filter",
			requestBody: model.DonationQueryRequest{
				Year:       2025,
				Period:     "first",
				DonorName:  "善心",
				SortBy:     "amount",
				SortOrder:  "desc",
				PageSize:   5,
				PageNumber: 0,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
		{
			name: "valid request with second period",
			requestBody: model.DonationQueryRequest{
				Year:       2024,
				Period:     "second",
				DonorName:  "",
				SortBy:     "time",
				SortOrder:  "asc",
				PageSize:   20,
				PageNumber: 1,
			},
			expectedCode:   http.StatusOK,
			expectedStatus: e.SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/donation/getDonationList", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			GetDonationList(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Code)
		})
	}
}

func TestGetDonationList_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody model.DonationQueryRequest
		description string
	}{
		{
			name: "invalid period",
			requestBody: model.DonationQueryRequest{
				Period:     "invalid",
				SortBy:     "time",
				SortOrder:  "desc",
				PageSize:   10,
				PageNumber: 0,
			},
			description: "invalid period value",
		},
		{
			name: "invalid sort_by",
			requestBody: model.DonationQueryRequest{
				Period:     "all",
				SortBy:     "invalid",
				SortOrder:  "desc",
				PageSize:   10,
				PageNumber: 0,
			},
			description: "invalid sort_by value",
		},
		{
			name: "invalid sort_order",
			requestBody: model.DonationQueryRequest{
				Period:     "all",
				SortBy:     "time",
				SortOrder:  "invalid",
				PageSize:   10,
				PageNumber: 0,
			},
			description: "invalid sort_order value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/donation/getDonationList", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			GetDonationList(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, e.INVALID_PARAMS, response.Code)
		})
	}
}

func TestCreateDonation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.DonationCreateRequest
		userID       string
		expectedCode int
	}{
		{
			name: "valid donation request",
			requestBody: model.DonationCreateRequest{
				DonorName: "善心人士",
				Amount:    100.50,
				Message:   "祝愿净莲阁越来越好",
			},
			userID:       "test_user",
			expectedCode: http.StatusOK,
		},
		{
			name: "valid donation with minimal info",
			requestBody: model.DonationCreateRequest{
				DonorName: "匿名",
				Amount:    50.00,
				Message:   "",
			},
			userID:       "test_user_2",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/donation/createDonation", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)

			CreateDonation(c)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestCreateDonation_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody model.DonationCreateRequest
		userID      string
		description string
	}{
		{
			name: "empty donor name",
			requestBody: model.DonationCreateRequest{
				DonorName: "",
				Amount:    100.00,
				Message:   "test message",
			},
			userID:      "test_user",
			description: "donor name is required",
		},
		{
			name: "invalid amount - zero",
			requestBody: model.DonationCreateRequest{
				DonorName: "测试用户",
				Amount:    0,
				Message:   "test message",
			},
			userID:      "test_user",
			description: "amount must be positive",
		},
		{
			name: "invalid amount - negative",
			requestBody: model.DonationCreateRequest{
				DonorName: "测试用户",
				Amount:    -10.00,
				Message:   "test message",
			},
			userID:      "test_user",
			description: "amount must be positive",
		},
		{
			name: "donor name too long",
			requestBody: model.DonationCreateRequest{
				DonorName: "这是一个非常非常非常非常非常非常长的捐赠者姓名，超过了系统限制的字符数量限制，应该被拒绝",
				Amount:    100.00,
				Message:   "test message",
			},
			userID:      "test_user",
			description: "donor name too long",
		},
		{
			name: "message too long",
			requestBody: model.DonationCreateRequest{
				DonorName: "测试用户",
				Amount:    100.00,
				Message:   "这是一个非常非常非常非常非常非常长的留言内容，超过了系统限制的字符数量限制，应该被拒绝。这个留言内容太长了，超出了预期的长度限制，系统应该返回错误信息。继续添加更多内容来确保超过限制...",
			},
			userID:      "test_user",
			description: "message too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/donation/createDonation", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)

			CreateDonation(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, e.INVALID_PARAMS, response.Code)
		})
	}
}

func TestCreateDonation_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	requestBody := model.DonationCreateRequest{
		DonorName: "测试用户",
		Amount:    100.00,
		Message:   "测试留言",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/donation/createDonation", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// 不设置user_id，模拟未认证状态

	CreateDonation(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.ERROR_AUTH, response.Code)
}

func TestGetDonationStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.DonationStatsRequest
		expectedCode int
	}{
		{
			name: "valid stats request for current year",
			requestBody: model.DonationStatsRequest{
				Year: time.Now().Year(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid stats request for specific year",
			requestBody: model.DonationStatsRequest{
				Year: 2024,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/donation/getDonationStats", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			GetDonationStats(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, e.SUCCESS, response.Code)
		})
	}
}

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  model.AuthRequest
		expectedCode int
	}{
		{
			name: "valid user authentication",
			requestBody: model.AuthRequest{
				UserID: "test_user_123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid user authentication with special characters",
			requestBody: model.AuthRequest{
				UserID: "wulongcha_test",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			AuthUser(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, e.SUCCESS, response.Code)

			// 验证响应包含token
			if response.Code == e.SUCCESS {
				dataMap, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				token, exists := dataMap["token"]
				assert.True(t, exists)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthUser_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody model.AuthRequest
		description string
	}{
		{
			name: "empty user ID",
			requestBody: model.AuthRequest{
				UserID: "",
			},
			description: "user ID is required",
		},
		{
			name: "user ID too long",
			requestBody: model.AuthRequest{
				UserID: "this_is_a_very_very_very_very_very_very_very_very_very_very_long_user_id_that_exceeds_the_maximum_allowed_length_for_user_identification_in_our_system",
			},
			description: "user ID too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			AuthUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response app.Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, e.INVALID_PARAMS, response.Code)
		})
	}
}

func TestAuthUser_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	AuthUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, e.INVALID_PARAMS, response.Code)
}
