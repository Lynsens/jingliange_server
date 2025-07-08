package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"

	"github.com/lynsens/jingliange_server/internal/model"
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
// @Description 获取净莲阁的活动列表，支持时间戳过滤和分页
// @Tags About
// @Accept json
// @Param query body model.ActivityQueryRequest true "查询参数" schemaexample({"timestamp":0,"page_number":0})
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.Activity} "{"code":200,"msg":"ok","data":[{"id":1,"title":"素食文化活动","content":"介绍素食文化","img":"/images/activity1.jpg","status":1}]}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/about/getActivityList [post]
func GetActivityList(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewAboutDb(db)

	// 使用 body 参数
	var queryReq model.ActivityQueryRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	// 设置默认值
	if queryReq.PageNumber < 0 {
		queryReq.PageNumber = 0
	}

	activityList, err := repo.GetActivityList(queryReq.Timestamp, queryReq.PageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activityList)
}
