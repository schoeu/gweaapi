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
	rsUrl := strings.Replace(config.DeleteCarUrl, "{name}", n)
	req, err := http.NewRequest("DELETE", rsUrl, nil)
	for k, v := range config.HeadersInfo{
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	fmt.Println("delete car: ", n)
}

func GetCarsInfo(n string) {
	
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
	id string
}

func AddCars(lpn, vin, esn string) string {
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

	resp, err := client.Do(req)
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	as := addStruct{}

	fmt.Println(string(respBytes))
}