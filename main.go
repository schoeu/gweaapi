package main

import (
	"./config"
	"./store"
	"./utils"
	"./violation"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

	db := utils.OpenDb("mysql", config.MysqlUrl)
	defer db.Close()

	result := utils.Rs{}
	cityr := utils.Cityrs{}
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
				rs := utils.GetJsonData(city)
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
			matchArr := utils.SearchText(key, cityr)
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
		openData := utils.GetOpenJSON(code)
		c.JSON(200, gin.H{
			"status": 0,
			"data":   openData.Openid,
		})
	})

	apis.GET("/violation", func(c *gin.Context) {
		uid := c.Query("uid")
		lpn := c.Query("lpn")
		vin := c.DefaultQuery("vin", "")
		esn := c.DefaultQuery("esn", "")

		carid := violation.AddCars(lpn, vin, esn)
		carString := strconv.Itoa(carid)
		fmt.Println(carString)
		rs := violation.GetCarsInfo(carString)

		if carid != 0 {
			violation.DeleteCars(carString)
		}

		b, err := json.Marshal(rs)
		utils.ErrHandle(err)

		jsonStr := string(b)
		// INSERT INTO services (user_id, violation, violation_info) VALUES ('122', '1', '{"a":1}') ON DUPLICATE KEY UPDATE violation_info= '{"a":3}
		_, dbErr := db.Exec(`INSERT INTO services (openid, violation, vioinfo) VALUES (?, 1, ?) ON DUPLICATE KEY UPDATE vioinfo= ?`, uid, jsonStr, jsonStr)
		utils.ErrHandle(dbErr)

		c.JSON(200, gin.H{
			"status": 0,
			"data":   rs,
		})
	})

	apis.GET("/deletecar", func(c *gin.Context) {
		n := c.Query("n")
		violation.DeleteCars(n)
		c.JSON(200, gin.H{
			"status": 0,
			"data":   "done",
		})
	})

	apis.GET("/vioinfo", func(c *gin.Context) {
		uid := c.DefaultQuery("uid", "")
		var info string
		if uid != "" {
			rows, err := db.Query(`select vioinfo from services where openid = ?`, uid)
			utils.ErrHandle(err)
			for rows.Next() {
				err := rows.Scan(&info)
				utils.ErrHandle(err)
			}
			err = rows.Err()
			utils.ErrHandle(err)
			defer rows.Close()
			fmt.Println(info)
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   info,
		})
	})

	app.Run(config.Port)
}
