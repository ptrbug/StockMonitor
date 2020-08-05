package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

//MaxHistorySize max history size
const MaxHistorySize = 20

type daily struct {
	maxPrice   float64
	maxPercent float64
	Close      [MaxHistorySize]float64
}

func getDaily(query int64, symbol string, cookies []*http.Cookie) (*daily, error) {
	remoteURL := "https://stock.xueqiu.com/v5/stock/chart/kline.json"
	values := url.Values{}
	values.Add("symbol", symbol)
	values.Add("begin", strconv.FormatInt(query*1000, 10))
	values.Add("period", "day")
	values.Add("type", "before")
	values.Add("count", fmt.Sprintf("-%d", MaxHistorySize))
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

	day := &daily{}
	for i := len(item) - 1; i >= 0; i-- {
		array := item[i].([]interface{})
		day.Close[len(item)-1-i] = array[5].(float64)
	}
	if day.Close[0] <= 0 {
		return nil, fmt.Errorf("data error")
	}
	day.maxPrice, day.maxPercent = calcMaxmaxPriceAndPercent(day.Close[0])
	return day, nil
}
