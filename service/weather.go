// weather.go
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiKey = "YOUR_API_KEY" // Replace with your actual API key
const apiUrl = "http://api.openweathermap.org/data/2.5/weather"

var cache = NewWeatherCache(10 * time.Minute) // Cache for 10 minutes

// GetWeather fetches the weather data for a given city
func GetWeather(city string) {
	if cachedResponse, found := cache.Get(city); found {
		fmt.Println("Fetching data from cache...")
		printWeatherData(cachedResponse.Data, city)
		return
	}

	escapedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?q=%s&appid=%s&units=imperial", apiUrl, escapedCity, apiKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var weatherData map[string]interface{}
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Store the data in the cache
	cache.Set(city, weatherData)
	printWeatherData(weatherData, city)
}

// Helper function to print weather data
func printWeatherData(weatherData map[string]interface{}, city string) {
	fmt.Printf("Weather in %s:\n", city)
	if main, ok := weatherData["main"].(map[string]interface{}); ok {
		fmt.Printf("Temperature: %.2f°F\n", main["temp"])
		fmt.Printf("Humidity: %.2f%%\n", main["humidity"])
	}
	if weather, ok := weatherData["weather"].([]interface{}); ok && len(weather) > 0 {
		fmt.Printf("Condition: %s\n", weather[0].(map[string]interface{})["description"])
	}
}
