package repo

import (
	"github.com/lynsens/jingliange_server/internal/model"
	"gorm.io/gorm"
)

type MenuDB struct {
	db *gorm.DB
}

func NewMenuDB(db *gorm.DB) *MenuDB {
	return &MenuDB{db: db}
}

func (m *MenuDB) GetMenuList(pageSize, pageNumber int, name string) ([]model.Menu, error) {
	var menus []model.Menu
	offset := pageNumber * pageSize

	query := m.db.Table("menu").Where("status = ?", 1) // 只查询状态为1的菜品

	// 添加模糊查询条件 - 搜索菜品名称和描述
	if name != "" {
		query = query.Where("(name LIKE ? OR `desc` LIKE ?)", "%"+name+"%", "%"+name+"%")
	}

	err := query.Offset(offset).Limit(pageSize).Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *MenuDB) GetMenuCount(name string) (int64, error) {
	var count int64
	query := m.db.Table("menu").Where("status = ?", 1) // 只统计状态为1的菜品

	// 添加模糊查询条件 - 搜索菜品名称和描述
	if name != "" {
		query = query.Where("(name LIKE ? OR `desc` LIKE ?)", "%"+name+"%", "%"+name+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *MenuDB) GetMenuByID(id int) (model.Menu, error) {
	var menu model.Menu
	err := m.db.Table("menu").Where("id = ?", id).First(&menu).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func (m *MenuDB) CreateMenu(menu model.Menu) error {
	// 检查是否已存在相同名称的菜品
	var existingMenu model.Menu
	err := m.db.Table("menu").Where("name = ? AND status = ?", menu.Name, 1).First(&existingMenu).Error
	if err == nil {
		// 找到了相同名称的菜品，返回错误
		return gorm.ErrDuplicatedKey
	}
	if err != gorm.ErrRecordNotFound {
		// 其他数据库错误
		return err
	}

	// 确保新创建的菜品状态为1（正常）
	menu.Status = 1
	err = m.db.Table("menu").Create(&menu).Error // 不使用 Table("menu")，直接使用模型
	if err != nil {
		return err
	}
	return nil
}

func (m *MenuDB) UpdateMenu(menu model.Menu) error {
	err := m.db.Table("menu").Where("id = ?", menu.ID).Updates(menu).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *MenuDB) DeleteMenu(id int) error {
	// 软删除：将状态标记为0（删除状态）而不是真正删除记录
	err := m.db.Table("menu").Where("id = ?", id).Update("status", 0).Error
	if err != nil {
		return err
	}
	return nil
}
