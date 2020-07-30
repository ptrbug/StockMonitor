package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//MaxHistorySize max history size
const MaxHistorySize = 5

type daily struct {
	maxPrice   float64
	maxPercent float64
	Close      [MaxHistorySize]float64
}

func calcMaxmaxPriceAndPercent(lastClosePrice float64) (float64, float64) {
	maxPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", lastClosePrice+0.005), 64)
	maxPercent = maxPrice / lastClosePrice
	fmt.Print(maxPrice)
}

func getDaily(symbol string, cookies []*http.Cookie) (*daily, error) {
	remoteURL := "https://stock.xueqiu.com/v5/stock/chart/kline.json"
	values := url.Values{}
	values.Add("symbol", symbol)
	values.Add("begin", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
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
	return day, nil
}
