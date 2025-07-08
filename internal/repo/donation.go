package repo

import (
	"fmt"

	"github.com/lynsens/jingliange_server/internal/model"
	"gorm.io/gorm"
)

type DonationDB struct {
	db *gorm.DB
}

func NewDonationDB(db *gorm.DB) *DonationDB {
	return &DonationDB{db: db}
}

// GetDonationList 获取功德榜列表
func (d *DonationDB) GetDonationList(req model.DonationQueryRequest) ([]model.Donation, error) {
	var donations []model.Donation
	offset := req.PageNumber * req.PageSize

	query := d.db.Table("donation").Where("is_visible = ?", 1)

	// 年份筛选
	if req.Year > 0 {
		query = query.Where("YEAR(donate_time) = ?", req.Year)
	}

	// 时间段筛选
	if req.Period == "first" {
		// 上半年 1-6月
		query = query.Where("MONTH(donate_time) BETWEEN 1 AND 6")
	} else if req.Period == "second" {
		// 下半年 7-12月
		query = query.Where("MONTH(donate_time) BETWEEN 7 AND 12")
	}
	// req.Period == "all" 不需要额外条件

	// 昵称搜索
	if req.DonorName != "" {
		query = query.Where("donor_name LIKE ?", "%"+req.DonorName+"%")
	}

	// 排序
	orderBy := "donate_time"
	if req.SortBy == "amount" {
		orderBy = "amount"
	}

	sortOrder := "DESC"
	if req.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(fmt.Sprintf("%s %s", orderBy, sortOrder))

	err := query.Offset(offset).Limit(req.PageSize).Find(&donations).Error
	return donations, err
}

// GetDonationCount 获取功德榜总数
func (d *DonationDB) GetDonationCount(req model.DonationQueryRequest) (int64, error) {
	var count int64
	query := d.db.Table("donation").Where("is_visible = ?", 1)

	// 年份筛选
	if req.Year > 0 {
		query = query.Where("YEAR(donate_time) = ?", req.Year)
	}

	// 时间段筛选
	if req.Period == "first" {
		query = query.Where("MONTH(donate_time) BETWEEN 1 AND 6")
	} else if req.Period == "second" {
		query = query.Where("MONTH(donate_time) BETWEEN 7 AND 12")
	}

	// 昵称搜索
	if req.DonorName != "" {
		query = query.Where("donor_name LIKE ?", "%"+req.DonorName+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

// CreateDonation 创建捐款记录
func (d *DonationDB) CreateDonation(donation model.Donation) error {
	return d.db.Table("donation").Create(&donation).Error
}

// GetDonationStats 获取捐款统计
func (d *DonationDB) GetDonationStats(req model.DonationStatsRequest) (model.DonationStats, error) {
	var stats model.DonationStats
	query := d.db.Table("donation").Where("is_visible = ?", 1)

	// 年份筛选
	if req.Year > 0 {
		query = query.Where("YEAR(donate_time) = ?", req.Year)
	}

	// 时间段筛选
	if req.Period == "first" {
		query = query.Where("MONTH(donate_time) BETWEEN 1 AND 6")
	} else if req.Period == "second" {
		query = query.Where("MONTH(donate_time) BETWEEN 7 AND 12")
	}

	// 统计总金额和总人次
	err := query.Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count").
		Scan(&stats).Error

	return stats, err
}

// GetDonationByID 根据ID获取捐款记录
func (d *DonationDB) GetDonationByID(id int) (model.Donation, error) {
	var donation model.Donation
	err := d.db.Table("donation").Where("id = ?", id).First(&donation).Error
	return donation, err
}

// UpdateDonationVisibility 更新捐款记录显示状态（管理员功能）
func (d *DonationDB) UpdateDonationVisibility(id int, isVisible int) error {
	return d.db.Table("donation").Where("id = ?", id).Update("is_visible", isVisible).Error
}

// DeleteDonation 删除捐款记录（管理员功能）
func (d *DonationDB) DeleteDonation(id int) error {
	return d.db.Table("donation").Where("id = ?", id).Delete(&model.Donation{}).Error
}

// UserDB 用户数据库操作
type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{db: db}
}

// CreateOrUpdateUser 创建或更新用户
func (u *UserDB) CreateOrUpdateUser(userID string) error {
	user := model.User{
		ID: userID,
	}

	// 使用 GORM 的 Upsert 功能
	return u.db.Table("user").Where("id = ?", userID).FirstOrCreate(&user).Error
}

// GetUserByID 根据ID获取用户
func (u *UserDB) GetUserByID(userID string) (model.User, error) {
	var user model.User
	err := u.db.Table("user").Where("id = ?", userID).First(&user).Error
	return user, err
}

// ValidateUser 验证用户是否存在
func (u *UserDB) ValidateUser(userID string) bool {
	var count int64
	u.db.Table("user").Where("id = ?", userID).Count(&count)
	return count > 0
}
