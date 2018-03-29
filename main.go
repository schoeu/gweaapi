package main

import (
	"./config"
	"./store"
	"./utils"
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
	SunriseTime        string `json:"sunriseTime"`
	SunsetTime         string `json:"sunsetTime"`
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

type cityrs struct {
	Data []cityinfo `json:"data"`
}

type cityinfo struct {
	K int      `json:"k"`
	N string   `json:"n"`
	S []string `json:"s"`
}

type rsType []string

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func main() {
	during := time.Minute * 30
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	store.GetRedis()
	result := rs{}
	cityr := cityrs{}
	apis := app.Group("/api")
	apis.GET("/weather/:city", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		city := c.Param("city")

		cityTemp := store.GetData(city)
		if cityTemp != "" {
			err := json.Unmarshal([]byte(cityTemp.(string)), &result)
			utils.ErrHandle(err)

			c.JSON(200, gin.H{
				"status": 0,
				"city":   city,
				"data":   cityTemp,
				"from":   "cache",
			})
		} else {
			rs := getJsonData(city)
			b, err := json.Marshal(rs)
			utils.ErrHandle(err)

			store.SetData(city, b, during)
			c.JSON(200, gin.H{
				"status": 0,
				"city":   city,
				"data":   rs,
			})
		}
	})

	apis.POST("/feedback", func(c *gin.Context) {
		fb := c.PostForm("feedback")
		if fb != "" {
			key := strings.Join(strings.Split(time.Now().String(), " ")[:2], "-")
			store.SetData(key, fb, 0)
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   "ok",
		})
	})

	apis.GET("/search", func(c *gin.Context) {
		key := c.Query("key")

		citys := store.GetData("citylist")
		err := json.Unmarshal([]byte(citys.(string)), &cityr)
		utils.ErrHandle(err)

		if key == "" {
			c.JSON(200, gin.H{
				"status": 0,
				"data":   cityr,
			})
		} else {
			matchArr := searchText(key, cityr)
			c.JSON(200, gin.H{
				"status": 0,
				"data":   matchArr,
			})
		}
	})

	app.Run(config.Port)
}

func getJsonData(city string) rs {
	city = url.QueryEscape(city + config.Suffix)
	wUrl := strings.Replace(config.WeatherUrl, "${city}", city, -1)

	res, err := http.Get(wUrl)
	utils.ErrHandle(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	utils.ErrHandle(err)

	result := rs{}
	dec := mahonia.NewDecoder("GB18030")
	_, cdate, transErr := dec.Translate(body, true)

	if transErr != nil {
		cdate = body
	}

	json.Unmarshal(cdate, &result)
	return result
}

func searchText(key string, data cityrs) []rsType {
	d := data.Data
	var rs []rsType
	for _, v := range d {
		var rsInfo rsType
		n := v.N
		s := v.S
		if strings.Contains(n, key) {
			rsInfo = append(rsInfo, n)
		}
		for _, val := range s {
			if strings.Contains(val, key) {
				rsInfo = append(rsInfo, n, val)
			}
			if len(rsInfo) > 0 {
				rs = append(rs, rsInfo)
				rsInfo = rsType{}
			}
		}
	}
	return rs
}
