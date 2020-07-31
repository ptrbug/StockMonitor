package main

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type realtime struct {
	symbol             string
	percent            float64
	current            float64
	currentYearPercent float64
	name               string
	flag               int
}

func getTopPercent(count int) ([]*realtime, error) {
	remoteURL := "https://xueqiu.com/service/v5/stock/screener/quote/list"
	values := url.Values{}
	values.Add("page", "1")
	values.Add("size", strconv.Itoa(count))
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
	curs := make([]*realtime, 0, len(list))
	for _, v := range list {

		item := v.(map[string]interface{})

		symbol, _ := item["symbol"].(string)
		if len(symbol) >= 5 && (symbol[2:5] == "688" || symbol[2:5] == "787" || symbol[2:5] == "789") {
			continue
		}
		name, _ := item["name"].(string)
		if name[0:1] == "N" {
			continue
		}

		cur := &realtime{}
		cur.symbol = symbol
		cur.percent, _ = item["percent"].(float64)
		cur.current, _ = item["current"].(float64)
		cur.currentYearPercent, _ = item["current_year_percent"].(float64)
		cur.name = name
		curs = append(curs, cur)
	}
	return curs, nil
}
