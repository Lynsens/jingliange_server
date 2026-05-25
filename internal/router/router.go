package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/router/api"
	"github.com/lynsens/jingliange_server/internal/router/api/admin"
	v1 "github.com/lynsens/jingliange_server/internal/router/api/v1"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/upload"
	"github.com/lynsens/jingliange_server/pkg/util"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/lynsens/jingliange_server/docs"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.LoggerWithWriter(logging.F))
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFS("/uploads/images", http.Dir(upload.GetImageFullPath()))

	// r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	// r.POST("/auth", api.GetAuth)
	r.POST("/uploadImage", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	{
		// 公开接口（不需要认证）
		//获取净莲阁介绍
		apiv1.GET("/about/getDescription", v1.GetDescription)
		//获取首页头图
		apiv1.GET("/about/getTopImage", v1.GetTopImage)
		//获取活动列表
		apiv1.POST("/about/getActivityList", v1.GetActivityList)
		//获取图片列表
		apiv1.POST("/about/getImageList", v1.GetImangeList)
		//获取单个菜品信息
		apiv1.POST("/menu/getMenuByID", v1.GetMenuByID)
		//获取功德榜列表
		apiv1.POST("/donation/getDonationList", v1.GetDonationList)
		//获取捐款统计
		apiv1.POST("/donation/getDonationStats", v1.GetDonationStats)
		//用户认证
		apiv1.POST("/auth/login", v1.AuthUser)
		//上传图片
		apiv1.POST("/uploadImage", api.UploadImage)
	}

	// 支持可选认证的接口（有token显示个人信息，无token显示公开信息）
	apiOptional := r.Group("/api/v1")
	apiOptional.Use(util.OptionalJWT())
	{
		//获取菜单（支持显示用户点赞状态）
		apiOptional.POST("/menu/getMenu", v1.GetMenu)
		//获取菜品评论列表（支持当前用户评论置顶）
		apiOptional.POST("/menu/getComments", v1.GetMenuComments)
		//获取当前套餐推荐
		apiOptional.GET("/combo/active", v1.GetActiveComboRecommendation)
		//提交建议箱内容
		apiOptional.POST("/suggestion/create", v1.CreateSuggestion)
	}

	// 需要认证的接口
	apiAuth := r.Group("/api/v1")
	apiAuth.Use(util.JWT())
	{
		//菜品点赞
		apiAuth.POST("/menu/like", v1.LikeMenu)
		//获取菜品点赞状态
		apiAuth.POST("/menu/getLikeStatus", v1.GetMenuLikeStatus)
		//菜品评论
		apiAuth.POST("/menu/comment", v1.CommentMenu)
		//删除自己的菜品评论
		apiAuth.DELETE("/menu/comment/delete", v1.DeleteMenuComment)
		//创建捐款记录
		apiAuth.POST("/donation/createDonation", v1.CreateDonation)
	}

	apiAdmin := r.Group("/api/admin")
	apiAdmin.POST("/login", admin.Login)
	apiAdmin.Use(util.AdminJWT())
	{
		// 管理员菜单列表
		apiAdmin.POST("/menu/list", admin.GetMenuItems)
		// 管理员评论列表
		apiAdmin.POST("/comment/list", admin.GetComments)
		// 管理员删除评论
		apiAdmin.DELETE("/comment/delete", admin.DeleteComment)
		// 上传菜品
		apiAdmin.POST("/uploadMenuItem", admin.UploadMenuItem)
		// 更新菜品
		apiAdmin.PUT("/updateMenuItem", admin.UpdateMenuItem)
		// 设置今日推荐菜品
		apiAdmin.PUT("/recommendMenuItem", admin.SetRecommendedMenuItem)
		// 下架或重新上架菜品
		apiAdmin.PUT("/archiveMenuItem", admin.ArchiveMenuItem)
		// 删除菜品
		apiAdmin.DELETE("/deleteMenuItem", admin.DeleteMenuItem)
		// 活动管理
		apiAdmin.POST("/activity/list", admin.GetActivityList)
		apiAdmin.POST("/activity/create", admin.CreateActivity)
		apiAdmin.PUT("/activity/update", admin.UpdateActivity)
		apiAdmin.DELETE("/activity/delete", admin.DeleteActivity)
		apiAdmin.PUT("/activity/top", admin.SetActivityTop)
		// 套餐推荐管理
		apiAdmin.POST("/combo/list", admin.GetComboRecommendations)
		apiAdmin.POST("/combo/create", admin.CreateComboRecommendation)
		apiAdmin.PUT("/combo/update", admin.UpdateComboRecommendation)
		apiAdmin.PUT("/combo/active", admin.SetComboRecommendationActive)
		apiAdmin.DELETE("/combo/delete", admin.DeleteComboRecommendation)
		// 建议箱管理
		apiAdmin.POST("/suggestion/list", admin.GetSuggestions)
		apiAdmin.PUT("/suggestion/status", admin.UpdateSuggestionStatus)
		// 只读运维维护面板
		apiAdmin.GET("/ops/summary", admin.GetOpsSummary)
		apiAdmin.GET("/ops/access-logs", admin.GetOpsAccessLogs)
		apiAdmin.GET("/ops/error-logs", admin.GetOpsErrorLogs)
		apiAdmin.GET("/ops/app-logs", admin.GetOpsAppLogs)
	}

	return r
}
