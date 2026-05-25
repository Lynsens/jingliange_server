package repo

import (
	"time"

	"github.com/lynsens/jingliange_server/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuDB struct {
	db *gorm.DB
}

func NewMenuDB(db *gorm.DB) *MenuDB {
	return &MenuDB{db: db}
}

func (m *MenuDB) GetMenuList(pageSize, pageNumber int, name string) ([]model.MenuWithLikes, error) {
	var menus []model.Menu
	offset := pageNumber * pageSize

	query := m.db.Table("menu").Where("status = ? AND is_archived = ?", 1, 0)

	// 添加模糊查询条件 - 搜索菜品名称和描述
	if name != "" {
		query = query.Where("(name LIKE ? OR `desc` LIKE ?)", "%"+name+"%", "%"+name+"%")
	}

	err := query.Order("is_recommended DESC, id ASC").Offset(offset).Limit(pageSize).Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 为每个菜品获取点赞数
	var menusWithLikes []model.MenuWithLikes
	for _, menu := range menus {
		likeCount, err := m.GetMenuLikeCount(menu.ID)
		if err != nil {
			return nil, err
		}

		menuWithLikes := model.MenuWithLikes{
			Menu:      menu,
			LikeCount: likeCount,
		}
		menusWithLikes = append(menusWithLikes, menuWithLikes)
	}

	return menusWithLikes, nil
}

func (m *MenuDB) GetMenuCount(name string) (int64, error) {
	var count int64
	query := m.db.Table("menu").Where("status = ? AND is_archived = ?", 1, 0)

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

func (m *MenuDB) GetAdminMenuList(pageSize, pageNumber int, keyword string, archiveStatus string) ([]model.MenuWithLikes, error) {
	var menus []model.Menu
	offset := pageNumber * pageSize

	query := m.db.Table("menu").Where("status = ?", 1)
	switch archiveStatus {
	case "active":
		query = query.Where("is_archived = ?", 0)
	case "archived":
		query = query.Where("is_archived = ?", 1)
	}

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(CAST(id AS CHAR) = ? OR name LIKE ? OR `desc` LIKE ?)", keyword, like, like)
	}

	err := query.
		Order("is_archived ASC, is_recommended DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	menusWithLikes := make([]model.MenuWithLikes, 0, len(menus))
	for _, menu := range menus {
		likeCount, err := m.GetMenuLikeCount(menu.ID)
		if err != nil {
			return nil, err
		}

		menusWithLikes = append(menusWithLikes, model.MenuWithLikes{
			Menu:      menu,
			LikeCount: likeCount,
		})
	}

	return menusWithLikes, nil
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

	return m.db.Transaction(func(tx *gorm.DB) error {
		// 确保新创建的菜品状态为1（正常）
		menu.Status = 1
		if menu.IsArchived == 1 {
			menu.IsRecommended = 0
			now := time.Now()
			menu.ArchiveTime = &now
		}
		if err := tx.Table("menu").Create(&menu).Error; err != nil {
			return err
		}

		if menu.IsRecommended == 1 {
			if err := tx.Table("menu").
				Where("id <> ? AND status = ? AND is_archived = ?", menu.ID, 1, 0).
				Update("is_recommended", 0).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *MenuDB) UpdateMenu(menu model.Menu) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		isRecommended := menu.IsRecommended
		if menu.IsArchived == 1 {
			isRecommended = 0
		}

		updates := map[string]interface{}{
			"name":           menu.Name,
			"image_url":      menu.Image_url,
			"desc":           menu.Desc,
			"nutrition":      menu.Nutrition,
			"ingredients":    menu.Ingredients,
			"is_recommended": isRecommended,
			"is_archived":    menu.IsArchived,
			"status":         menu.Status,
		}

		if menu.IsArchived == 1 {
			now := time.Now()
			updates["archive_time"] = now
		} else {
			updates["archive_time"] = nil
		}

		if err := tx.Table("menu").Where("id = ?", menu.ID).Updates(updates).Error; err != nil {
			return err
		}

		if isRecommended == 1 {
			return tx.Table("menu").
				Where("id <> ? AND status = ? AND is_archived = ?", menu.ID, 1, 0).
				Update("is_recommended", 0).Error
		}

		return nil
	})
}

func (m *MenuDB) DeleteMenu(id int) error {
	// 软删除：将状态标记为0（删除状态）而不是真正删除记录
	return m.db.Table("menu").Where("id = ?", id).Updates(map[string]interface{}{
		"status":         0,
		"is_recommended": 0,
	}).Error
}

func (m *MenuDB) SetRecommendedMenu(id int, isRecommended int) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if isRecommended == 1 {
			if err := tx.Table("menu").
				Where("status = ? AND is_archived = ?", 1, 0).
				Update("is_recommended", 0).Error; err != nil {
				return err
			}

			return tx.Table("menu").
				Where("id = ? AND status = ? AND is_archived = ?", id, 1, 0).
				Update("is_recommended", 1).Error
		}

		return tx.Table("menu").
			Where("id = ? AND status = ?", id, 1).
			Update("is_recommended", 0).Error
	})
}

