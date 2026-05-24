package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"gorm.io/gorm"
)

func normalizeBinaryFlag(value int) int {
	if value == 1 {
		return 1
	}
	return 0
}

func validateArchiveStatus(value string) string {
	switch value {
	case "active", "archived":
		return value
	default:
		return "all"
	}
}

// @Summary 管理员获取菜单列表
// @Description 管理员菜单列表，可查看上架中和已下架菜品。
// @Tags Menu
// @Accept json
// @Param query body model.AdminMenuListRequest true "查询参数" schemaexample({"keyword":"","archive_status":"all","page_size":20,"page_number":0})
// @Produce json
// @Success 200 {object} app.Response{data=[]model.MenuWithLikes}
// @Router /api/admin/menu/list [post]
func GetMenuItems(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	var req model.AdminMenuListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.PageSize <= 0 {
		req.PageSize = setting.AppSetting.PageSize
	}
	if req.PageNumber < 0 {
		req.PageNumber = 0
	}
	req.ArchiveStatus = validateArchiveStatus(req.ArchiveStatus)

	menus, err := menuRepo.GetAdminMenuList(req.PageSize, req.PageNumber, strings.TrimSpace(req.Keyword), req.ArchiveStatus)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, menus)
}

// @Summary 管理员获取评论列表
// @Description 管理员评论列表，可按评论内容、用户 ID、菜品名称或 ID 搜索。
// @Tags Menu
// @Accept json
// @Param query body model.AdminCommentListRequest true "查询参数" schemaexample({"keyword":"","page_size":20,"page_number":0})
// @Produce json
// @Success 200 {object} app.Response{data=[]model.AdminCommentItem}
// @Router /api/admin/comment/list [post]
func GetComments(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	var req model.AdminCommentListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.PageSize <= 0 {
		req.PageSize = setting.AppSetting.PageSize
	}
	if req.PageNumber < 0 {
		req.PageNumber = 0
	}

	comments, err := menuRepo.GetAdminCommentList(req.PageSize, req.PageNumber, strings.TrimSpace(req.Keyword))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, comments)
}

// @Summary 管理员删除评论
// @Description 管理员删除评论仅清空评论内容，不影响用户点赞状态。
// @Tags Menu
// @Accept json
// @Param delete body model.DeleteCommentRequest true "删除参数" schemaexample({"id":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/comment/delete [delete]
func DeleteComment(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	var req model.DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid comment ID")
		return
	}

	if err := menuRepo.ClearMenuComment(req.ID); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 上传菜品
// @Description 提供菜品的名称和简介，图片url、营养价值表、主要成分可选。营养价值表和主要成分需要以JSON字符串格式提供。
// @Tags Menu
// @Accept json
// @Param menu body model.Menu true "菜品信息" schemaexample({"name":"豆腐汤","image_url":"/images/tofusoup.jpg","desc":"清淡营养的素食汤品，豆腐嫩滑，口感清香，富含植物蛋白","nutrition":"{\"calories\":\"120kcal\",\"protein\":\"12g\",\"carbs\":\"8g\",\"fat\":\"6g\",\"fiber\":\"2g\",\"sodium\":\"600mg\"}","ingredients":"{\"tofu\":\"200g\",\"seaweed\":\"20g\",\"green_onion\":\"10g\",\"mushroom\":\"50g\",\"soy_sauce\":\"10ml\",\"sesame_oil\":\"5ml\",\"salt\":\"3g\",\"white_pepper\":\"1g\",\"vegetable_broth\":\"400ml\"}"})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":{"id":1,"name":"豆腐汤","image_url":"/images/tofusoup.jpg","desc":"清淡营养的素食汤品，豆腐嫩滑，口感清香，富含植物蛋白","nutrition":"...","ingredients":"...","status":1}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Name and description are required"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/admin/uploadMenuItem [post]
func UploadMenuItem(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewMenuDB(db)

	menuItem := model.Menu{
		Status: 1, // 默认状态为1（正常）
	}

	if err := c.ShouldBindJSON(&menuItem); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if menuItem.Name == "" || menuItem.Desc == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Name and description are required")
		return
	}
	menuItem.IsRecommended = normalizeBinaryFlag(menuItem.IsRecommended)
	menuItem.IsArchived = normalizeBinaryFlag(menuItem.IsArchived)

	if err := repo.CreateMenu(menuItem); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			appG.Response(http.StatusConflict, e.INVALID_PARAMS, "Menu item with this name already exists")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, menuItem)
}

