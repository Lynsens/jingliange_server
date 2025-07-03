package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"

	"fmt"

	routers "github.com/lynsens/jingliange_server/internal/router"
)

func init() {
	// Load the configuration file
	setting.Setup()
	logging.Setup()

}

// @title Jingliange Server API
// @version 1.0
// @description This is the backend API for Jingliange.
// @termsOfService TO IMPLEMENT

// @host https://jingliange.com
// @BasePath /api/v1
// @schemes http
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)
	defer logging.Close()
	server.ListenAndServe()

}
