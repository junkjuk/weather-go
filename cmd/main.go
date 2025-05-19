package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"weather/configs"
	"weather/internal/api/handler"
	"weather/internal/domain"

	"github.com/caarlos0/env/v11"
	"github.com/go-co-op/gocron/v2"
	weatherClient "github.com/kashifkhan0771/go-weather"
)

func main() {
	cfg := configs.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Println("cfg")
	fmt.Println(cfg)
	client, _ := weatherClient.NewWeatherAPIConfig(cfg.WeatherApiToken)
	weatherService := domain.NewWeatherServiceImpl(client)
	subsRepo := domain.NewSubscriptionRepo(&cfg)
	emailSender := domain.NewEmailSender(&cfg)
	subscriptionService := domain.NewSubscriptionServiceImpl(subsRepo, emailSender)

	s, _ := gocron.NewScheduler()

	s.NewJob(
		gocron.DurationJob(
			1*time.Hour,
		),
		gocron.NewTask(
			emailSender.SendWeather,
			subsRepo,
			weatherService,
		),
	)
	s.Start()

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	router := gin.Default()

	router.GET("/weather", weatherHandler.GetWeather)
	router.POST("/subscribe", subscriptionHandler.Subscribe)
	router.GET("/confirm/:token", subscriptionHandler.Confirm)
	router.GET("/unsubscribe/:token", subscriptionHandler.Unsubscribe)

	router.Run("0.0.0.0:8080")
}
