package main

import (
	"./config"
	"./store"
	"./utils"
	"./violation"
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
	City       string `json:"city"`
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

type sessionBody struct {
	Session string `json:"session_key"`
	Openid  string `json:"openid"`
}

func main() {
	during := time.Minute * 30
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	store.GetRedis()

	db := utils.OpenDb("mysql", config.MysqlUrl)
	defer db.Close()

	result := rs{}
	cityr := cityrs{}
	apis := app.Group("/weather/api")
	apis.GET("/weather/:city", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		city := c.Param("city")

		// muti citys
		mcity := strings.Split(city, ",")
		if len(mcity) > 1 {

		} else {
			cityTemp := store.GetData(city)
			if cityTemp == nil {
				err := json.Unmarshal([]byte(cityTemp.(string)), &result)
				utils.ErrHandle(err)

				c.JSON(200, gin.H{
					"status": 0,
					"city":   city,
					"data":   result,
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
		}
	})

	apis.POST("/feedback", func(c *gin.Context) {
		fb := c.PostForm("feedback")
		username := c.Query("username")
		if fb != "" {
			_, err := db.Exec(`insert into userinfo (username, comments) values (?, ?)`, username, fb)
			utils.ErrHandle(err)
		}

		c.JSON(200, gin.H{
			"status": 0,
			"data":   "ok",
		})
	})

	apis.GET("/search", func(c *gin.Context) {
		key := c.Query("key")

		citys := store.GetData("citylist")
		if citys == "" {
			citys = config.CityMap
			store.SetData("citylist", citys, 0)
		}

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

	apis.GET("/addcity", func(c *gin.Context) {
		username := c.Query("username")
		citys := c.Query("citys")
		var name string
		if username != "" {
			err := db.QueryRow(`select citylist from cityinfo where username = ?`, username).Scan(&name)
			utils.ErrHandle(err)
		}
		if citys != "" {
			fmt.Println("name", name)
			if name != "" {
				fmt.Println("update", username, citys)
				_, err := db.Exec(`update weathers.cityinfo set citylist = ? where username = ?`, citys, username)
				utils.ErrHandle(err)
			} else {
				fmt.Println("insert", username, citys)
				_, err := db.Exec(`insert into cityinfo (username, citylist) values (?, ?)`, username, citys)
				utils.ErrHandle(err)
			}
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   "ok",
		})
	})

	apis.GET("/getcity", func(c *gin.Context) {
		username := c.Query("username")
		var citys string
		if username != "" {
			rows, err := db.Query(`select citylist from cityinfo where username = ?`, username)
			utils.ErrHandle(err)
			for rows.Next() {
				err := rows.Scan(&citys)
				utils.ErrHandle(err)
			}
			err = rows.Err()
			utils.ErrHandle(err)
			defer rows.Close()
			c.JSON(200, gin.H{
				"status": 0,
				"data":   strings.Split(citys, ","),
			})
		} else {
			c.JSON(200, gin.H{
				"status": 1,
				"data":   "no data.",
			})
		}
	})

	apis.GET("/getopenid", func(c *gin.Context) {
		code := c.Query("code")
		openData := getOpenJSON(code)
		c.JSON(200, gin.H{
			"status": 0,
			"data":   openData.Openid,
		})
	})

	apis.GET("/violation", func(c *gin.Context) {
		lpn := c.Query("lpn")
		vin := c.DefaultQuery("vin", "")
		esn := c.DefaultQuery("esn", "")
		
		fmt.Println(lpn, vin, esn)
		violation.Get(lpn, vin, esn)
		// c.JSON(200, gin.H{
		// 	"status": 0,
		// 	"data":   openData.Openid,
		// })
	})

	app.Run(config.Port)
}

func getJsonData(city string) rs {
	city = url.QueryEscape(city + config.Suffix)
	wUrl := strings.Replace(config.WeatherUrl, "${city}", city, -1)

	body := get(wUrl)
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

func getOpenJSON(code string) sessionBody {
	rsUrl := strings.Replace(config.CodeUrl, "${jscode}", code, -1)
	body := get(rsUrl)
	s := sessionBody{}
	json.Unmarshal(body, &s)
	return s
}

func get(url string) []byte {
	res, err := http.Get(url)
	utils.ErrHandle(err)
	body, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()
	utils.ErrHandle(err)
	return body
}
