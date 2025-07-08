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
	Img        string    `gorm:"type:varchar(32);not null;default:''" json:"img"` // Image address
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

// ActivityQueryRequest 活动查询请求结构体
type ActivityQueryRequest struct {
	Timestamp  int `json:"timestamp" example:"0"`   // 时间戳
	PageNumber int `json:"page_number" example:"0"` // 页码，从零开始
}
