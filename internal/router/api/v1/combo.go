package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/repo"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"gorm.io/gorm"
)

// @Summary 获取当前套餐推荐
// @Description 获取首页当前启用的套餐推荐，包含套餐内可见菜品列表。
// @Tags Combo
// @Accept json
// @Param Authorization header string false "Bearer token (可选)"
// @Produce json
// @Success 200 {object} app.Response "{"code":200,"msg":"ok","data":{"id":1,"title":"清淡养胃套餐","items":[]}}"
// @Router /api/v1/combo/active [get]
func GetActiveComboRecommendation(c *gin.Context) {
	appG := app.Gin{C: c}

	db, err := repo.ConnectDb()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DB, nil)
		return
	}

	userID := ""
	if value, exists := c.Get("user_id"); exists {
		if id, ok := value.(string); ok {
			userID = id
		}
	}

	menuRepo := repo.NewMenuDB(db)
	combo, err := menuRepo.GetActiveComboRecommendation(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, combo)
}
