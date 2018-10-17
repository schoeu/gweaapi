package main

import (
	"./config"
	"./lunar"
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
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	layout = "2006-01-02 15:04:05"
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
				_, err := db.Exec(`update weathers.cityinfo set citylist = ? where username = ?`, citys, username)
				utils.ErrHandle(err)
			} else {
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
		city := c.Query("city")
		vin := c.DefaultQuery("vin", "")
		esn := c.DefaultQuery("esn", "")

		rs := violation.GetVioInfo(city, lpn, vin, esn)

		b, err := json.Marshal(rs)
		utils.ErrHandle(err)

		jsonStr := string(b)
		// INSERT INTO services (user_id, violation, violation_info) VALUES ('122', '1', '{"a":1}') ON DUPLICATE KEY UPDATE violation_info= '{"a":3}
		_, dbErr := db.Exec(`INSERT INTO services (openid, violation, vioinfo, lpn, vin ,esn, city) VALUES (?, 1, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE vioinfo= ?`, uid, jsonStr, lpn, vin, esn, city, jsonStr)
		utils.ErrHandle(dbErr)

		c.JSON(200, gin.H{
			"status": 0,
			"data":   rs,
		})
	})

	apis.GET("/carpre", func(c *gin.Context) {
		lpn := c.Query("lpn")

		rs := violation.GetCarPre(lpn)

		c.JSON(200, gin.H{
			"status": 0,
			"data":   rs,
		})
	})

	type vioRsInfo struct {
		FineCount  int    `json:"fineCount"`
		ItemCount  int    `json:"itemCount"`
		Lpn        string `json:"lpn"`
		PointCount int    `json:"pointCount"`
	}
	apis.GET("/vioinfo", func(c *gin.Context) {
		uid := c.DefaultQuery("uid", "")
		if uid != "" {
			var info, city, lpn, vin, esn, date string
			err := db.QueryRow(`select vioinfo, city, lpn, vin, esn, update_date from services where openid = ?`, uid).Scan(&info, &city, &lpn, &vin, &esn, &date)
			utils.ErrHandle(err)

			if lpn == "" {
				c.JSON(200, gin.H{
					"status": 0,
					"msg":    "No data.",
					"data":   []string{},
				})
				return
			}

			// t, _ := time.Parse(layout, date)
			t, _ := time.ParseInLocation(layout, date, time.Local)
			timeNow := time.Now()
			sub := timeNow.Sub(t)

			fmt.Println("sub time: ", sub)

			vi := violation.VioInfo{}
			if config.SubTime*time.Hour < sub && city != "" && lpn != "" {
				fmt.Println("new request.", city, lpn, vin, esn)
				vi = violation.GetVioInfo(city, lpn, vin, esn)
				b, err := json.Marshal(vi)
				utils.ErrHandle(err)

				jsonStr := string(b)
				_, dbErr := db.Exec(`update services set vioinfo = ?, update_date  = ? where openid = ?`, jsonStr, timeNow.Format(layout), uid)
				utils.ErrHandle(dbErr)
			} else {
				err = json.Unmarshal([]byte(info), &vi)
				utils.ErrHandle(err)
				fmt.Println("cache.")
			}

			vri := vioRsInfo{}

			if vi.ErrCode == 0 {
				result := vi.Result
				lists := result.Lists
				vri.Lpn = result.Hphm
				vri.ItemCount = len(lists)

				for _, v := range lists {
					money, _ := strconv.Atoi(v.Money)
					fen, _ := strconv.Atoi(v.Fen)
					vri.FineCount += money
					vri.PointCount += fen
				}
				fmt.Println(vri)
			}
			c.JSON(200, gin.H{
				"status": 0,
				"msg":    "",
				"data":   vri,
			})
		} else {
			c.JSON(200, gin.H{
				"status": 1,
				"msg":    "no uid.",
				"data":   "",
			})
		}
	})

	apis.GET("/gettime", func(c *gin.Context) {
		t := []string{
			time.Now().Format("2006-01-02"),
		}
		rs := lunar.Lunar(time.Now().Format("20060102"))
		t = append(t, rs...)
		c.JSON(200, gin.H{
			"status": 0,
			"data":   t,
		})
	})

	app.Run(config.Port)
}
