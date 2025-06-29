package model

import "time"

type Description struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Content    string    `gorm:"type:varchar(512);not null;default:''" json:"content"`
	Status     uint      `gorm:"type:int unsigned;not null;default:0" json:"status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime(3)" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime(3)" json:"update_time"`
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
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime(3)" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime(3)" json:"update_time"`
}

func (Activity) TableName() string {
	return "activity"
}

type Image struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Address    string    `gorm:"type:varchar(64);not null;default:''" json:"address" example:"/images/example.jpg"`
	Desc       string    `gorm:"type:varchar(512);not null;default:''" json:"desc" example:"A beautiful image"`
	Status     uint      `gorm:"type:int unsigned;not null;default:0" json:"status" example:"1"`                         // Status: 0 deleted, 1 normal
	TopPic     int       `gorm:"type:tinyint(1);not null;default:0" json:"top_pic" example:"1"`                          // Whether it is a top image
	Type       uint      `gorm:"type:int unsigned;not null;default:0" json:"type" example:"0"`                           // Image type: 0 activity, 1 restaurant introduction
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime(3)" json:"create_time" example:"2024-06-01T12:00:00Z"` // Creation time
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime(3)" json:"update_time" example:"2024-06-01T12:30:00Z"` // Update time
}

func (Image) TableName() string {
	return "images"
}
