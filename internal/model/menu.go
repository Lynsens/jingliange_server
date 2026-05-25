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
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	MenuID        int       `json:"menu_id" gorm:"type:bigint(20) unsigned;not null;default:0" example:"1"` // 菜单 ID
	UserID        string    `json:"user_id" gorm:"type:varchar(64);not null;default:''" example:"user123"`  // 用户 ID
	UserNickname  string    `json:"user_nickname" gorm:"type:varchar(64);not null;default:''" example:"莲友"`
	UserAvatarURL string    `json:"user_avatar_url" gorm:"type:varchar(256);not null;default:''" example:"/uploads/images/avatar.jpg"`
	Preference    uint      `json:"preference" gorm:"type:int unsigned;not null;default:0" example:"1"` // 状态：0 默认，1 喜欢，2 不喜欢
	Comment       string    `json:"comment" gorm:"type:varchar(128);not null;default:''" example:"非常好吃"`
	Status        uint      `json:"status" gorm:"type:int unsigned;not null;default:0" example:"1"`                                                   // 状态：0 删除，1 正常
	CreateTime    time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Creation time
	UpdateTime    time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Update time
	IsMine        bool      `json:"is_mine" gorm:"-" example:"true"`
}

// AdminCommentItem 管理员评论列表项
type AdminCommentItem struct {
	ID            int       `json:"id"`
	MenuID        int       `json:"menu_id"`
	MenuName      string    `json:"menu_name"`
	UserID        string    `json:"user_id"`
	UserNickname  string    `json:"user_nickname"`
	UserAvatarURL string    `json:"user_avatar_url"`
	Comment       string    `json:"comment"`
	Preference    uint      `json:"preference"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
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

// AdminCommentListRequest 管理员评论查询请求结构体
type AdminCommentListRequest struct {
	Keyword    string `json:"keyword" example:"好吃"`    // 评论、用户 ID、菜品名称或 ID
	PageSize   int    `json:"page_size" example:"20"`  // 每页数量
	PageNumber int    `json:"page_number" example:"0"` // 页码，从零开始
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

// DeleteCommentRequest 删除评论请求结构体
type DeleteCommentRequest struct {
	ID int `json:"id" example:"1"` // 评论反馈 ID
}

// MenuLikeRequest 菜品点赞请求结构体（用户ID从JWT token中获取）
type MenuLikeRequest struct {
	MenuID int `json:"menu_id" example:"1"` // 菜品ID
}

// MenuCommentRequest 菜品评论请求结构体（用户ID从JWT token中获取）
type MenuCommentRequest struct {
	MenuID        int    `json:"menu_id" example:"1"`        // 菜品ID
	Comment       string `json:"comment" example:"非常好吃的菜品！"` // 评论内容
	UserNickname  string `json:"user_nickname" example:"莲友"` // 评论展示昵称
	UserAvatarURL string `json:"user_avatar_url" example:"/uploads/images/avatar.jpg"`
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

type ComboRecommendation struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Title       string    `json:"title" gorm:"type:varchar(64);not null;default:''" example:"清淡养胃套餐"`
	Description string    `json:"description" gorm:"type:varchar(512);not null;default:''" example:"适合午餐，搭配均衡"`
	IsActive    int       `json:"is_active" gorm:"type:tinyint(1);not null;default:0" example:"1"`
	Status      uint      `json:"status" gorm:"type:int unsigned;not null;default:1" example:"1"`
	CreateTime  time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"`
	UpdateTime  time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"`
}

type ComboRecommendationItem struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	ComboID    uint64    `json:"combo_id" gorm:"type:bigint(20) unsigned;not null;default:0" example:"1"`
	MenuID     int       `json:"menu_id" gorm:"type:bigint(20) unsigned;not null;default:0" example:"1"`
	SortOrder  int       `json:"sort_order" gorm:"type:int unsigned;not null;default:0" example:"0"`
	Status     uint      `json:"status" gorm:"type:int unsigned;not null;default:1" example:"1"`
	CreateTime time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"`
	UpdateTime time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"`
}

type ComboRecommendationMenuItem struct {
	MenuWithUserLikes
	ComboItemID uint64 `json:"combo_item_id" example:"1"`
	SortOrder   int    `json:"sort_order" example:"0"`
}

type ComboRecommendationResponse struct {
	ComboRecommendation
	Items []ComboRecommendationMenuItem `json:"items"`
}

type ComboRecommendationRequest struct {
	ID          uint64 `json:"id" example:"1"`
	Title       string `json:"title" example:"清淡养胃套餐"`
	Description string `json:"description" example:"适合午餐，搭配均衡"`
	IsActive    int    `json:"is_active" example:"1"`
	MenuIDs     []int  `json:"menu_ids" example:"1,2,3"`
}

type ComboRecommendationListRequest struct {
	Keyword    string `json:"keyword" example:"清淡"`
	PageSize   int    `json:"page_size" example:"20"`
	PageNumber int    `json:"page_number" example:"0"`
}

type ComboRecommendationIDRequest struct {
	ID uint64 `json:"id" example:"1"`
}

type ComboRecommendationActiveRequest struct {
	ID       uint64 `json:"id" example:"1"`
	IsActive int    `json:"is_active" example:"1"`
}

// MenuByIDRequest 获取单个菜品请求结构体
type MenuByIDRequest struct {
	MenuID int `json:"menu_id" example:"1"` // 菜品ID
}
