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

	if err := a.db.Create(&image).Error; err != nil {
		return err
	}

	return nil
}
