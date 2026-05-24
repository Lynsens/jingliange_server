package model

import "time"

type Menu struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Name          string     `json:"name" gorm:"type:varchar(64);not null;default:''" example:"紫菜汤"` // 菜品名称
	Image_url     string     `json:"image_url" gorm:"type:varchar(256);not null;default:''" example:"/images/menu.jpg"`
	Desc          string     `json:"desc" gorm:"type:varchar(512);not null;default:''" example:"美味的菜品"`
	Nutrition     string     `json:"nutrition" gorm:"type:varchar(512);not null;default:''" example:"{\"protein\": \"10g\", \"carbs\": \"20g\", \"fat\": \"5g\"}"`
	Ingredients   string     `json:"ingredients" gorm:"type:varchar(512);not null;default:''" example:"{\"米\", \"豆腐\"}"`
	IsRecommended int        `json:"is_recommended" gorm:"type:tinyint(1);not null;default:0" example:"1"`                                             // 是否今日推荐：0 否，1 是
	IsArchived    int        `json:"is_archived" gorm:"type:tinyint(1);not null;default:0" example:"0"`                                                // 是否下架：0 上架，1 下架
	ArchiveTime   *time.Time `json:"archive_time" gorm:"type:datetime(3)" example:"2012-1-1"`                                                          // 下架时间
	Status        uint       `json:"status" gorm:"type:int unsigned;not null;default:1" example:"1"`                                                   // 状态：0 删除，1 正常
	CreateTime    time.Time  `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Creation time
	UpdateTime    time.Time  `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Update time
}

type MenuFeedback struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	MenuID     int       `json:"menu_id" gorm:"type:bigint(20) unsigned;not null;default:0" example:"1"` // 菜单 ID
	UserID     string    `json:"user_id" gorm:"type:varchar(64);not null;default:''" example:"user123"`  // 用户	 ID
	Preference uint      `json:"preference" gorm:"type:int unsigned;not null;default:0" example:"1"`     // 状态：0 默认，1 喜欢，2 不喜欢
	Comment    string    `json:"comment" gorm:"type:varchar(128);not null;default:''" example:"非常好吃"`
	Status     uint      `json:"status" gorm:"type:int unsigned;not null;default:0" example:"1"`                                                   // 状态：0 删除，1 正常
	CreateTime time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Creation time
	UpdateTime time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Update time
}

// MenuQueryRequest 菜单查询请求结构体
type MenuQueryRequest struct {
	Name       string `json:"name" example:"紫菜汤"`      // 菜单名称，支持模糊匹配
	PageSize   int    `json:"page_size" example:"10"`  // 每页数量
	PageNumber int    `json:"page_number" example:"0"` // 页码，从零开始
}

// AdminMenuListRequest 管理员菜单查询请求结构体
type AdminMenuListRequest struct {
	Keyword       string `json:"keyword" example:"紫菜汤"`        // 菜单名称、介绍或 ID
	ArchiveStatus string `json:"archive_status" example:"all"` // all, active, archived
	PageSize      int    `json:"page_size" example:"20"`       // 每页数量
	PageNumber    int    `json:"page_number" example:"0"`      // 页码，从零开始
}

// DeleteMenuRequest 删除菜单请求结构体
type DeleteMenuRequest struct {
	ID int `json:"id" example:"1"` // 菜品ID
}

// SetMenuRecommendationRequest 设置今日推荐请求结构体
type SetMenuRecommendationRequest struct {
	ID            int `json:"id" example:"1"`             // 菜品ID
	IsRecommended int `json:"is_recommended" example:"1"` // 是否设为今日推荐：0 否，1 是
}

// ArchiveMenuRequest 下架/上架菜单请求结构体
type ArchiveMenuRequest struct {
	ID         int `json:"id" example:"1"`          // 菜品ID
	IsArchived int `json:"is_archived" example:"1"` // 是否下架：0 上架，1 下架
}

// MenuLikeRequest 菜品点赞请求结构体（用户ID从JWT token中获取）
type MenuLikeRequest struct {
	MenuID int `json:"menu_id" example:"1"` // 菜品ID
}

// MenuCommentRequest 菜品评论请求结构体（用户ID从JWT token中获取）
type MenuCommentRequest struct {
	MenuID  int    `json:"menu_id" example:"1"`        // 菜品ID
	Comment string `json:"comment" example:"非常好吃的菜品！"` // 评论内容
}

// MenuCommentsQueryRequest 获取菜品评论请求结构体
type MenuCommentsQueryRequest struct {
	MenuID     int `json:"menu_id" example:"1"`     // 菜品ID
	PageSize   int `json:"page_size" example:"10"`  // 每页数量
	PageNumber int `json:"page_number" example:"0"` // 页码，从零开始
}

// MenuWithLikes 包含点赞数的菜品信息
type MenuWithLikes struct {
	Menu
	LikeCount int64 `json:"like_count" example:"5"` // 点赞数
}

// MenuWithUserLikes 包含点赞数和用户点赞状态的菜品信息
type MenuWithUserLikes struct {
	Menu
	LikeCount int64 `json:"like_count" example:"5"` // 点赞数
	Liked     bool  `json:"liked" example:"true"`   // 当前用户是否已点赞
}

// MenuByIDRequest 获取单个菜品请求结构体
type MenuByIDRequest struct {
	MenuID int `json:"menu_id" example:"1"` // 菜品ID
}
