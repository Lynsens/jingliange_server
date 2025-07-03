package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
)

// @Summary 获取净莲阁的菜单
// @Description 获取全部净莲阁菜单，返回菜单项列表，每个菜单项包含名称、图片 url、营养价值表 json等信息。输入名称过滤菜单项，支持模糊匹配。
// @Tags  Menu
// @Param name query string false "菜单名称，支持模糊匹配" default("")
// @Param pageSize query string false "每页数量" default(10)
// @Param pageNumber query string false "页码, 从零开始" default(0)
// @Produce  json
// @Success 200 {object} app.Response{data=model.Menu}
// @Failure 500 {object} app.Response
// @Router /api/v1/menu/getMenu [post]
func GetMenu(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
	}

	repo := repo.NewMenuDB(db)

	name := c.DefaultQuery("name", "")

	pageNumberStr := c.DefaultQuery("pageNumber", "0")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid pageNumber")
		return
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid pageSize")
		return
	}
	if pageSize <= 0 {
		pageSize = setting.AppSetting.PageSize // Default page size
	}

	activityList, err := repo.GetMenuList(pageSize, pageNumber, name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activityList)
}
