package violation

import (
	"../config"
	"../utils"
	"fmt"
	"github.com/json-iterator/go"
	"net/url"
)

type VioInfo struct {
	ErrCode    int    `json:"error_code"`
	Reason     string `json:"reason"`
	Resultcode string `json:"resultcode"`
	Result     VioRs  `json:"result"`
}

type VioRs struct {
	City     string    `json:"city"`
	Hphm     string    `json:"hphm"`
	Hpzl     string    `json:"hpzl"`
	Lists    []VioList `json:"lists"`
	Province string    `json:"province"`
}

type VioList struct {
	Act       string `json:"act"`
	Archiveno string `json:"archiveno"`
	Area      string `json:"area"`
	Code      string `json:"code"`
	Date      string `json:"date"`
	Fen       string `json:"fen"`
	Handled   string `json:"handled"`
	Money     string `json:"money"`
	Wzcity    string `json:"wzcity"`
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func GetVioInfo(city, lpn, vin, esn string) VioInfo {
	u, _ := url.Parse(config.JHUrl)
	q := u.Query()
	q.Set("city", city)
	q.Set("hphm", lpn)
	q.Set("classno", vin)
	q.Set("engineno", esn)
	q.Set("key", config.JHKey)
	u.RawQuery = q.Encode()

	respBytes := utils.Get(u.String())

	vi := VioInfo{}
	err := json.Unmarshal(respBytes, &vi)
	utils.ErrHandle(err)
	fmt.Println("package violation, vi: ", vi)
	return vi
}
