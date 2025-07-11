package v1

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/lynsens/jingliange_server/pkg/util"
)

// @Summary 获取功德榜列表
// @Description 获取捐款记录列表，支持年份、时间段、昵称筛选，支持按时间或金额排序，支持分页
// @Tags Donation
// @Accept json
// @Param query body model.DonationQueryRequest true "查询参数" schemaexample({"year":2025,"period":"all","donor_name":"","sort_by":"time","sort_order":"desc","page_size":10,"page_number":0})
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.Donation} "{"code":200,"msg":"ok","data":[{"id":1,"user_id":"user123","donor_name":"善心人士","amount":100.00,"donate_time":"2025-07-08T10:00:00Z","is_visible":1,"message":"祝愿净莲阁越来越好"}]}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/donation/getDonationList [post]
func GetDonationList(c *gin.Context) {
	appG := app.Gin{C: c}
	logging.Info("GetDonationList - 开始处理获取功德榜列表请求")

	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("GetDonationList - 数据库连接失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	donationRepo := repo.NewDonationDB(db)

	// 使用 body 参数
	var queryReq model.DonationQueryRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		logging.Error("GetDonationList - 参数绑定失败:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	logging.Info("GetDonationList - 请求参数:", fmt.Sprintf("year=%d, period=%s, donor_name=%s, sort_by=%s, sort_order=%s, page_size=%d, page_number=%d",
		queryReq.Year, queryReq.Period, queryReq.DonorName, queryReq.SortBy, queryReq.SortOrder, queryReq.PageSize, queryReq.PageNumber))

	// 设置默认值
	if queryReq.PageSize <= 0 {
		queryReq.PageSize = setting.AppSetting.PageSize
	}
	if queryReq.PageNumber < 0 {
		queryReq.PageNumber = 0
	}
	if queryReq.SortBy == "" {
		queryReq.SortBy = "time"
	}
	if queryReq.SortOrder == "" {
		queryReq.SortOrder = "desc"
	}
	if queryReq.Period == "" {
		queryReq.Period = "all"
	}
	// 如果没有指定年份，使用当前年份
	if queryReq.Year == 0 {
		queryReq.Year = time.Now().Year()
	}

	// 验证参数
	if queryReq.Period != "all" && queryReq.Period != "first" && queryReq.Period != "second" {
		logging.Error("GetDonationList - 参数验证失败: period无效")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid period, must be: all, first, or second")
		return
	}
	if queryReq.SortBy != "time" && queryReq.SortBy != "amount" {
		logging.Error("GetDonationList - 参数验证失败: sort_by无效")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid sort_by, must be: time or amount")
		return
	}
	if queryReq.SortOrder != "asc" && queryReq.SortOrder != "desc" {
		logging.Error("GetDonationList - 参数验证失败: sort_order无效")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid sort_order, must be: asc or desc")
		return
	}

	// 获取捐款列表
	logging.Info("GetDonationList - 开始查询捐款列表")
	donations, err := donationRepo.GetDonationList(queryReq)
	if err != nil {
		logging.Error("GetDonationList - 查询捐款列表失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	logging.Info("GetDonationList - 成功获取捐款列表, 数量:", len(donations))
	appG.Response(http.StatusOK, e.SUCCESS, donations)
}

// @Summary 创建捐款记录
// @Description 用户支付成功后创建捐款记录（需要JWT认证）
// @Tags Donation
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param donation body model.DonationCreateRequest true "捐款信息" schemaexample({"donor_name":"善心人士","amount":100.00,"message":"祝愿净莲阁越来越好"})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"Donation created successfully"}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 401 {object} app.Response "{"code":401,"msg":"unauthorized","data":"Token required"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/donation/createDonation [post]
func CreateDonation(c *gin.Context) {
	appG := app.Gin{C: c}
	logging.Info("CreateDonation - 开始处理创建捐款记录请求")

	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("CreateDonation - 数据库连接失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	donationRepo := repo.NewDonationDB(db)

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logging.Error("CreateDonation - JWT中未找到用户ID")
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "User ID not found in token")
		return
	}

	logging.Info("CreateDonation - 获取到用户ID:", userID)

	// 使用 body 参数
	var createReq model.DonationCreateRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		logging.Error("CreateDonation - 参数绑定失败:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	logging.Info("CreateDonation - 请求参数:", fmt.Sprintf("donor_name=%s, amount=%.2f, message=%s",
		createReq.DonorName, createReq.Amount, createReq.Message))

	// 验证必填参数
	if createReq.DonorName == "" {
		logging.Error("CreateDonation - 参数验证失败: donor_name为空")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Donor name is required")
		return
	}
	if createReq.Amount <= 0 {
		logging.Error("CreateDonation - 参数验证失败: amount必须大于0")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Amount must be greater than 0")
		return
	}
	if createReq.Amount > 10000 {
		logging.Error("CreateDonation - 参数验证失败: amount不能超过10000")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Amount must not exceed 10000")
		return
	}

	// 使用 utf8.RuneCountInString 来正确计算Unicode字符数
	donorNameLength := utf8.RuneCountInString(createReq.DonorName)
	if donorNameLength < 1 || donorNameLength > 10 {
		logging.Error("CreateDonation - 参数验证失败: donor_name长度不符合要求")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Donor name must be 1-10 characters")
		return
	}

	// 创建捐款记录
	donation := model.Donation{
		UserID:     userID.(string),
		DonorName:  createReq.DonorName,
		Amount:     createReq.Amount,
		DonateTime: time.Now(),
		IsVisible:  1, // 默认显示
		Message:    createReq.Message,
		Remarks:    "", // 管理员备注，初始为空
	}

	logging.Info("CreateDonation - 开始创建捐款记录")
	if err := donationRepo.CreateDonation(donation); err != nil {
		logging.Error("CreateDonation - 创建捐款记录失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	logging.Info("CreateDonation - 成功创建捐款记录, 用户:", userID, "金额:", createReq.Amount)
	appG.Response(http.StatusOK, e.SUCCESS, "Donation created successfully")
}

// @Summary 获取捐款统计
// @Description 获取指定时间段的捐款统计信息，包括总金额和总人次
// @Tags Donation
// @Accept json
// @Param stats body model.DonationStatsRequest true "统计参数" schemaexample({"year":2025,"period":"all"})
// @Produce  json
// @Success 200 {object} app.Response{data=model.DonationStats} "{"code":200,"msg":"ok","data":{"total_amount":5000.00,"total_count":50}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/donation/getDonationStats [post]
func GetDonationStats(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	donationRepo := repo.NewDonationDB(db)

	// 使用 body 参数
	var statsReq model.DonationStatsRequest
	if err := c.ShouldBindJSON(&statsReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	// 设置默认值
	if statsReq.Period == "" {
		statsReq.Period = "all"
	}
	if statsReq.Year == 0 {
		statsReq.Year = time.Now().Year()
	}

	// 验证参数
	if statsReq.Period != "all" && statsReq.Period != "first" && statsReq.Period != "second" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid period, must be: all, first, or second")
		return
	}

	// 获取统计信息
	stats, err := donationRepo.GetDonationStats(statsReq)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, stats)
}

