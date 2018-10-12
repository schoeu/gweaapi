package violation

import (
	"../config"
	"../utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type dataStruct struct {
	ModelId     string `json:"modelId"`
	Vin         string `json:"vin"`
	RootbrandId string `json:"rootbrandId"`
	Lpn         string `json:"lpn"`
	Esn         string `json:"esn"`
	SubbrandId  string `json:"subbrandId"`
}

type addStruct struct {
	Id int `json:"id"`
}

type vioDetail struct {
	Lpn              string `json:"lpn"`
	CityCode         int    `json:"cityCode"`
	ViolationTime    int    `json:"violationTime"`
	Address          string `json:"address"`
	Behavior         string `json:"behavior"`
	DeductPoints     int    `json:"deductPoints"`
	FineNumber       int    `json:"fineNumber"`
	CollectionAgency string `json:"collectionAgency"`
	Code             string `json:"code"`
	ProcessState     int    `json:"processState"`
}

type rs struct {
	PointCount int         `json:"pointCount"`
	Details    []vioDetail `json:"details"`
	Lpn        string      `json:"lpn"`
	ItemCount  int         `json:"itemCount"`
	FineCount  int         `json:"fineCount"`
}

func DeleteCars(n string) {
	rsUrl := strings.Replace(config.DeleteCarUrl, "{name}", n, -1)
	req, err := http.NewRequest("DELETE", rsUrl, nil)
	for k, v := range config.HeadersInfo {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.ErrHandle(err)

	defer resp.Body.Close()
	fmt.Println("delete car: ", n)
}

func GetCarsInfo(n string) rs {
	rsUrl := strings.Replace(config.InfoUrl, "{num}", n, -1)
	req, err := http.NewRequest("GET", rsUrl, nil)
	for k, v := range config.HeadersInfo {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.ErrHandle(err)

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	utils.ErrHandle(err)

	r := rs{}
	err = json.Unmarshal(respBytes, &r)
	utils.ErrHandle(err)

	fmt.Println("get car info: ", r)

	return r
}

func AddCars(lpn, vin, esn string) int {
	ds := dataStruct{}
	ds.Vin = vin
	ds.Lpn = lpn
	ds.Esn = esn

	jsonStr, err := json.Marshal(ds)
	fmt.Println(string(jsonStr))
	utils.ErrHandle(err)

	reader := bytes.NewReader(jsonStr)

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.AddCarUrl, reader)
	for k, v := range config.HeadersInfo {
		req.Header.Add(k, v)
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	utils.ErrHandle(err)

	as := addStruct{}
	err = json.Unmarshal(respBytes, &as)
	utils.ErrHandle(err)
	fmt.Println(string(respBytes))
	return as.Id
}
