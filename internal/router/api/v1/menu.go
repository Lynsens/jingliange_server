package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"gorm.io/gorm"
)

// @Summary 获取净莲阁的菜单
// @Description 获取全部净莲阁菜单，返回菜单项列表，每个菜单项包含名称、图片 url、营养价值表 json、点赞数等信息。输入名称过滤菜单项，支持模糊匹配。如果提供了JWT token，会返回用户的点赞状态。
// @Tags  Menu
// @Accept json
// @Param Authorization header string false "Bearer token (可选)"
// @Param query body model.MenuQueryRequest true "查询参数" schemaexample({"name":"紫菜汤","page_size":10,"page_number":0})
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.MenuWithUserLikes} "{"code":200,"msg":"ok","data":[{"id":1,"name":"紫菜汤","image_url":"/images/menu.jpg","desc":"美味的菜品","nutrition":"...","ingredients":"...","status":1,"like_count":5,"liked":true}]}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/getMenu [post]
func GetMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	logging.Info("GetMenu - 开始处理获取菜单列表请求")

	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("GetMenu - 数据库连接失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	repo := repo.NewMenuDB(db)

	// 使用 body 参数
	var queryReq model.MenuQueryRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		logging.Error("GetMenu - 参数绑定失败:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	logging.Info("GetMenu - 请求参数:", fmt.Sprintf("name=%s, page_size=%d, page_number=%d",
		queryReq.Name, queryReq.PageSize, queryReq.PageNumber))

	// 设置默认值
	if queryReq.PageSize <= 0 {
		queryReq.PageSize = setting.AppSetting.PageSize // Default page size
	}
	if queryReq.PageNumber < 0 {
		queryReq.PageNumber = 0
	}

	// 检查是否有用户认证信息
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		userIDStr = userID.(string)
		logging.Info("GetMenu - 检测到已认证用户:", userIDStr)
	} else {
		logging.Info("GetMenu - 未检测到用户认证，返回公开信息")
	}

	// 使用包含用户点赞状态的方法获取菜单列表
	menuList, err := repo.GetMenuListWithUserLikes(queryReq.PageSize, queryReq.PageNumber, queryReq.Name, userIDStr)
	if err != nil {
		logging.Error("GetMenu - 查询菜单列表失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	logging.Info("GetMenu - 成功获取菜单列表, 数量:", len(menuList))
	appG.Response(http.StatusOK, e.SUCCESS, menuList)
}

// @Summary 菜品点赞
// @Description 用户为喜欢的菜品点赞或取消点赞（需要JWT认证）
// @Tags Menu
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param like body model.MenuLikeRequest true "点赞参数" schemaexample({"menu_id":1})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"liked successfully"}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 401 {object} app.Response "{"code":401,"msg":"unauthorized","data":"Token required"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/like [post]
func LikeMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	logging.Info("LikeMenu - 开始处理菜品点赞请求")

	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("LikeMenu - 数据库连接失败:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logging.Error("LikeMenu - JWT中未找到用户ID")
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "User ID not found in token")
		return
	}

	logging.Info("LikeMenu - 获取到用户ID:", userID)

	// 使用 body 参数
	var likeReq model.MenuLikeRequest
	if err := c.ShouldBindJSON(&likeReq); err != nil {
		logging.Error("LikeMenu - 参数绑定失败:", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	logging.Info("LikeMenu - 请求参数: menu_id=", likeReq.MenuID)

	// 验证参数
	if likeReq.MenuID <= 0 {
		logging.Error("LikeMenu - 参数验证失败: menu_id无效")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	// 检查菜品是否存在且状态正常
	existingMenu, err := menuRepo.GetMenuByID(likeReq.MenuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Error("LikeMenu - 菜品不存在:", likeReq.MenuID)
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	// 检查菜品状态
	if existingMenu.Status == 0 || existingMenu.IsArchived == 1 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item not available")
		return
	}

	// 检查当前点赞状态
	currentStatus, err := menuRepo.GetMenuLikeStatus(likeReq.MenuID, userID.(string))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	// 如果当前是喜欢状态，则取消点赞；否则点赞
	if currentStatus == 1 {
		// 取消点赞
		if err := menuRepo.UnlikeMenu(likeReq.MenuID, userID.(string)); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, "unliked successfully")
	} else {
		// 点赞
		if err := menuRepo.LikeMenu(likeReq.MenuID, userID.(string)); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, "liked successfully")
	}
}

// @Summary 获取菜品点赞状态
// @Description 获取用户对特定菜品的点赞状态（需要JWT认证）
// @Tags Menu
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param status body model.MenuLikeRequest true "查询参数" schemaexample({"menu_id":1})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":{"liked":true,"preference":1}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 401 {object} app.Response "{"code":401,"msg":"unauthorized","data":"Token required"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/getLikeStatus [post]
func GetMenuLikeStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "User ID not found in token")
		return
	}

	// 使用 body 参数
	var statusReq model.MenuLikeRequest
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	// 验证参数
	if statusReq.MenuID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	// 检查菜品是否存在且可见
	existingMenu, err := menuRepo.GetMenuByID(statusReq.MenuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}
	if existingMenu.Status == 0 || existingMenu.IsArchived == 1 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item not available")
		return
	}

	// 获取点赞状态
	preference, err := menuRepo.GetMenuLikeStatus(statusReq.MenuID, userID.(string))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	// 构建响应数据
	response := map[string]interface{}{
		"liked":      preference == 1,
		"preference": preference,
	}

	appG.Response(http.StatusOK, e.SUCCESS, response)
}

