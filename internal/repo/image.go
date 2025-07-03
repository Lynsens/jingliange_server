package repo

import "github.com/lynsens/jingliange_server/internal/model"

func (a *aboutDB) SaveImageUrlToDB(url string, desc string, status uint, topPic int, imageType uint) error {
	image := model.Image{
		Address: url,
		Desc:    desc,
		Status:  status,
		TopPic:  topPic,
		Type:    imageType,
		// CreateTime and UpdateTime are handled automatically by GORM
	}

	if err := a.db.Table("images").Create(&image).Error; err != nil {
		return err
	}

	return nil
}

func (a *aboutDB) GetImageList(imageType int, topPic int, pageNumber int, pageSize int) ([]model.Image, error) {
	m := []model.Image{}
	a.db = a.db.Debug()
	query := a.db.Table("images").Where("status = ? AND type = ?", 1, imageType)

	if topPic >= 0 {
		query = query.Where("top_pic = ?", topPic)
	}

	offset := pageNumber * pageSize
	err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}

func (a *aboutDB) GetTopImage() (model.Image, error) {
	m := model.Image{}
	err := a.db.Table("images").
		Where("status = ? AND top_pic = ?", 1, 1).
		Order("create_time DESC").
		First(&m).Error
	if err != nil {
		return m, err
	}
	return m, nil
}
