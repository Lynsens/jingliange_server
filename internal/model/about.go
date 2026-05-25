package model

import "time"

type Description struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Content    string    `gorm:"type:varchar(512);not null;default:''" json:"content"`
	Status     uint      `gorm:"type:int unsigned;not null;default:0" json:"status"`
	CreateTime time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Creation time
	UpdateTime time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Update time
}

func (Description) TableName() string {
	return "description"
}

type Activity struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(64);not null;default:''" json:"title"`
	Content    string    `gorm:"type:varchar(512);not null;default:''" json:"content"`
	EventTime  string    `gorm:"type:varchar(64);not null;default:''" json:"event_time"`
	Place      string    `gorm:"type:varchar(128);not null;default:''" json:"place"`
	Img        string    `gorm:"type:varchar(32);not null;default:''" json:"img"` // Image address
	IsTop      int       `gorm:"type:tinyint(1);not null;default:0" json:"is_top"`
	Status     uint      `gorm:"type:int unsigned;not null;default:0" json:"status"`
	CreateTime time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Creation time
	UpdateTime time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Update time
}

func (Activity) TableName() string {
	return "activity"
}

type Image struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Address    string    `gorm:"type:varchar(256);not null;default:''" json:"address" example:"/images/example.jpg"`
	Desc       string    `gorm:"type:varchar(512);not null;default:''" json:"desc" example:"A beautiful image"`
	Status     uint      `gorm:"type:int unsigned;not null;default:0" json:"status" example:"1"`                                // Status: 0 deleted, 1 normal
	TopPic     int       `gorm:"type:tinyint(1);not null;default:0" json:"top_pic" example:"1"`                                 // Whether it is a top image
	Type       uint      `gorm:"type:int unsigned;not null;default:0" json:"type" example:"0"`                                  // Image type: 0 activity, 1 restaurant introduction
	CreateTime time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Creation time
	UpdateTime time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Update time
}

func (Image) TableName() string {
	return "images"
}

type Suggestion struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	UserID        string    `gorm:"type:varchar(64);not null;default:''" json:"user_id" example:"user123"`
	UserNickname  string    `gorm:"type:varchar(64);not null;default:''" json:"user_nickname" example:"莲友"`
	UserAvatarURL string    `gorm:"type:varchar(256);not null;default:''" json:"user_avatar_url" example:"/uploads/images/avatar.jpg"`
	Content       string    `gorm:"type:varchar(1024);not null;default:''" json:"content" example:"希望增加更多清淡菜品"`
	Contact       string    `gorm:"type:varchar(128);not null;default:''" json:"contact" example:"微信或手机号"`
	HandleStatus  uint      `gorm:"type:int unsigned;not null;default:0" json:"handle_status" example:"0"`                         // 处理状态：0 未处理，1 已处理
	Status        uint      `gorm:"type:int unsigned;not null;default:1" json:"status" example:"1"`                                // 记录状态：0 删除，1 正常
	CreateTime    time.Time `json:"create_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Creation time
	UpdateTime    time.Time `json:"update_time" gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null"` // Update time
}

func (Suggestion) TableName() string {
	return "suggestion"
}

// ActivityQueryRequest 活动查询请求结构体
type ActivityQueryRequest struct {
	Timestamp  int `json:"timestamp" example:"0"`   // 时间戳
	PageNumber int `json:"page_number" example:"0"` // 页码，从零开始
}

// AdminActivityListRequest 管理员活动列表请求结构体
type AdminActivityListRequest struct {
	Keyword    string `json:"keyword" example:"茶会"`    // 标题、地点、内容关键词
	PageSize   int    `json:"page_size" example:"20"`  // 每页数量
	PageNumber int    `json:"page_number" example:"0"` // 页码，从零开始
}

// AdminActivityRequest 管理员活动新增/更新请求结构体
type AdminActivityRequest struct {
	ID        uint64 `json:"id" example:"1"`
	Title     string `json:"title" example:"莲花茶会"`
	Content   string `json:"content" example:"净心茶会与素食交流"`
	EventTime string `json:"event_time" example:"周日 14:00"`
	Place     string `json:"place" example:"二楼茶室"`
	IsTop     int    `json:"is_top" example:"1"`
}

// AdminActivityIDRequest 管理员活动 ID 请求结构体
type AdminActivityIDRequest struct {
	ID uint64 `json:"id" example:"1"`
}

// AdminActivityTopRequest 管理员活动置顶请求结构体
type AdminActivityTopRequest struct {
	ID    uint64 `json:"id" example:"1"`
	IsTop int    `json:"is_top" example:"1"`
}

// SuggestionCreateRequest 建议箱提交请求结构体
type SuggestionCreateRequest struct {
	Content       string `json:"content" example:"希望增加更多清淡菜品"`
	Contact       string `json:"contact" example:"微信或手机号"`
	UserNickname  string `json:"user_nickname" example:"莲友"`
	UserAvatarURL string `json:"user_avatar_url" example:"/uploads/images/avatar.jpg"`
}

// AdminSuggestionListRequest 管理员建议列表请求结构体
type AdminSuggestionListRequest struct {
	Keyword      string `json:"keyword" example:"清淡"`
	HandleStatus string `json:"handle_status" example:"all"` // all, pending, handled
	PageSize     int    `json:"page_size" example:"20"`
	PageNumber   int    `json:"page_number" example:"0"`
}

// AdminSuggestionStatusRequest 管理员建议处理状态请求结构体
type AdminSuggestionStatusRequest struct {
	ID           uint64 `json:"id" example:"1"`
	HandleStatus uint   `json:"handle_status" example:"1"` // 0 未处理，1 已处理
}