func (m *MenuDB) ArchiveMenu(id int, isArchived int) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"is_archived": isArchived,
		}

		if isArchived == 1 {
			updates["is_recommended"] = 0
			updates["archive_time"] = time.Now()
		} else {
			updates["archive_time"] = nil
		}

		return tx.Table("menu").
			Where("id = ? AND status = ?", id, 1).
			Updates(updates).Error
	})
}

func (m *MenuDB) LikeMenu(menuID int, userID string) error {
	// 检查用户是否已经对该菜品有反馈
	var existingFeedback model.MenuFeedback
	err := m.db.Table("menu_feedback").Where("menu_id = ? AND user_id = ? AND status = ?", menuID, userID, 1).First(&existingFeedback).Error

	if err == nil {
		// 如果已存在反馈，更新为喜欢
		err = m.db.Table("menu_feedback").Where("id = ?", existingFeedback.ID).Update("preference", 1).Error
		return err
	} else if err == gorm.ErrRecordNotFound {
		// 如果不存在反馈，创建新的喜欢记录
		feedback := model.MenuFeedback{
			MenuID:     int(menuID),
			UserID:     userID,
			Preference: 1, // 1 表示喜欢
			Status:     1, // 1 表示正常状态
		}
		err = m.db.Table("menu_feedback").Create(&feedback).Error
		return err
	} else {
		// 其他数据库错误
		return err
	}
}

func (m *MenuDB) UnlikeMenu(menuID int, userID string) error {
	// 将用户对该菜品的反馈设置为默认状态（取消点赞）
	err := m.db.Table("menu_feedback").Where("menu_id = ? AND user_id = ? AND status = ?", menuID, userID, 1).Update("preference", 0).Error
	return err
}

func (m *MenuDB) GetMenuLikeStatus(menuID int, userID string) (int, error) {
	var feedback model.MenuFeedback
	err := m.db.Table("menu_feedback").Where("menu_id = ? AND user_id = ? AND status = ?", menuID, userID, 1).First(&feedback).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil // 0 表示默认状态（未点赞）
	} else if err != nil {
		return 0, err
	}

	return int(feedback.Preference), nil
}

func (m *MenuDB) GetMenuLikeCount(menuID int) (int64, error) {
	var count int64
	err := m.db.Table("menu_feedback").Where("menu_id = ? AND preference = ? AND status = ?", menuID, 1, 1).Count(&count).Error
	return count, err
}

