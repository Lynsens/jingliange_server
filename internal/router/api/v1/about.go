package v1

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/e"
)

// @Summary Retrieve the description of Jingliange
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/getDescription [get]
func GetDescription(c *gin.Context) {

	appG := app.Gin{C: c}

	// to implement: db model and fetch from db
	rsp := "净莲阁成立于，是非营利性一家素食餐厅。"

	appG.Response(http.StatusOK, e.SUCCESS, rsp)
}
