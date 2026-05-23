package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lynsens/jingliange_server/pkg/setting"
)

var jwtSecret []byte

const AdminRole = "admin"

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.StandardClaims
}

// InitJWT 初始化JWT密钥
func InitJWT() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// GenerateToken 生成JWT token
func GenerateToken(userID string) (string, error) {
	return generateToken(userID, "")
}

// GenerateAdminToken 生成管理员JWT token
func GenerateAdminToken(username string) (string, error) {
	return generateToken(username, AdminRole)
}

func generateToken(userID, role string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(setting.AppSetting.JwtExpire) * time.Hour)

	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    "jingliange_server",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// ValidateToken 验证token是否有效
func ValidateToken(token string) bool {
	_, err := ParseToken(token)
	return err == nil
}
