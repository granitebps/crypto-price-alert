package vendors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/granitebps/crypto-price-alert/types"
)

func GetPriceIndodax(alert types.Alert) (int, error) {
	var data types.IndodaxResponse
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://indodax.com/api/ticker/%s", alert.Pair), nil)
	if err != nil {
		return 0, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	lastPrice, _ := strconv.Atoi((fmt.Sprintf("%v", data["ticker"]["last"])))
	return lastPrice, nil
}
