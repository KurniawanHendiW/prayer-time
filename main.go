package main

import (
	"net/http"

	"github.com/prayer-time/client/waktusholat"
	"github.com/prayer-time/config"
	"github.com/prayer-time/handler"
	"github.com/prayer-time/service/prayerTime"
	"github.com/prayer-time/util"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Get()

	router := initRouter(cfg)

	util.RunServerGracefully(cfg.RestPort, cfg.TimeOut, router)
}

func initRouter(cfg config.Config) *gin.Engine {
	// init service
	waktuSholatSvc := waktusholat.NewService(cfg.WaktuSholatHost)
	prayerTimeSvc := prayerTime.NewService(waktuSholatSvc, cfg.PassKey)

	// init handler
	prayerTimeHandler := handler.NewHandler(prayerTimeSvc)

	router := gin.Default()

	router.LoadHTMLGlob("views/pages/*")
	router.Static("/css", "./views/assets/css")
	router.Static("/js", "./views/assets/js")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Prayer Time For Calendar",
		})
	})

	route := router.Group("/prayer-time")
	{
		route.POST("/get-key", prayerTimeHandler.GetKeyPrayerTime)
		route.GET("/get", prayerTimeHandler.GetDataPrayerTime)
	}

	return router
}
