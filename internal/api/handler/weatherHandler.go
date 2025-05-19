package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"weather/internal/domain"
)

type WeatherHandler struct {
	service domain.WeatherService
}

func NewWeatherHandler(service domain.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		service: service,
	}
}

func (w WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}
	currentWeather, err := w.service.GetWeather(city)
	if err != nil {
		fmt.Println("Error fetching weather:", err)
		c.JSON(http.StatusNotFound, "City not found")
		return
	}

	c.IndentedJSON(http.StatusOK, currentWeather)
}
