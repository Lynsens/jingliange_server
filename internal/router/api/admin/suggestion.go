package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
)

func normalizeSuggestionHandleStatus(value string) string {
	switch value {
	case "pending", "handled":
		return value
	default:
		return "all"
	}
}

// @Summary 管理员获取建议列表
// @Description 管理员查看用户建议，可按内容、联系方式、用户搜索。
// @Tags Suggestion
// @Accept json
// @Param query body model.AdminSuggestionListRequest true "查询参数" schemaexample({"keyword":"","handle_status":"all","page_size":20,"page_number":0})
// @Produce json
// @Success 200 {object} app.Response{data=[]model.Suggestion}
// @Router /api/admin/suggestion/list [post]
func GetSuggestions(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	aboutRepo := repo.NewAboutDb(db)

	var req model.AdminSuggestionListRequest
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
	req.HandleStatus = normalizeSuggestionHandleStatus(req.HandleStatus)

	suggestions, err := aboutRepo.GetAdminSuggestionList(strings.TrimSpace(req.Keyword), req.HandleStatus, req.PageSize, req.PageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, suggestions)
}

// @Summary 管理员更新建议处理状态
// @Description 管理员将建议标记为未处理或已处理。
// @Tags Suggestion
// @Accept json
// @Param status body model.AdminSuggestionStatusRequest true "处理状态" schemaexample({"id":1,"handle_status":1})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":null}"
// @Router /api/admin/suggestion/status [put]
func UpdateSuggestionStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	aboutRepo := repo.NewAboutDb(db)

	var req model.AdminSuggestionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}
	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid suggestion ID")
		return
	}

	handleStatus := uint(0)
	if req.HandleStatus == 1 {
		handleStatus = 1
	}

	if err := aboutRepo.SetSuggestionHandleStatus(req.ID, handleStatus); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
