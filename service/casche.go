// cache.go
package service

import (
	"sync"
	"time"
)

// WeatherResponse represents the structure of the weather data
type WeatherResponse struct {
	Timestamp time.Time
	Data      map[string]interface{}
}

// WeatherCache holds the cached weather data
type WeatherCache struct {
	data map[string]WeatherResponse
	mu   sync.RWMutex
	ttl  time.Duration // Time to live for cached data
}

// NewWeatherCache creates a new WeatherCache
func NewWeatherCache(ttl time.Duration) *WeatherCache {
	return &WeatherCache{
		data: make(map[string]WeatherResponse),
		ttl:  ttl,
	}
}

// Get retrieves weather data from the cache
func (c *WeatherCache) Get(city string) (WeatherResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	response, found := c.data[city]
	if found && time.Since(response.Timestamp) < c.ttl {
		return response, true
	}
	return WeatherResponse{}, false
}

// Set stores weather data in the cache
func (c *WeatherCache) Set(city string, data map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[city] = WeatherResponse{
		Timestamp: time.Now(),
		Data:      data,
	}
}
