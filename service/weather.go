package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiKey = "97653189e01ec25518b58f4c5510c25d" // Replace with your actual API key
const apiUrl = "http://api.openweathermap.org/data/2.5/weather"
const forecastUrl = "http://api.openweathermap.org/data/2.5/forecast" // New forecast endpoint

var cache = NewWeatherCache(10 * time.Minute) // Cache for 10 minutes

// GetWeather fetches the weather data for a given city
func GetWeather(city string) map[string]interface{} {
	if cachedResponse, found := cache.Get(city); found {
		fmt.Println("Fetching data from cache...")
		return cachedResponse.Data // Return cached data

	}

	escapedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?q=%s&appid=%s&units=imperial", apiUrl, escapedCity, apiKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	var weatherData map[string]interface{}
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Store the data in the cache
	cache.Set(city, weatherData)
	return weatherData // Return the fetched data
}

// GetForecast fetches the 5-day weather forecast for a given city
func GetForecast(city string) []map[string]interface{} {
	escapedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?q=%s&appid=%s&units=imperial", forecastUrl, escapedCity, apiKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error fetching forecast data:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	var forecastData map[string]interface{}
	err = json.Unmarshal(body, &forecastData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Extract the relevant forecast data
	var forecasts []map[string]interface{}
	if list, ok := forecastData["list"].([]interface{}); ok {
		for _, item := range list {
			if forecast, ok := item.(map[string]interface{}); ok {
				forecasts = append(forecasts, forecast)
			}
		}
	}
	return forecasts // Return the forecast data
}
