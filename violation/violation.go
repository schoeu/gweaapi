package violation

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"strings"
	"io/ioutil"
	"../utils"
	"../config"
)

func DeleteCars(n string) {
	rsUrl := strings.Replace(config.DeleteCarUrl, "{name}", n, -1)
	req, err := http.NewRequest("DELETE", rsUrl, nil)
	for k, v := range config.HeadersInfo{
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.ErrHandle(err)

	defer resp.Body.Close()
	fmt.Println("delete car: ", n)
}

func GetCarsInfo(n string) {
	rsUrl := strings.Replace(config.InfoUrl, "{num}", n, -1)
	req, err := http.NewRequest("GET", rsUrl, nil)
	for k, v := range config.HeadersInfo{
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.ErrHandle(err)

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	utils.ErrHandle(err)

	fmt.Println("get car info: ", string(respBytes))
}

type dataStruct struct {
	ModelId string `json:"modelId"`
	Vin string `json:"vin"`
	RootbrandId string `json:"rootbrandId"`
	Lpn string `json:"lpn"`
	Esn string `json:"esn"`
	SubbrandId string `json:"subbrandId"`
}

type addStruct struct {
	Id int `json:"id"`
}

func AddCars(lpn, vin, esn string) int {
	ds := dataStruct{}
	ds.Vin = lpn
	ds.Lpn = vin
	ds.Esn = esn

	jsonStr, err := json.Marshal(ds)
	utils.ErrHandle(err)

	reader := bytes.NewReader(jsonStr)

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.AddCarUrl, reader)
	for k, v := range config.HeadersInfo{
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