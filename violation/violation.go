package violation

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"../utils"
	"../config"
)

func AddCars(num, vin, esn string) string {
	return ""
}

func DeleteCars(n string) {

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

func Get(lpn, vin, esn string) {
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
	fmt.Println(string(respBytes))
}