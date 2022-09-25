package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"user-balance/consts"
	"user-balance/models"
)

func ConvertCurrency(currency string, amount float64) (float64, error) {
	url := consts.ExchangeCurrencyUrl
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	data := string(body)

	var curr models.CurrencyConverter
	json.Unmarshal([]byte(data), &curr)

	UsdRub := curr.Rates["RUB"]
	USDamound := amount / UsdRub
	USDcurr, ok := curr.Rates[currency]
	if !ok {
		return 0, fmt.Errorf("Wrong currency")
	}
	converted := USDamound * USDcurr

	return math.Floor(converted*100) / 100, nil
}
