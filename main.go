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
	item := data["item"].([]interface{})

	history := &stockHistory{}
	first := len(item) - 1
	last := first - 4
	if last < 0 {
		last = 0
	}
	for i := first; i >= last; i-- {
		array := item[i].([]interface{})
		history.close[first-i] = array[5].(float64)
	}
	return history, nil
}

func getAllHistory(cookies []*http.Cookie) (map[string]*stockHistory, error) {
	todays, err := getStockToday(1000)
	if err != nil {
		return nil, fmt.Errorf("error")
	}

	allHistory := make(map[string]*stockHistory, len(todays))
	for _, v := range todays {
		history, err := getStockHistory(v.symbol, cookies)
		if err != nil {
			fmt.Printf("getStockHistory : %s error\n", v.symbol)
		} else {
			allHistory[v.symbol] = history
		}
	}
	return allHistory, nil
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

	allHistory, _ := getAllHistory(cookies)
	fmt.Println(allHistory)

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
