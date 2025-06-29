package v1

import (
	"net/http"
	"strconv"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/gin-gonic/gin"

	_ "github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/e"
)

// @Summary 获取净莲阁介绍
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/getDescription [get]
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

// @Summary 获取图片列表
// @Description 获取净莲阁的图片列表。返回图片地址、简介、是否为头图等信息。类型为 0 表示活动图片，1 表示餐厅介绍图片。
// @Param type query int false "图片类型：0 活动，1 餐厅介绍" default(0)
// @Param top_pic query int false "是否为头图：0 否，1 是" default(0)
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.Image} "成功返回图片列表"
// @Failure 500 {object} app.Response
// @Router /api/v1/getImageList [post]
func GetImangeList(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
	}

	repo := repo.NewAboutDb(db)
	typeStr := c.DefaultQuery("type", "0")
	imageType, err := strconv.Atoi(typeStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid type")
		return
	}
	topPicStr := c.DefaultQuery("top_pic", "0")
	topPic, err := strconv.Atoi(topPicStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid top_pic")
		return
	}
	imageList, err := repo.GetImageList(imageType, topPic)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, imageList)

}

// @Summary 获取活动列表
// @Description 获取净莲阁的活动列表，输入时间戳
// @Param timestamp query int true "时间戳" default(0)
// @Param pageNumber query int false "页码, 从零开始" default(0)
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/getActivityList [post]
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
