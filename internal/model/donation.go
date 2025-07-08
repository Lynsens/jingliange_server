package model

import "time"

// Donation 捐款记录
type Donation struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	UserID     string    `json:"user_id" gorm:"type:varchar(64);not null;default:''" example:"user123"`    // 用户ID（必填）
	DonorName  string    `json:"donor_name" gorm:"type:varchar(32);not null;default:''" example:"善心人士"`    // 捐款人昵称（必填）
	Amount     float64   `json:"amount" gorm:"type:decimal(10,2);not null;default:0" example:"100.00"`     // 捐款金额（必填）
	DonateTime time.Time `json:"donate_time" gorm:"type:datetime;not null" example:"2025-07-08T10:00:00Z"` // 捐款时间（必填）
	IsVisible  int       `json:"is_visible" gorm:"type:tinyint(1);not null;default:1" example:"1"`         // 是否显示在榜单：0隐藏，1显示
	Message    string    `json:"message" gorm:"type:varchar(256);not null;default:''" example:"祝愿净莲阁越来越好"` // 留言
	Remarks    string    `json:"remarks" gorm:"type:varchar(256);not null;default:''" example:""`          // 备注
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" example:"2025-07-08T10:00:00Z"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" example:"2025-07-08T10:00:00Z"`
}

// DonationQueryRequest 功德榜查询请求
type DonationQueryRequest struct {
	Year       int    `json:"year" example:"2025"`       // 年份，0表示不限制
	Period     string `json:"period" example:"all"`      // 时间范围：all全年，first上半年，second下半年
	DonorName  string `json:"donor_name" example:"善心"`   // 捐款人昵称搜索（模糊匹配）
	SortBy     string `json:"sort_by" example:"time"`    // 排序方式：time时间，amount金额
	SortOrder  string `json:"sort_order" example:"desc"` // 排序顺序：desc降序，asc升序
	PageSize   int    `json:"page_size" example:"10"`    // 每页数量
	PageNumber int    `json:"page_number" example:"0"`   // 页码，从零开始
}

// DonationCreateRequest 创建捐款记录请求（用户ID从JWT token中获取）
type DonationCreateRequest struct {
	DonorName string  `json:"donor_name" example:"善心人士"`   // 捐款人昵称
	Amount    float64 `json:"amount" example:"100.00"`     // 捐款金额
	Message   string  `json:"message" example:"祝愿净莲阁越来越好"` // 留言
}

// DonationStatsRequest 捐款统计请求
type DonationStatsRequest struct {
	Year   int    `json:"year" example:"2025"`  // 年份，0表示不限制
	Period string `json:"period" example:"all"` // 时间范围：all全年，first上半年，second下半年
}

// DonationStats 捐款统计响应
type DonationStats struct {
	TotalAmount float64 `json:"total_amount" example:"5000.00"` // 总捐款金额
	TotalCount  int64   `json:"total_count" example:"50"`       // 总捐款人次
}

// User 用户信息（简化版，基于微信小程序）
type User struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(64)" example:"user123"`                                                                        // 用户ID（来自微信）
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" example:"2025-07-08T10:00:00Z"`                             // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" example:"2025-07-08T10:00:00Z"` // 更新时间
}

// AuthRequest 用户认证请求
type AuthRequest struct {
	UserID string `json:"user_id" example:"user123"` // 用户ID（来自微信小程序）
}
