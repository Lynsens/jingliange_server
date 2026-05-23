package model

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"jingliange-admin"`
}