func (m *MenuDB) CommentMenu(menuID int, userID string, comment string, userNickname string, userAvatarURL string) error {
	// 检查用户是否已经对该菜品有反馈
	var existingFeedback model.MenuFeedback
	err := m.db.Table("menu_feedback").Where("menu_id = ? AND user_id = ? AND status = ?", menuID, userID, 1).First(&existingFeedback).Error

	if err == nil {
		// 如果已存在反馈，更新评论
		err = m.db.Table("menu_feedback").Where("id = ?", existingFeedback.ID).Updates(map[string]interface{}{
			"comment":         comment,
			"user_nickname":   userNickname,
			"user_avatar_url": userAvatarURL,
		}).Error
		return err
	} else if err == gorm.ErrRecordNotFound {
		// 如果不存在反馈，创建新的评论记录
		feedback := model.MenuFeedback{
			MenuID:        int(menuID),
			UserID:        userID,
			UserNickname:  userNickname,
			UserAvatarURL: userAvatarURL,
			Preference:    0, // 0 表示默认状态
			Comment:       comment,
			Status:        1, // 1 表示正常状态
		}
		err = m.db.Table("menu_feedback").Create(&feedback).Error
		return err
	} else {
		// 其他数据库错误
		return err
	}
}

func (m *MenuDB) GetMenuComments(menuID int, pageSize, pageNumber int, currentUserID string) ([]model.MenuFeedback, error) {
	var comments []model.MenuFeedback
	offset := pageNumber * pageSize

	query := m.db.Table("menu_feedback").
		Where("menu_id = ? AND status = ? AND comment != ?", menuID, 1, "").
		Offset(offset).
		Limit(pageSize)

	if currentUserID != "" {
		query = query.Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "CASE WHEN user_id = ? THEN 0 ELSE 1 END",
				Vars:               []interface{}{currentUserID},
				WithoutParentheses: true,
			},
		})
	}

	err := query.
		Order("update_time DESC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	for i := range comments {
		comments[i].IsMine = currentUserID != "" && comments[i].UserID == currentUserID
	}
	if err := m.fillCommentLikeState(comments, currentUserID); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m *MenuDB) fillCommentLikeState(comments []model.MenuFeedback, currentUserID string) error {
	if len(comments) == 0 {
		return nil
	}

	commentIDs := make([]int, 0, len(comments))
	indexByID := make(map[int]int, len(comments))
	for i, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
		indexByID[comment.ID] = i
	}

	type likeCountRow struct {
		CommentID int
		Count     int64
	}
	var counts []likeCountRow
	if err := m.db.Table("menu_comment_like").
		Select("comment_id, COUNT(*) AS count").
		Where("comment_id IN ? AND status = ?", commentIDs, 1).
		Group("comment_id").
		Scan(&counts).Error; err != nil {
		return err
	}
	for _, row := range counts {
		if index, ok := indexByID[row.CommentID]; ok {
			comments[index].LikeCount = row.Count
		}
	}

	if currentUserID == "" {
		return nil
	}

	var likedRows []model.MenuCommentLike
	if err := m.db.Table("menu_comment_like").
		Where("comment_id IN ? AND user_id = ? AND status = ?", commentIDs, currentUserID, 1).
		Find(&likedRows).Error; err != nil {
		return err
	}
	for _, row := range likedRows {
		if index, ok := indexByID[row.CommentID]; ok {
			comments[index].Liked = true
		}
	}

	return nil
}

func (m *MenuDB) GetMenuCommentCount(menuID int) (int64, error) {
	var count int64
	err := m.db.Table("menu_feedback").
		Where("menu_id = ? AND status = ? AND comment != ?", menuID, 1, "").
		Count(&count).Error
	return count, err
}

func (m *MenuDB) GetAdminCommentList(pageSize, pageNumber int, keyword string) ([]model.AdminCommentItem, error) {
	var comments []model.AdminCommentItem
	offset := pageNumber * pageSize

	query := m.db.Table("menu_feedback AS f").
		Select("f.id, f.menu_id, COALESCE(m.name, '') AS menu_name, f.user_id, f.user_nickname, f.user_avatar_url, f.comment, f.preference, f.create_time, f.update_time").
		Joins("LEFT JOIN menu AS m ON m.id = f.menu_id").
		Where("f.status = ? AND f.comment != ?", 1, "")

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(CAST(f.id AS CHAR) = ? OR CAST(f.menu_id AS CHAR) = ? OR f.user_id LIKE ? OR f.user_nickname LIKE ? OR f.comment LIKE ? OR m.name LIKE ?)", keyword, keyword, like, like, like, like)
	}

	err := query.
		Order("f.update_time DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&comments).Error
	return comments, err
}

