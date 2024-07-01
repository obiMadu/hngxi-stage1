package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// define constants
const webPort = ":8082"
const ipInfoURL = "https://ipinfo.io"
const ipWeatherURL = "https://api.openweathermap.org/data/2.5/weather"

type ipLocInfo struct {
	City string `json:"city"`
	Loc  string `json:"loc"`
}

type ipWeatherInfo struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v\n", err)
	}

	// make a gin router
	mux := gin.Default()

	// use cors
	mux.Use(cors.Default())

	// setup route
	mux.GET("/api/hello", func(c *gin.Context) {
		// get client IP
		ip := c.ClientIP()

		// get visitor_name query
		visitor_name := strings.Trim(c.Query("visitor_name"), "\"")
		if visitor_name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "you did not provide a visitor_name query",
			})
			return
		}

		// get IP location details from ipinfo.io
		var ipLocInfo ipLocInfo
		url := fmt.Sprintf("%s/%s/json", ipInfoURL, ip)

		// call for IP location info
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to get ip location info from ipinfo.io",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			err := json.NewDecoder(resp.Body).Decode(&ipLocInfo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "failed to process ip location info from ipinfo.io",
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to get ip location info from ipinfo.io",
			})
			return
		}

		// get temperation info from OpenWeatherMap
		// get latitude and longitude
		var lat, lon float64
		fmt.Sscanf(ipLocInfo.Loc, "%f,%f", &lat, &lon)

		// get api key from environment
		apiKey := os.Getenv("OPENWEATHERMAP_APIKEY")

		// call OpenWeatherMap
		url = fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", ipWeatherURL, lat, lon, apiKey)
		resp, err = http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to contact ip weather service",
			})
			return
		}

		var ipWeatherInfo ipWeatherInfo
		if resp.StatusCode == 200 {
			if err := json.NewDecoder(resp.Body).Decode(&ipWeatherInfo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "failed to parse request from weather api",
				})
				return
			}
		} else {
			var errResp map[string]any
			_ = json.NewDecoder(resp.Body).Decode(&errResp)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "got an error response from weather api",
				"error": gin.H{
					"responseCode": resp.StatusCode,
					"error":        errResp,
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"client_ip": ip,
			"location":  ipLocInfo.City,
			"greeting":  fmt.Sprintf("Hello, %s!, the temperature is %.0f degrees Celcius in %s", visitor_name, ipWeatherInfo.Main.Temp, ipLocInfo.City),
		})
	})

	// exec router
	mux.Run(webPort)
}
