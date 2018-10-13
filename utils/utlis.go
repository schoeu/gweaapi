package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/json-iterator/go"
	"github.com/axgle/mahonia"
	"net/url"
	"net/http"
	"strings"
	"io/ioutil"
	"../config"
)

type rsType []string

type sessionBody struct {
	Session string `json:"session_key"`
	Openid  string `json:"openid"`
}

type Cityrs struct {
	Data []cityinfo `json:"data"`
}

type cityinfo struct {
	K int      `json:"k"`
	N string   `json:"n"`
	S []string `json:"s"`
}

type Rs struct {
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

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// 创建数据库链接
func OpenDb(dbTyepe string, dbStr string) *sql.DB {
	if dbTyepe == "" {
		dbTyepe = "mysql"
	}
	db, err := sql.Open(dbTyepe, dbStr)
	ErrHandle(err)

	err = db.Ping()
	ErrHandle(err)
	return db
}

func GetJsonData(city string) Rs {
	city = url.QueryEscape(city + config.Suffix)
	wUrl := strings.Replace(config.WeatherUrl, "${city}", city, -1)

	body := get(wUrl)
	result := Rs{}
	dec := mahonia.NewDecoder("GB18030")
	_, cdate, transErr := dec.Translate(body, true)

	if transErr != nil {
		cdate = body
	}

	json.Unmarshal(cdate, &result)
	return result
}

func SearchText(key string, data Cityrs) []rsType {
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

func GetOpenJSON(code string) sessionBody {
	rsUrl := strings.Replace(config.CodeUrl, "${jscode}", code, -1)
	body := get(rsUrl)
	s := sessionBody{}
	json.Unmarshal(body, &s)
	return s
}

func get(url string) []byte {
	res, err := http.Get(url)
	ErrHandle(err)
	body, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()
	ErrHandle(err)
	return body
}