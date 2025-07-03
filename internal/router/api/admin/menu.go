package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"gorm.io/gorm"
)

// @Summary 上传菜品
// @Description 提供菜品的名称、图片url、简介、营养价值表、主要成分，上传到数据库。营养价值表和主要成分需要以JSON字符串格式提供。
// @Tags Menu
// @Accept json
// @Param menu body model.Menu true "菜品信息" schemaexample({"name":"豆腐汤","image_url":"/images/tofusoup.jpg","desc":"清淡营养的素食汤品，豆腐嫩滑，口感清香，富含植物蛋白","nutrition":"{\"calories\":\"120kcal\",\"protein\":\"12g\",\"carbs\":\"8g\",\"fat\":\"6g\",\"fiber\":\"2g\",\"sodium\":\"600mg\"}","ingredients":"{\"tofu\":\"200g\",\"seaweed\":\"20g\",\"green_onion\":\"10g\",\"mushroom\":\"50g\",\"soy_sauce\":\"10ml\",\"sesame_oil\":\"5ml\",\"salt\":\"3g\",\"white_pepper\":\"1g\",\"vegetable_broth\":\"400ml\"}"})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":{"id":1,"name":"豆腐汤","image_url":"/images/tofusoup.jpg","desc":"清淡营养的素食汤品，豆腐嫩滑，口感清香，富含植物蛋白","nutrition":"...","ingredients":"...","status":1}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"All fields are required"}"
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

	if menuItem.Name == "" || menuItem.Image_url == "" || menuItem.Desc == "" || menuItem.Nutrition == "" || menuItem.Ingredients == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "All fields are required")
		return
	}

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
