package util

import (
	"testing"
	"time"

	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化测试环境的配置
	setting.AppSetting = &setting.App{
		JwtSecret: "test_secret_key_for_jwt_testing_123456",
		JwtExpire: 720, // 720小时
	}
	InitJWT()
}

func TestJWTTokenGeneration(t *testing.T) {
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
			name:   "numeric user ID",
			userID: "123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// 验证生成的token可以被正确解析
			claims, err := ParseToken(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.True(t, claims.ExpiresAt > time.Now().Unix())
		})
	}
}

func TestJWTTokenParsing(t *testing.T) {
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

func TestJWTTokenExpiration(t *testing.T) {
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

func TestJWTTokenRoundTrip(t *testing.T) {
	// 测试完整的token生成和解析流程
	originalUserID := "test_user_round_trip"

	// 生成token
	token, err := GenerateToken(originalUserID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 解析token
	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// 验证数据一致性
	assert.Equal(t, originalUserID, claims.UserID)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestJWTMultipleTokens(t *testing.T) {
	// 测试为不同用户生成多个token
	users := []string{"user1", "user2", "user3", "admin", "test_123"}
	tokens := make(map[string]string)

	// 为每个用户生成token
	for _, userID := range users {
		token, err := GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		tokens[userID] = token
	}

	// 验证每个token都能正确解析回原用户ID
	for userID, token := range tokens {
		claims, err := ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	}
}
