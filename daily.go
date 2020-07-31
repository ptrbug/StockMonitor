package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

//MaxHistorySize max history size
const MaxHistorySize = 5

type daily struct {
	maxPrice   float64
	maxPercent float64
	Close      [MaxHistorySize]float64
}

func calcMaxmaxPriceAndPercent(lastClosePrice float64) (float64, float64) {
	n10 := math.Pow10(2)
	maxPrice := math.Trunc((lastClosePrice*1.1+0.5/n10)*n10) / n10
	maxPercent := math.Trunc((maxPrice*100/lastClosePrice+0.5/n10)*n10) / n10
	return maxPrice, maxPercent
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
	first := len(item) - 1
	last := first - 4
	if last < 0 {
		last = 0
	}
	for i := first; i >= last; i-- {
		array := item[i].([]interface{})
		day.Close[first-i] = array[5].(float64)
	}
	if day.Close[0] <= 0 {
		return nil, fmt.Errorf("data error")
	}
	day.maxPrice, day.maxPercent = calcMaxmaxPriceAndPercent(day.Close[0])
	return day, nil
}
