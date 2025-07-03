package repo

import (
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"gorm.io/gorm"
)

type aboutDB struct {
	db *gorm.DB
}

func NewAboutDb(db *gorm.DB) *aboutDB {
	return &aboutDB{db: db}
}

func (a *aboutDB) GetDescription() (model.Description, error) {
	m := model.Description{}
	err := a.db.Table("description").Select("content").Where("status = ?", 1).First(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}

func (a *aboutDB) GetActivityList(timestamp int, pageNumber int) ([]model.Activity, error) {
	pageSize := setting.AppSetting.PageSize
	m := []model.Activity{}
	err := a.db.Table("activity").
		Where("status = ? AND create_time > ?", 1, timestamp).
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Order("create_time DESC").
		Find(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}
