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
		Order("is_top DESC, create_time DESC").
		Find(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}

func (a *aboutDB) GetAdminActivityList(keyword string, pageSize int, pageNumber int) ([]model.Activity, error) {
	if pageSize <= 0 {
		pageSize = setting.AppSetting.PageSize
	}
	if pageNumber < 0 {
		pageNumber = 0
	}

	m := []model.Activity{}
	query := a.db.Table("activity").Where("status = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR content LIKE ? OR place LIKE ?)", like, like, like)
	}

	err := query.
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Order("is_top DESC, create_time DESC").
		Find(&m).Error
	return m, err
}

func (a *aboutDB) GetActivityByID(id uint64) (model.Activity, error) {
	m := model.Activity{}
	err := a.db.Table("activity").Where("id = ?", id).First(&m).Error
	return m, err
}

func (a *aboutDB) CreateActivity(activity model.Activity) error {
	activity.Status = 1
	return a.db.Table("activity").Create(&activity).Error
}

func (a *aboutDB) UpdateActivity(activity model.Activity) error {
	return a.db.Table("activity").
		Where("id = ? AND status = ?", activity.ID, 1).
		Updates(map[string]interface{}{
			"title":      activity.Title,
			"content":    activity.Content,
			"event_time": activity.EventTime,
			"place":      activity.Place,
			"is_top":     activity.IsTop,
		}).Error
}

func (a *aboutDB) DeleteActivity(id uint64) error {
	return a.db.Table("activity").Where("id = ?", id).Updates(map[string]interface{}{
		"status": 0,
		"is_top": 0,
	}).Error
}

func (a *aboutDB) SetActivityTop(id uint64, isTop int) error {
	return a.db.Table("activity").
		Where("id = ? AND status = ?", id, 1).
		Update("is_top", isTop).Error
}

func (a *aboutDB) CreateSuggestion(suggestion model.Suggestion) error {
	suggestion.Status = 1
	suggestion.HandleStatus = 0
	return a.db.Table("suggestion").Create(&suggestion).Error
}

func (a *aboutDB) GetAdminSuggestionList(keyword string, handleStatus string, pageSize int, pageNumber int) ([]model.Suggestion, error) {
	if pageSize <= 0 {
		pageSize = setting.AppSetting.PageSize
	}
	if pageNumber < 0 {
		pageNumber = 0
	}

	suggestions := []model.Suggestion{}
	query := a.db.Table("suggestion").Where("status = ?", 1)
	switch handleStatus {
	case "pending":
		query = query.Where("handle_status = ?", 0)
	case "handled":
		query = query.Where("handle_status = ?", 1)
	}

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(content LIKE ? OR contact LIKE ? OR user_id LIKE ? OR user_nickname LIKE ?)", like, like, like, like)
	}

	err := query.
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Order("handle_status ASC, create_time DESC").
		Find(&suggestions).Error
	return suggestions, err
}

func (a *aboutDB) SetSuggestionHandleStatus(id uint64, handleStatus uint) error {
	return a.db.Table("suggestion").
		Where("id = ? AND status = ?", id, 1).
		Update("handle_status", handleStatus).Error
}
