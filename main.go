package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"./config"
	"github.com/json-iterator/go"
	"github.com/gin-gonic/gin"
)

type rs struct {
	Status string `json:"status"`
	Data []dataStru `json:"data"`
}

type dataStru struct {
	Pm25 string `json:"ps_pm25"`
	Forecast6d struct {
		Info []forecastStru `json:"info"`
	} `json:"forecast6d"`
	Observe observeStru `json:"observe"`
}

type forecastStru struct {
	Date string `json:"date"`
	TemperatureDay string `json:"temperature_day"`
	TemperatureNight string `json:"temperature_night"`
	WeatherDay string `json:"weather_day"`
	WeatherNight string `json:"weather_night"`
	WindDirectionDay string `json:"wind_direction_day"`
	WindDirectionNight string `json:"wind_direction_night"`
	WindPowerDay string `json:"wind_power_day"`
	WindPowerNight string `json:"wind_power_night"`
}

type observeStru struct {
	Humidity string `json:"humidity"`
	Temperature string `json:"temperature"`
	Weather string `json:"weather"`
	WindDirection string `json:"wind_direction"`
	WindDirectionNum string `json:"wind_direction_num"`
	WindPowerNum string `json:"wind_power_num"`
	
}

func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	res, err := http.Get(config.WeatherUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	result := rs{}
	json.Unmarshal(body, &result)

	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	apis := app.Group("/api")

	apis.GET("/weather/:city", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "application/json;charset=gbk")

		city := c.Param("city")

		c.JSON(200, gin.H{
			"status": "1",
			"city": city,
			"data":   result,
		})
	})

	app.Run(config.Port)
}
