package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/upload"
)

// @Summary 上传图片
// @Description 上传净莲阁的图片。支持图片格式检查和大小限制，返回图片的访问地址、保存路径和文件名。
// @Param image formData file true "图片文件"
// @Param desc formData string false "图片描述"
// @Param top_pic formData int true "是否为头图：0 否，1 是"
// @Param type formData int true "图片类型：0 活动，1 餐厅介绍"
// @Produce json
// @Success 200 {object} app.Response{data=map[string]string} "成功"
// @Failure 400 {object} app.Response "参数错误"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /api/v1/uploadImage [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Error("UploadImage: FormFile error", err, nil)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid image file")
		return
	}

	if image == nil {
		logging.Warn("UploadImage: No image provided")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "no image provided")
		return
	}

	// Parse image context
	desc, topPic, imageType := ParseImgContext(c)
	if desc == "" {
		logging.Warn("UploadImage: Image description is required")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "image description is required")
		return
	}
	if topPic < 0 || topPic > 1 {
		logging.Warn("UploadImage: Invalid top_pic value", topPic)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid top_pic value")
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Error("UploadImage: CheckImage error", err, nil)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Error("UploadImage: SaveUploadedFile error", err, nil)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	logging.Info("UploadImage: Image uploaded successfully", imageName, nil)
	imageUrl := upload.GetImageFullUrl(imageName)

	// Save image URL to database
	err = SaveImageToDB(imageUrl, desc, 1, topPic, imageType)
	if err != nil {
		logging.Error("UploadImage: SaveImageToDB error", err, nil)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_url":  imageUrl,
		"save_path":  savePath,
		"image_name": imageName,
	})
	return
}

// ParseImgContext parses the image upload context and returns the necessary parameters.
// It validates the input and returns appropriate error responses if any required fields are missing or invalid.
func ParseImgContext(c *gin.Context) (desc string, topPic int, imageType int) {
	appG := app.Gin{C: c}

	desc = c.PostForm("desc")
	if desc == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "image description is required")
		return
	}

	topPicStr := c.PostForm("top_pic")
	if topPicStr == "" {
		topPic = 0
	} else {
		var convErr error
		topPic, convErr = strconv.Atoi(topPicStr)
		if convErr != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid top_pic value")
			return
		}
	}

	imageTypeStr := c.PostForm("type")
	if imageTypeStr == "" {
		imageType = 0
	} else {
		var convErr error
		imageType, convErr = strconv.Atoi(imageTypeStr)
		if convErr != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "invalid type value")
			return
		}
	}

	return desc, topPic, imageType
}

// save img url to database
func SaveImageToDB(url string, desc string, status int, topPic int, imageType int) error {
	db, err := repo.ConnectDb()
	if err != nil {
		logging.Error("SaveImageToDB: ConnectDb error", err)
		return err
	}

	repo := repo.NewAboutDb(db)
	err = repo.SaveImageUrlToDB(url, desc, uint(status), topPic, uint(imageType))
	if err != nil {
		return err
	}
	logging.Info("SaveImageToDB: Image URL saved successfully", url)
	return nil
}
