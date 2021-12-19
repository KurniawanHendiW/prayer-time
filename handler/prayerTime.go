package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prayer-time/service/prayerTime"
)

type Handler interface {
	GetKeyPrayerTime(c *gin.Context)
	GetDataPrayerTime(c *gin.Context)
	GetCityByName(c *gin.Context)
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
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	resp, err := h.prayerTime.GetKeyPrayerTime(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) GetDataPrayerTime(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "param key is required",
		})
		return
	}

	dataPrayTimeRequest := prayerTime.DataPrayerTimeRequest{
		Key: key,
	}

	resp, err := h.prayerTime.GetDataPrayerTime(dataPrayTimeRequest)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	downloadName := "Jakarta_Indonesia.ics"

	c.Writer.Header().Set("Content-type", "text/calendar;charset=UTF-8")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
	c.Writer.Write([]byte(resp.Data))
}

func (h *handler) GetCityByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "city name is required",
		})
		return
	}

	resp, err := h.prayerTime.GetCityByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
