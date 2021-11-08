package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prayer-time/service/prayerTime"
)

type Handler interface {
	GetKeyPrayerTime(c *gin.Context)
	GetDataPrayerTime(c *gin.Context)
}

type handler struct {
	prayerTime prayerTime.Service
}

func NewHandler(prayerTime prayerTime.Service) Handler {
	return &handler{
		prayerTime: prayerTime,
	}
}

func (h *handler) GetKeyPrayerTime(c *gin.Context) {
	var req prayerTime.KeyPrayerTimeRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.prayerTime.GetKeyPrayerTime(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) GetDataPrayerTime(c *gin.Context) {
	//var dataPrayTimeRequest prayerTime.DataPrayerTimeRequest
	//if err := c.Bind(&dataPrayTimeRequest); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//	return
	//}

	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, "param key is required")
		return
	}

	dataPrayTimeRequest := prayerTime.DataPrayerTimeRequest{
		Key: key,
	}

	resp, err := h.prayerTime.GetDataPrayerTime(dataPrayTimeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	downloadName := "Jakarta_Indonesia.ics"

	c.Writer.Header().Set("Content-type", "text/calendar;charset=UTF-8")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
	c.Writer.Write([]byte(resp.Data))
}