func (m *MenuDB) ClearMenuComment(id int) error {
	return m.db.Table("menu_feedback").
		Where("id = ? AND status = ?", id, 1).
		Update("comment", "").Error
}

func (m *MenuDB) ClearOwnMenuComment(id int, userID string) error {
	result := m.db.Table("menu_feedback").
		Where("id = ? AND user_id = ? AND status = ?", id, userID, 1).
		Update("comment", "")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *MenuDB) ToggleMenuCommentLike(commentID int, userID string) (bool, error) {
	returnedLiked := false
	err := m.db.Transaction(func(tx *gorm.DB) error {
		var comment model.MenuFeedback
		if err := tx.Table("menu_feedback").
			Where("id = ? AND status = ? AND comment != ?", commentID, 1, "").
			First(&comment).Error; err != nil {
			return err
		}

		var existing model.MenuCommentLike
		err := tx.Table("menu_comment_like").
			Where("comment_id = ? AND user_id = ?", commentID, userID).
			First(&existing).Error
		if err == nil {
			newStatus := uint(1)
			if existing.Status == 1 {
				newStatus = 0
			}
			returnedLiked = newStatus == 1
			return tx.Table("menu_comment_like").
				Where("id = ?", existing.ID).
				Update("status", newStatus).Error
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		returnedLiked = true
		return tx.Table("menu_comment_like").Create(&model.MenuCommentLike{
			CommentID: commentID,
			UserID:    userID,
			Status:    1,
		}).Error
	})

	return returnedLiked, err
}

func (m *MenuDB) GetMenuByIDWithLikes(id int) (model.MenuWithLikes, error) {
	var menu model.Menu
	err := m.db.Table("menu").Where("id = ? AND status = ? AND is_archived = ?", id, 1, 0).First(&menu).Error
	if err != nil {
		return model.MenuWithLikes{}, err
	}

	// 获取点赞数
	likeCount, err := m.GetMenuLikeCount(id)
	if err != nil {
		return model.MenuWithLikes{}, err
	}

	menuWithLikes := model.MenuWithLikes{
		Menu:      menu,
		LikeCount: likeCount,
	}

	return menuWithLikes, nil
}

// GetMenuListWithUserLikes 获取包含用户点赞状态的菜单列表
func (m *MenuDB) GetMenuListWithUserLikes(pageSize, pageNumber int, name string, userID string) ([]model.MenuWithUserLikes, error) {
	var menus []model.Menu
	offset := pageNumber * pageSize

	query := m.db.Table("menu").Where("status = ? AND is_archived = ?", 1, 0)

	// 添加模糊查询条件 - 搜索菜品名称和描述
	if name != "" {
		query = query.Where("(name LIKE ? OR `desc` LIKE ?)", "%"+name+"%", "%"+name+"%")
	}

	err := query.Order("is_recommended DESC, id ASC").Offset(offset).Limit(pageSize).Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 为每个菜品获取点赞数和用户点赞状态
	var menusWithUserLikes []model.MenuWithUserLikes
	for _, menu := range menus {
		likeCount, err := m.GetMenuLikeCount(menu.ID)
		if err != nil {
			return nil, err
		}

		// 检查用户是否已点赞
		var liked bool
		if userID != "" {
			preference, err := m.GetMenuLikeStatus(menu.ID, userID)
			if err != nil {
				// 如果查询出错，默认为未点赞
				liked = false
			} else {
				liked = (preference == 1)
			}
		} else {
			liked = false
		}

		menuWithUserLikes := model.MenuWithUserLikes{
			Menu:      menu,
			LikeCount: likeCount,
			Liked:     liked,
		}
		menusWithUserLikes = append(menusWithUserLikes, menuWithUserLikes)
	}

	return menusWithUserLikes, nil
}
