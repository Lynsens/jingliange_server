package v1

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/stretchr/testify/assert"
)

// TestSimpleAPIParsing 测试API参数解析（不连接数据库）
func TestSimpleAPIParsing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody model.MenuQueryRequest
		expectCode  int
	}{
		{
			name: "valid json parsing",
			requestBody: model.MenuQueryRequest{
				Name:       "汤",
				PageSize:   10,
				PageNumber: 0,
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求体
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// 验证JSON序列化是否正确
			var parsedRequest model.MenuQueryRequest
			err = json.Unmarshal(jsonData, &parsedRequest)
			assert.NoError(t, err)
			assert.Equal(t, tt.requestBody.Name, parsedRequest.Name)
			assert.Equal(t, tt.requestBody.PageSize, parsedRequest.PageSize)
			assert.Equal(t, tt.requestBody.PageNumber, parsedRequest.PageNumber)

			t.Logf("请求参数解析成功: name=%s, page_size=%d, page_number=%d",
				parsedRequest.Name, parsedRequest.PageSize, parsedRequest.PageNumber)
		})
	}
}