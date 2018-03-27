package main

import (
	"./config"
	"./store"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type rs struct {
	Data []dataStru `json:"data"`
}

type dataStru struct {
	Pm25       string `json:"ps_pm25"`
	Forecast6d struct {
		Info []forecastStru `json:"info"`
	} `json:"forecast6d"`
	Observe observeStru `json:"observe"`
}

type forecastStru struct {
	Date               string `json:"date"`
	TemperatureDay     string `json:"temperature_day"`
	TemperatureNight   string `json:"temperature_night"`
	WeatherDay         string `json:"weather_day"`
	WeatherNight       string `json:"weather_night"`
	WindDirectionDay   string `json:"wind_direction_day"`
	WindDirectionNight string `json:"wind_direction_night"`
	WindPowerDay       string `json:"wind_power_day"`
	WindPowerNight     string `json:"wind_power_night"`
}

type observeStru struct {
	Humidity         string `json:"humidity"`
	Temperature      string `json:"temperature"`
	Weather          string `json:"weather"`
	WindDirection    string `json:"wind_direction"`
	WindDirectionNum string `json:"wind_direction_num"`
	WindPowerNum     string `json:"wind_power_num"`
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	redisClient := store.GetRedis()

	apis := app.Group("/apis")

	apis.GET("/weather/:city", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		city := c.Param("city")
		rs := getJsonData(city)

		c.JSON(200, gin.H{
			"status": 0,
			"city":   city,
			"data":   rs,
		})
	})

	apis.POST("/feedback", func(c *gin.Context) {
		fb := c.PostForm("feedback")
		if fb != "" {
			key := strings.Join(strings.Split(time.Now().String(), " ")[:2], "-")
			err := redisClient.Set(key, fb, 0).Err()
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   "ok",
		})
	})

	app.Run(config.Port)
}

func getJsonData(city string) rs {
	city = url.QueryEscape(city + config.Suffix)
	fmt.Println(city)
	wUrl := strings.Replace(config.WeatherUrl, "${city}", city, -1)

	res, err := http.Get(wUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	result := rs{}

	dec := mahonia.NewDecoder("GB18030")
	_, cdate, transErr := dec.Translate(body, true)

	if transErr != nil {
		cdate = body
	}

	json.Unmarshal(cdate, &result)
	return result
}
