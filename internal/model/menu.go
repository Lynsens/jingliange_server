package model

import "time"

type Menu struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Name        string    `json:"name" gorm:"type:varchar(64);not null;default:''" example:"紫菜汤"` // 菜品名称
	Image_url   string    `json:"image_url" gorm:"type:varchar(256);not null;default:''" example:"/images/menu.jpg"`
	Desc        string    `json:"desc" gorm:"type:varchar(512);not null;default:''" example:"美味的菜品"`
	Nutrition   string    `json:"nutrition" gorm:"type:varchar(512);not null;default:''" example:"{\"protein\": \"10g\", \"carbs\": \"20g\", \"fat\": \"5g\"}"`
	Ingredients string    `json:"ingredients" gorm:"type:varchar(512);not null;default:''" example:"{\"米\", \"豆腐\"}"`
	Status      uint      `json:"status" gorm:"type:int unsigned;not null;default:1" example:"1"`                                // 状态：0 删除，1 正常
	CreateTime  time.Time `gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Creation time
	UpdateTime  time.Time `gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Update time
}

type MenuFeedback struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	MenuID     int       `json:"menu_id" gorm:"type:bigint(20) unsigned;not null;default:0" example:"1"` // 菜单 ID
	UserID     string    `json:"user_id" gorm:"type:varchar(64);not null;default:''" example:"user123"`  // 用户	 ID
	Preference uint      `json:"preference" gorm:"type:int unsigned;not null;default:0" example:"1"`     // 状态：0 默认，1 喜欢，2 不喜欢
	Comment    string    `json:"comment" gorm:"type:varchar(128);not null;default:''" example:"非常好吃"`
	Status     uint      `json:"status" gorm:"type:int unsigned;not null;default:0" example:"1"`                                // 状态：0 删除，1 正常
	CreateTime time.Time `gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Creation time
	UpdateTime time.Time `gorm:"type:bigint(20);precision:19;scale:0;default:CURRENT_TIMESTAMP(3);not null" example:"2012-1-1"` // Update time
}
