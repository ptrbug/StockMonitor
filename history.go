package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func isFileExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func downloadAllDaily(cookies []*http.Cookie) (map[string]*daily, error) {

	curs, err := getTopPercent(10000)
	if err != nil {
		return nil, fmt.Errorf("error")
	}

	all := make(map[string]*daily, len(curs))
	for _, v := range curs {
		day, err := getDaily(v.symbol, cookies)
		if err != nil {
			fmt.Printf("getStockHistory : %s error\n", v.symbol)
		} else {
			all[v.symbol] = day
		}
	}
	return all, nil
}

func loadAllDaily() (map[string]*daily, error) {

	resp, err := http.Get("https://xueqiu.com/")
	if err != nil {
		return nil, fmt.Errorf("get cookies failed")
	}
	cookies := resp.Cookies()
	resp.Body.Close()

	tmLocal := time.Now().Local()
	day := tmLocal.Day()
	hour := tmLocal.Hour()
	minute := tmLocal.Minute()
	if day == 6 || day == 7 || (hour > 15 || (hour == 15 && minute > 30)) {
		remoteURL := "https://stock.xueqiu.com/v5/stock/chart/kline.json"
		values := url.Values{}
		values.Add("symbol", "SH000001")
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

		fmt.Println(item)
	}
	year, month, day := time.Now().Date()
	filepath := fmt.Sprintf("daily/%d-%d-%d.json", year, month, day)

	if isFileExist(filepath) {
		data, err := ioutil.ReadFile(filepath)
		if err == nil {
			all := make(map[string]*daily, 10000)
			err = json.Unmarshal(data, &all)
			if err == nil {
				return all, nil
			}
		}
	}

	alllHistory, err := downloadAllDaily(cookies)
	if err == nil {
		data, err := json.Marshal(&alllHistory)
		if err == nil {
			os.MkdirAll(path.Dir(filepath), os.ModePerm)
			ioutil.WriteFile(filepath, data, os.ModePerm)
		}
		return alllHistory, err
	}

	return nil, fmt.Errorf("loadAllDaily failed")
}
