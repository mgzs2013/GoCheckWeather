package main

import (
	"GoCheckWeather/service"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a city name.")
		return
	}
	city := os.Args[1]
	service.GetWeather(city)
}
