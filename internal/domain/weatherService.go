package domain

import (
	"fmt"
	weatherClient "github.com/kashifkhan0771/go-weather"
)

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherService interface {
	GetWeather(city string) (Weather, error)
}

type WeatherServiceImpl struct {
	weatherApi *weatherClient.WeatherAPIConfig
}

func NewWeatherServiceImpl(weatherApi *weatherClient.WeatherAPIConfig) *WeatherServiceImpl {
	return &WeatherServiceImpl{
		weatherApi: weatherApi,
	}
}

func (w WeatherServiceImpl) GetWeather(city string) (Weather, error) {
	options := weatherClient.Options{
		Query: city,
	}
	currentWeather, err := w.weatherApi.GetCurrentWeather(options)
	if err != nil {
		fmt.Println("Error fetching weather:", err)
		return Weather{}, err
	}

	return Weather{
		Temperature: currentWeather.Current.TempC,
		Humidity:    currentWeather.Current.Humidity,
		Description: fmt.Sprintf(`Feels like %.1f°C`, currentWeather.Current.FeelsLikeC),
	}, nil
}