// @Summary 菜品评论
// @Description 用户对菜品进行评论（需要JWT认证）
// @Tags Menu
// @Accept json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param comment body model.MenuCommentRequest true "评论参数" schemaexample({"menu_id":1,"comment":"非常好吃的菜品！"})
// @Produce  json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"comment added successfully"}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 401 {object} app.Response "{"code":401,"msg":"unauthorized","data":"Token required"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/comment [post]
func CommentMenu(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, "User ID not found in token")
		return
	}

	// 使用 body 参数
	var commentReq model.MenuCommentRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	// 验证参数
	if commentReq.MenuID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}
	if commentReq.Comment == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Comment content is required")
		return
	}
	commentReq.Comment = strings.TrimSpace(commentReq.Comment)
	commentReq.UserNickname = strings.TrimSpace(commentReq.UserNickname)
	commentReq.UserAvatarURL = strings.TrimSpace(commentReq.UserAvatarURL)
	if commentReq.Comment == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Comment content is required")
		return
	}
	if commentReq.UserNickname == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "User nickname is required")
		return
	}
	if commentReq.UserAvatarURL == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "User avatar is required")
		return
	}
	if utf8.RuneCountInString(commentReq.UserNickname) > 64 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "User nickname must not exceed 64 characters")
		return
	}
	if len(commentReq.UserAvatarURL) > 256 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "User avatar URL must not exceed 256 characters")
		return
	}

	// 检查菜品是否存在且状态正常
	existingMenu, err := menuRepo.GetMenuByID(commentReq.MenuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	// 检查菜品状态
	if existingMenu.Status == 0 || existingMenu.IsArchived == 1 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item not available")
		return
	}

	// 添加评论
	if err := menuRepo.CommentMenu(commentReq.MenuID, userID.(string), commentReq.Comment, commentReq.UserNickname, commentReq.UserAvatarURL); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "comment added successfully")
}

// @Summary 获取菜品评论列表
// @Description 获取指定菜品的评论列表，支持分页
// @Tags Menu
// @Accept json
// @Param query body model.MenuCommentsQueryRequest true "查询参数" schemaexample({"menu_id":1,"page_size":10,"page_number":0})
// @Produce  json
// @Success 200 {object} app.Response{data=[]model.MenuFeedback} "{"code":200,"msg":"ok","data":[{"id":1,"menu_id":1,"user_id":"user123","preference":1,"comment":"非常好吃！","status":1}]}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/getComments [post]
func GetMenuComments(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	// 使用 body 参数
	var queryReq model.MenuCommentsQueryRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	// 验证参数
	if queryReq.MenuID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	// 设置默认值
	if queryReq.PageSize <= 0 {
		queryReq.PageSize = setting.AppSetting.PageSize
	}
	if queryReq.PageNumber < 0 {
		queryReq.PageNumber = 0
	}

	// 检查菜品是否存在且可见
	existingMenu, err := menuRepo.GetMenuByID(queryReq.MenuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}
	if existingMenu.Status == 0 || existingMenu.IsArchived == 1 {
		appG.Response(http.StatusNotFound, e.ERROR, "Menu item not available")
		return
	}

	// 获取评论列表
	comments, err := menuRepo.GetMenuComments(queryReq.MenuID, queryReq.PageSize, queryReq.PageNumber)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, comments)
}

// @Summary 获取单个菜品信息
// @Description 根据菜品ID获取单个菜品的详细信息，包括点赞数
// @Tags Menu
// @Accept json
// @Param query body model.MenuByIDRequest true "查询参数" schemaexample({"menu_id":1})
// @Produce  json
// @Success 200 {object} app.Response{data=model.MenuWithLikes} "{"code":200,"msg":"ok","data":{"id":1,"name":"紫菜汤","image_url":"/images/menu.jpg","desc":"美味的菜品","nutrition":"...","ingredients":"...","status":1,"like_count":5,"create_time":"2012-1-1","update_time":"2012-1-1"}}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Invalid input data"}"
// @Failure 404 {object} app.Response "{"code":404,"msg":"not found","data":"Menu item not found"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/menu/getMenuByID [post]
func GetMenuByID(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	menuRepo := repo.NewMenuDB(db)

	// 使用 body 参数
	var queryReq model.MenuByIDRequest
	if err := c.ShouldBindJSON(&queryReq); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("Invalid input data: %v", err))
		return
	}

	// 验证参数
	if queryReq.MenuID <= 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid menu ID")
		return
	}

	// 获取菜品信息
	menuWithLikes, err := menuRepo.GetMenuByIDWithLikes(queryReq.MenuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusNotFound, e.ERROR, "Menu item not found")
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, menuWithLikes)
}
