package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"

	_ "github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/e"
)

// @Summary 获取净莲阁介绍
// @Produce  json
// @Tags About
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/about/getDescription [get]
func GetDescription(c *gin.Context) {
	appG := app.Gin{C: c}

	// db, err := d.ConnectDb()
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
	// }
	// repo := d.NewAboutDb(db)

	// desc, err := repo.GetDescription()
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	// 	return
	// }

	desc := "净莲阁成立于 2018 年，是一家非营利性素食餐厅。"
	appG.Response(http.StatusOK, e.SUCCESS, desc)
}

// @Summary 获取头图
// @Description 获取净莲阁的头图。返回图片地址、简介、是否为头图等信息。
// @Tags About
// @Produce  json
// @Success 200 {object} app.Response{data=model.Image} "成功返回头图信息"
// @Failure 500 {object} app.Response
// @Router /api/v1/about/getTopImage [get]
func GetTopImage(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewAboutDb(db)
	image, err := repo.GetTopImage()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, image)
}

// @Summary 获取活动列表
// @Description 获取净莲阁的活动列表，输入时间戳
// @Tags About
// @Param timestamp query string true "时间戳" default(0)
// @Param pageNumber query string false "页码, 从零开始" default(0)
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/about/getActivityList [post]
func GetActivityList(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
	}

	repo := repo.NewAboutDb(db)

	timestampStr := c.DefaultQuery("timestamp", "0")
	pageNumberStr := c.DefaultQuery("pageNumber", "0")

	timestamp, err := strconv.Atoi(timestampStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid timestamp")
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid pageNumber")
		return
	}

	activityList, err := repo.GetActivityList(timestamp, pageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activityList)
}
