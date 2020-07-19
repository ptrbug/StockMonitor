package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func getStockHistory(symbol string) (*stockHistory, error) {
	history := &stockHistory{}
	return history, nil
}

func getStockToday(topPercentCount int) ([]stockToday, error) {
	remoteURL := "https://xueqiu.com/service/v5/stock/screener/quote/list"
	values := url.Values{}
	values.Add("page", "1")
	values.Add("size", strconv.Itoa(topPercentCount))
	values.Add("order", "desc")
	values.Add("orderby", "percent")
	values.Add("order_by", "percent")
	values.Add("market", "CN")
	values.Add("type", "sh_sz")
	body, err := fetch(remoteURL, values)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	var todays []stockToday
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	fmt.Print(m)
	return todays, nil
}

func main() {

	todays, _ := getStockToday(10)
	fmt.Println(todays)

	/*
		t := time.NewTimer(time.Second * 5)
		defer t.Stop()
		for {
			<-t.C
			todays, err := getStockToday(500)
			if err != nil {
			} else {
				for _, today := range todays {
					fmt.Println(today)
				}
			}

			t.Reset(time.Second * 5)
		}
	*/
}
