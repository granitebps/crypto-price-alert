package vendors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/granitebps/crypto-price-alert/helper"
	"github.com/granitebps/crypto-price-alert/types"
)

func GetPriceCoingecko(alert types.Alert) (int, error) {
	var data types.CoingeckoResponse

	pair, err := aliasCoingecko(alert.Pair)
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=idr&include_last_updated_at=true", pair), nil)
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

	lastPrice, _ := strconv.Atoi((fmt.Sprintf("%v", data[pair]["idr"])))
	return lastPrice, nil
}

func aliasCoingecko(pair string) (string, error) {
	var alias types.Alias

	filename := "alias.json"
	jsonData, err := helper.ReadJsonFile(filename)
	if err != nil {
		return "", err
	}

	result, err := helper.ParseAliasData(jsonData, alias)
	if err != nil {
		return "", err
	}

	return result["coingecko"][pair], nil
}
