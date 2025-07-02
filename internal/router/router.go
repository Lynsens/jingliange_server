package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/router/api"
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
		apiv1.GET("/getDescription", v1.GetDescription)
		//获取活动列表
		apiv1.POST("/getActivityList", v1.GetActivityList)
		//获取图片列表
		apiv1.POST("/getImageList", v1.GetImangeList)
		// //新建标签
		// apiv1.POST("/tags", v1.AddTag)
		// //更新指定标签
		// apiv1.PUT("/tags/:id", v1.EditTag)
		// //删除指定标签
		// apiv1.DELETE("/tags/:id", v1.DeleteTag)
		// //导出标签
		// r.POST("/tags/export", v1.ExportTag)
		// //导入标签
		// r.POST("/tags/import", v1.ImportTag)

		// //获取文章列表
		// apiv1.GET("/articles", v1.GetArticles)
		// //获取指定文章
		// apiv1.GET("/articles/:id", v1.GetArticle)
		// //新建文章
		// apiv1.POST("/articles", v1.AddArticle)
		// //更新指定文章
		// apiv1.PUT("/articles/:id", v1.EditArticle)
		// //删除指定文章
		// apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		// //生成文章海报
		// apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
