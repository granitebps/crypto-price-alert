package helper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/granitebps/crypto-price-alert/types"
)

func ReadJsonFile(filename string) ([]byte, error) {
	var emptyData []byte

	jsonFile, err := os.Open(filename)
	if err != nil {
		return emptyData, err
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return emptyData, err
	}

	return jsonData, nil
}

func ParseAlertData(jsonData []byte, data []types.Alert) ([]types.Alert, error) {
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return data, err
	}

	return data, nil
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func UpdateFile(alerts []types.Alert) {
	file, _ := json.MarshalIndent(alerts, "", " ")

	_ = ioutil.WriteFile("alert.json", file, 0644)
}
