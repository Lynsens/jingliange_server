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

func normalizeComboActive(value int) int {
	if value == 1 {
		return 1
	}
	return 0
}

func validateComboPayload(req model.ComboRecommendationRequest) string {
	if strings.TrimSpace(req.Title) == "" {
		return "Combo title is required"
	}
	if len(req.MenuIDs) == 0 {
		return "Combo menu items are required"
	}
	if len(req.MenuIDs) > 4 {
		return "Combo supports at most 4 menu items"
	}
	return ""
}

func getComboRepo(c *gin.Context, appG app.Gin) (*repo.MenuDB, bool) {
	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return nil, false
	}

	return repo.NewMenuDB(db), true
}

// @Summary 管理员获取套餐推荐列表
// @Description 管理员查看套餐推荐列表。
// @Tags Combo
// @Accept json
// @Param query body model.ComboRecommendationListRequest true "查询参数" schemaexample({"keyword":"","page_size":20,"page_number":0})
// @Produce json
// @Success 200 {object} app.Response{data=[]model.ComboRecommendationResponse}
// @Router /api/admin/combo/list [post]
func GetComboRecommendations(c *gin.Context) {
	appG := app.Gin{C: c}
	menuRepo, ok := getComboRepo(c, appG)
	if !ok {
		return
	}

	var req model.ComboRecommendationListRequest
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

	combos, err := menuRepo.GetAdminComboRecommendationList(strings.TrimSpace(req.Keyword), req.PageSize, req.PageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, combos)
}

// @Summary 管理员创建套餐推荐
// @Description 创建一个套餐推荐，最多包含 4 个上架菜品。
// @Tags Combo
// @Accept json
// @Param combo body model.ComboRecommendationRequest true "套餐内容" schemaexample({"title":"清淡养胃套餐","description":"适合午餐","is_active":1,"menu_ids":[1,2,3]})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/combo/create [post]
func CreateComboRecommendation(c *gin.Context) {
	appG := app.Gin{C: c}
	menuRepo, ok := getComboRepo(c, appG)
	if !ok {
		return
	}

	var req model.ComboRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}
	if msg := validateComboPayload(req); msg != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, msg)
		return
	}

	combo := model.ComboRecommendation{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		IsActive:    normalizeComboActive(req.IsActive),
		Status:      1,
	}

	if err := menuRepo.CreateComboRecommendation(combo, req.MenuIDs); err != nil {
		if errors.Is(err, repo.ErrInvalidComboItems) {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid combo menu items")
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 管理员更新套餐推荐
// @Description 更新套餐推荐标题、说明、菜品组合和启用状态。
// @Tags Combo
// @Accept json
// @Param combo body model.ComboRecommendationRequest true "套餐内容" schemaexample({"id":1,"title":"清淡养胃套餐","description":"适合午餐","is_active":1,"menu_ids":[1,2,3]})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/combo/update [put]
func UpdateComboRecommendation(c *gin.Context) {
	appG := app.Gin{C: c}
	menuRepo, ok := getComboRepo(c, appG)
	if !ok {
		return
	}

	var req model.ComboRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}
	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid combo ID")
		return
	}
	if msg := validateComboPayload(req); msg != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, msg)
		return
	}

	combo := model.ComboRecommendation{
		ID:          req.ID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		IsActive:    normalizeComboActive(req.IsActive),
	}

	if err := menuRepo.UpdateComboRecommendation(combo, req.MenuIDs); err != nil {
		if errors.Is(err, repo.ErrInvalidComboItems) {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid combo menu items")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Combo not found")
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 管理员设置首页套餐推荐
// @Description 设置或取消一个套餐作为首页套餐推荐。
// @Tags Combo
// @Accept json
// @Param active body model.ComboRecommendationActiveRequest true "启用状态" schemaexample({"id":1,"is_active":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/combo/active [put]
func SetComboRecommendationActive(c *gin.Context) {
	appG := app.Gin{C: c}
	menuRepo, ok := getComboRepo(c, appG)
	if !ok {
		return
	}

	var req model.ComboRecommendationActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}
	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid combo ID")
		return
	}

	if err := menuRepo.SetComboRecommendationActive(req.ID, normalizeComboActive(req.IsActive)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Combo not found")
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 管理员删除套餐推荐
// @Description 删除套餐推荐。
// @Tags Combo
// @Accept json
// @Param delete body model.ComboRecommendationIDRequest true "删除参数" schemaexample({"id":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/combo/delete [delete]
func DeleteComboRecommendation(c *gin.Context) {
	appG := app.Gin{C: c}
	menuRepo, ok := getComboRepo(c, appG)
	if !ok {
		return
	}

	var req model.ComboRecommendationIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}
	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid combo ID")
		return
	}

	if err := menuRepo.DeleteComboRecommendation(req.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Combo not found")
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
