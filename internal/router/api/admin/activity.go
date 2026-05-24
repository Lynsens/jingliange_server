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
	"gorm.io/gorm"
)

type activityRepository interface {
	GetAdminActivityList(keyword string, pageSize int, pageNumber int) ([]model.Activity, error)
	GetActivityByID(id uint64) (model.Activity, error)
	CreateActivity(activity model.Activity) error
	UpdateActivity(activity model.Activity) error
	DeleteActivity(id uint64) error
	SetActivityTop(id uint64, isTop int) error
}

func normalizeTopValue(value int) int {
	if value == 1 {
		return 1
	}
	return 0
}

func validateActivityPayload(req model.AdminActivityRequest) string {
	if strings.TrimSpace(req.Title) == "" {
		return "Activity title is required"
	}
	if strings.TrimSpace(req.EventTime) == "" {
		return "Activity time is required"
	}
	if strings.TrimSpace(req.Place) == "" {
		return "Activity place is required"
	}
	if strings.TrimSpace(req.Content) == "" {
		return "Activity content is required"
	}
	return ""
}

func getActivityRepo(c *gin.Context, appG app.Gin) (activityRepository, bool) {
	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return nil, false
	}

	return repo.NewAboutDb(db), true
}

// GetActivityList 获取管理员活动列表
func GetActivityList(c *gin.Context) {
	appG := app.Gin{C: c}
	activityRepo, ok := getActivityRepo(c, appG)
	if !ok {
		return
	}

	var req model.AdminActivityListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	activities, err := activityRepo.GetAdminActivityList(strings.TrimSpace(req.Keyword), req.PageSize, req.PageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activities)
}

// CreateActivity 新增活动
func CreateActivity(c *gin.Context) {
	appG := app.Gin{C: c}
	activityRepo, ok := getActivityRepo(c, appG)
	if !ok {
		return
	}

	var req model.AdminActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if msg := validateActivityPayload(req); msg != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, msg)
		return
	}

	activity := model.Activity{
		Title:     strings.TrimSpace(req.Title),
		Content:   strings.TrimSpace(req.Content),
		EventTime: strings.TrimSpace(req.EventTime),
		Place:     strings.TrimSpace(req.Place),
		IsTop:     normalizeTopValue(req.IsTop),
		Status:    1,
	}

	if err := activityRepo.CreateActivity(activity); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

// UpdateActivity 更新活动
func UpdateActivity(c *gin.Context) {
	appG := app.Gin{C: c}
	activityRepo, ok := getActivityRepo(c, appG)
	if !ok {
		return
	}

	var req model.AdminActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid activity ID")
		return
	}

	if msg := validateActivityPayload(req); msg != "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, msg)
		return
	}

	existingActivity, err := activityRepo.GetActivityByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Activity not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}
	if existingActivity.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Activity already deleted")
		return
	}

	activity := model.Activity{
		ID:        req.ID,
		Title:     strings.TrimSpace(req.Title),
		Content:   strings.TrimSpace(req.Content),
		EventTime: strings.TrimSpace(req.EventTime),
		Place:     strings.TrimSpace(req.Place),
		IsTop:     normalizeTopValue(req.IsTop),
	}

	if err := activityRepo.UpdateActivity(activity); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

// DeleteActivity 删除活动
func DeleteActivity(c *gin.Context) {
	appG := app.Gin{C: c}
	activityRepo, ok := getActivityRepo(c, appG)
	if !ok {
		return
	}

	var req model.AdminActivityIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid activity ID")
		return
	}

	existingActivity, err := activityRepo.GetActivityByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Activity not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}
	if existingActivity.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Activity already deleted")
		return
	}

	if err := activityRepo.DeleteActivity(req.ID); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// SetActivityTop 设置活动置顶
func SetActivityTop(c *gin.Context) {
	appG := app.Gin{C: c}
	activityRepo, ok := getActivityRepo(c, appG)
	if !ok {
		return
	}

	var req model.AdminActivityTopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	if req.ID == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid activity ID")
		return
	}

	isTop := normalizeTopValue(req.IsTop)
	existingActivity, err := activityRepo.GetActivityByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Activity not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}
	if existingActivity.Status == 0 {
		appG.Response(http.StatusNotFound, e.ERROR, "Activity already deleted")
		return
	}

	if err := activityRepo.SetActivityTop(req.ID, isTop); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
