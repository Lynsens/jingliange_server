package v1

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
)

// @Summary 提交建议箱内容
// @Description 用户提交对餐厅、菜单或小程序的建议，联系方式可选。
// @Tags Suggestion
// @Accept json
// @Param suggestion body model.SuggestionCreateRequest true "建议内容" schemaexample({"content":"希望增加更多清淡菜品","contact":"微信或手机号"})
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":"suggestion submitted successfully"}"
// @Failure 400 {object} app.Response "{"code":400,"msg":"invalid params","data":"Suggestion content is required"}"
// @Failure 500 {object} app.Response "{"code":500,"msg":"internal server error","data":null}"
// @Router /api/v1/suggestion/create [post]
func CreateSuggestion(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	aboutRepo := repo.NewAboutDb(db)

	var req model.SuggestionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Invalid input data")
		return
	}

	content := strings.TrimSpace(req.Content)
	contact := strings.TrimSpace(req.Contact)
	userNickname := strings.TrimSpace(req.UserNickname)
	userAvatarURL := strings.TrimSpace(req.UserAvatarURL)

	if content == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Suggestion content is required")
		return
	}
	if utf8.RuneCountInString(content) > 500 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Suggestion content is too long")
		return
	}
	if utf8.RuneCountInString(contact) > 64 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "Contact is too long")
		return
	}

	userID := ""
	if value, exists := c.Get("user_id"); exists {
		if id, ok := value.(string); ok {
			userID = id
		}
	}

	suggestion := model.Suggestion{
		UserID:        userID,
		UserNickname:  userNickname,
		UserAvatarURL: userAvatarURL,
		Content:       content,
		Contact:       contact,
	}

	if err := aboutRepo.CreateSuggestion(suggestion); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "suggestion submitted successfully")
}
