package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/router/api"
	"github.com/lynsens/jingliange_server/internal/router/api/admin"
	v1 "github.com/lynsens/jingliange_server/internal/router/api/v1"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/upload"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/lynsens/jingliange_server/docs"
	// "github.com/EDDYCJY/go-gin-example/middleware/jwt"
	// "github.com/EDDYCJY/go-gin-example/pkg/export"
	// "github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	// "github.com/EDDYCJY/go-gin-example/pkg/upload"
	// "github.com/EDDYCJY/go-gin-example/routers/api"
	// "github.com/EDDYCJY/go-gin-example/routers/api/v1"
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
	// apiv1.Use(jwt.JWT())
	{
		//获取净莲阁介绍
		apiv1.GET("/about/getDescription", v1.GetDescription)
		//获取首页头图
		apiv1.GET("/about/getTopImage", v1.GetTopImage)
		//获取活动列表
		apiv1.POST("/about/getActivityList", v1.GetActivityList)
		//获取图片列表
		apiv1.POST("/about/getImageList", v1.GetImangeList)
		//获取菜单
		apiv1.POST("/menu/getMenu", v1.GetMenu)
		//获取单个菜品信息
		apiv1.POST("/menu/getMenuByID", v1.GetMenuByID)
		//菜品点赞
		apiv1.POST("/menu/like", v1.LikeMenu)
		//获取菜品点赞状态
		apiv1.POST("/menu/getLikeStatus", v1.GetMenuLikeStatus)
		//菜品评论
		apiv1.POST("/menu/comment", v1.CommentMenu)
		//获取菜品评论列表
		apiv1.POST("/menu/getComments", v1.GetMenuComments)
		// //更新指定文章
		// apiv1.PUT("/articles/:id", v1.EditArticle)
		// //删除指定文章
		// apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	apiAdmin := r.Group("/api/admin")
	// apiAdmin.Use(jwt.JWT())
	{
		// 上传菜品
		apiAdmin.POST("/uploadMenuItem", admin.UploadMenuItem)
		// 删除菜品
		apiAdmin.DELETE("/deleteMenuItem", admin.DeleteMenuItem)
	}

	return r
}
