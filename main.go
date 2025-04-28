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

		if weatherData != nil {
			c.HTML(http.StatusOK, "index.html", WeatherData{
				City:        city,
				Temperature: weatherData["main"].(map[string]interface{})["temp"].(float64),
				Humidity:    weatherData["main"].(map[string]interface{})["humidity"].(float64),
				Condition:   weatherData["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
			})
		} else {
			c.HTML(http.StatusInternalServerError, "index.html", nil)
		}
	})

	// Start the server
	r.Run(":8080") // Listen on port 8080
}
