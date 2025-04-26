package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiKey = "97653189e01ec25518b58f4c5510c25d"
const apiUrl = "http://api.openweathermap.org/data/2.5/weather"

func GetWeather(city string) {
	escapedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=imperial", apiUrl, escapedCity, apiKey)
	resp, err := http.Get(url)
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

	// Print the raw JSON response for debugging
	fmt.Println("Raw response:", string(body))

	var weatherData map[string]interface{}
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Weather Data:", weatherData)
}

// Main function
func main() {
	city := "Los Angeles" // Example city
	GetWeather(city)
}
