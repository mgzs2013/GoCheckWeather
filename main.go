// main.go
package main

import (
	"GoCheckWeather/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WeatherData holds the data to be displayed in the HTML template
type WeatherData struct {
	City        string
	Temperature float64
	Humidity    float64
	Condition   string
	Error       string         // Add an Error field to hold error messages
	Forecast    []ForecastData // New field for forecast data
}

// ForecastData holds forecast information
type ForecastData struct {
	Date        string  // Date of the forecast
	Temperature float64 // Forecast temperature
	Condition   string  // Forecast condition
}

func main() {
	r := gin.Default()

	// Serve the HTML page
	r.LoadHTMLGlob("templates/*")

	// Render the main page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Handle form submission
	r.POST("/weather", func(c *gin.Context) {
		city := c.PostForm("city")
		weatherData := service.GetWeather(city)
		forecastData := service.GetForecast(city) // Get the forecast data

		if weatherData != nil {
			c.HTML(http.StatusOK, "index.html", WeatherData{
				City:        city,
				Temperature: weatherData["main"].(map[string]interface{})["temp"].(float64),
				Humidity:    weatherData["main"].(map[string]interface{})["humidity"].(float64),
				Condition:   weatherData["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
				Forecast:    parseForecastData(forecastData), // Parse forecast data
			})
		} else {
			// Pass an error message to the template
			c.HTML(http.StatusOK, "index.html", WeatherData{
				Error: "Could not retrieve weather data. Please check the city name and try again.",
			})
		}
	})

	// Start the server
	r.Run(":8080") // Listen on port 8080
}

// Function to parse forecast data into a more usable format
func parseForecastData(forecastData []map[string]interface{}) []ForecastData {
	var forecasts []ForecastData
	for _, item := range forecastData {
		if main, ok := item["main"].(map[string]interface{}); ok {
			date := item["dt_txt"].(string) // Get the date from the forecast
			temperature := main["temp"].(float64)
			condition := item["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)

			forecasts = append(forecasts, ForecastData{
				Date:        date,
				Temperature: temperature,
				Condition:   condition,
			})
		}
	}
	return forecasts
}