// @Summary 用户认证
// @Description 简单的用户认证，基于微信小程序的用户ID，返回JWT token
// @Tags Auth
// @Accept json
// @Param auth body model.AuthRequest true "认证参数" schemaexample({"user_id":"user123"})
// @Produce  json
// @Success 200 {object} app.Response{data=map[string]interface{}} "{"code":200,"msg":"ok","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","user_id":"user123"}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/auth/login [post]
func AuthUser(c *gin.Context) {
	appG := app.Gin{C: c}
	logging.Info("AuthUser - 开始处理用户认证请求")

	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("AuthUser - 数据库连接失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	userRepo := repo.NewUserDB(db)

	// 使用 body 参数
	var authReq model.AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		logging.Error("AuthUser - 参数绑定失败:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	logging.Info("AuthUser - 认证用户:", authReq.UserID)

	// 验证参数
	if authReq.UserID == "" {
		logging.Error("AuthUser - 参数验证失败: user_id为空")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "User ID is required")
		return
	}

	// 创建或更新用户记录
	logging.Info("AuthUser - 开始创建或更新用户记录")
	if err := userRepo.CreateOrUpdateUser(authReq.UserID); err != nil {
		logging.Error("AuthUser - 创建或更新用户记录失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	// 生成JWT token
	logging.Info("AuthUser - 开始生成JWT token")
	token, err := util.GenerateToken(authReq.UserID)
	if err != nil {
		logging.Error("AuthUser - 生成JWT token失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	// 返回token和用户ID
	data := map[string]interface{}{
		"token":   token,
		"user_id": authReq.UserID,
	}

	logging.Info("AuthUser - 用户认证成功, 用户ID:", authReq.UserID)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 以下是管理员功能接口，暂时注释，后续实现
/*
// @Summary 更新捐款记录显示状态（管理员）
// @Description 管理员更新捐款记录的显示状态
// @Tags Admin
// @Accept json
// @Param update body map[string]interface{} true "更新参数" schemaexample({"id":1,"is_visible":0})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"Status updated successfully"}"
// @Router /api/admin/donation/updateVisibility [post]
func UpdateDonationVisibility(c *gin.Context) {
	// TODO: 实现管理员更新捐款记录显示状态
}

// @Summary 删除捐款记录（管理员）
// @Description 管理员删除捐款记录
// @Tags Admin
// @Accept json
// @Param delete body map[string]interface{} true "删除参数" schemaexample({"id":1})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"Donation deleted successfully"}"
// @Router /api/admin/donation/deleteDonation [post]
func DeleteDonation(c *gin.Context) {
	// TODO: 实现管理员删除捐款记录
}
*/
