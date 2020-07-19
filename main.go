package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getStockHistory(symbol string, cookies []*http.Cookie) (*stockHistory, error) {
	remoteURL := "https://stock.xueqiu.com/v5/stock/chart/kline.json"
	values := url.Values{}
	values.Add("symbol", symbol)
	values.Add("begin", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	values.Add("period", "day")
	values.Add("type", "before")
	values.Add("count", "-7")
	values.Add("indicator", "line,pe,pb,ps,pcf,market_capital,agt,ggt,balance")
	body, err := fetch(remoteURL, values, cookies)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	data := m["data"].(map[string]interface{})
	column := data["column"].([]interface{})
	fmt.Println(column)

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
	body, err := fetch(remoteURL, values, nil)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	data := m["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	todays := make([]stockToday, len(list))
	for i, v := range list {
		item := v.(map[string]interface{})
		todays[i].symbol, _ = item["symbol"].(string)
		todays[i].percent, _ = item["percent"].(float64)
		todays[i].current, _ = item["current"].(float64)
		todays[i].currentYearPercent, _ = item["current_year_percent"].(float64)
		todays[i].name, _ = item["name"].(string)
	}
	return todays, nil
}

func main() {

	//var token string
	resp, err := http.Get("https://xueqiu.com/")
	if err != nil {
		return
	}
	cookies := resp.Cookies()
	for _, v := range cookies {
		fmt.Println(v.Name)
		fmt.Println(v.Value)
	}

	history, _ := getStockHistory("SH601798", cookies)
	fmt.Println(history)

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
