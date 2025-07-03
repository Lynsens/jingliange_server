package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
)

// @Summary 获取净莲阁的菜单
// @Description 获取全部净莲阁菜单，返回菜单项列表，每个菜单项包含名称、图片 url、营养价值表 json等信息。输入名称过滤菜单项，支持模糊匹配。
// @Tags  Menu
// @Accept json
// @Param query body model.MenuQueryRequest true "查询参数" schemaexample({"name":"紫菜汤","pageSize":10,"pageNumber":0})
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.Menu} "{"code":200,"msg":"ok","data":[{"id":1,"name":"紫菜汤","image_url":"/images/menu.jpg","desc":"美味的菜品","nutrition":"...","ingredients":"...","status":1}]}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/getMenu [post]
func GetMenu(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewMenuDB(db)

	// 使用 body 参数
	var queryReq model.MenuQueryRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	// 设置默认值
	if queryReq.PageSize <= 0 {
		queryReq.PageSize = setting.AppSetting.PageSize // Default page size
	}
	if queryReq.PageNumber < 0 {
		queryReq.PageNumber = 0
	}

	activityList, err := repo.GetMenuList(queryReq.PageSize, queryReq.PageNumber, queryReq.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activityList)
}