// @Summary 更新菜品
// @Description 根据菜品ID更新名称、简介、图片url、营养价值表和主要成分。图片url、营养价值表和主要成分可选，营养价值表和主要成分需要以JSON字符串格式提供。
// @Tags Menu
// @Accept json
// @Param menu body model.Menu true "菜品信息"
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":{"id":1,"name":"豆腐汤","image_url":"/images/tofusoup.jpg","desc":"清淡营养的素食汤品","nutrition":"...","ingredients":"...","status":1}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Name and description are required"}"
// @Failure 404 {object} app.Response "{"code":500,"msg":"fail","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/admin/updateMenuItem [put]
func UpdateMenuItem(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewMenuDB(db)

	var menuItem model.Menu
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if menuItem.ID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	if menuItem.Name == "" || menuItem.Desc == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Name and description are required")
		return
	}
	menuItem.IsRecommended = normalizeBinaryFlag(menuItem.IsRecommended)
	menuItem.IsArchived = normalizeBinaryFlag(menuItem.IsArchived)
	if menuItem.IsArchived == 1 && menuItem.IsRecommended == 1 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Archived menu item cannot be recommended")
		return
	}

	existingMenu, err := repo.GetMenuByID(menuItem.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	if existingMenu.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item already deleted")
		return
	}

	menuItem.Status = 1
	if err := repo.UpdateMenu(menuItem); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if menuItem.IsRecommended == 1 {
		if err := repo.SetRecommendedMenu(menuItem.ID, 1); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, menuItem)
}

// @Summary 设置今日推荐菜品
// @Description 设置或取消菜品今日推荐。设置为推荐时会自动取消其它菜品的今日推荐，保证最多只有一个今日推荐菜品。
// @Tags Menu
// @Accept json
// @Param recommend body model.SetMenuRecommendationRequest true "今日推荐参数" schemaexample({"id":1,"is_recommended":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid menu ID"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/admin/recommendMenuItem [put]
func SetRecommendedMenuItem(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	var req model.SetMenuRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	if req.IsRecommended != 0 && req.IsRecommended != 1 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid recommendation value")
		return
	}

	existingMenu, err := menuRepo.GetMenuByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	if existingMenu.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item already deleted")
		return
	}
	if existingMenu.IsArchived == 1 && req.IsRecommended == 1 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Archived menu item cannot be recommended")
		return
	}

	if err := menuRepo.SetRecommendedMenu(req.ID, req.IsRecommended); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 下架或重新上架菜品
// @Description 设置菜品 archive 状态。下架今日推荐菜品时会自动取消今日推荐。
// @Tags Menu
// @Accept json
// @Param archive body model.ArchiveMenuRequest true "下架参数" schemaexample({"id":1,"is_archived":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/archiveMenuItem [put]
func ArchiveMenuItem(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	var req model.ArchiveMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}
	req.IsArchived = normalizeBinaryFlag(req.IsArchived)

	existingMenu, err := menuRepo.GetMenuByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	if existingMenu.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item already deleted")
		return
	}

	if err := menuRepo.ArchiveMenu(req.ID, req.IsArchived); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除菜品
// @Description 根据菜品ID删除菜品（软删除，将状态标记为0）
// @Tags Menu
// @Accept json
// @Param delete body model.DeleteMenuRequest true "删除参数" schemaexample({"id":1})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid menu ID"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/admin/deleteMenuItem [delete]
func DeleteMenuItem(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewMenuDB(db)

	// 使用 body 参数
	var deleteReq model.DeleteMenuRequest
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	// 验证ID
	if deleteReq.ID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	// 检查菜品是否存在且状态为正常
	existingMenu, err := repo.GetMenuByID(deleteReq.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	// 检查菜品是否已经被删除
	if existingMenu.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item already deleted")
		return
	}

	// 执行软删除
	if err := repo.DeleteMenu(deleteReq.ID); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
