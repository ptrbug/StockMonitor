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
	"sync"
	"time"
)

type history struct {
	query int64

	due   int64
	count int
	mutex sync.Mutex
	data  map[string]*daily
}

func newHisory() *history {
	return &history{data: make(map[string]*daily, 10000)}
}

func (p *history) getAllStockSymbol() (map[string]struct{}, error) {
	curs, err := getTopPercent(10000)
	if err != nil {
		return nil, err
	}
	symbols := make(map[string]struct{}, len(curs))
	for _, v := range curs {
		symbols[v.symbol] = struct{}{}
	}
	return symbols, nil
}

func (p *history) download(query int64, symbols []string, cookies []*http.Cookie) (map[string]*daily, error) {

	fmt.Println("download start")
	percent := 0
	data := make(map[string]*daily, len(symbols))
	for _, v := range symbols {
		day, err := getDaily(query, v, cookies)
		if err != nil {
			fmt.Printf("getStockHistory : %s error\n", v)
		} else {
			data[v] = day
		}
		curPercent := len(data) * 100 / len(symbols)
		if curPercent > percent {
			percent = curPercent
			fmt.Printf("download complete:%d%%\n", curPercent)
		}
		time.Sleep(time.Millisecond * 200)
	}
	fmt.Println("download finished")
	return data, nil
}

func (p *history) save(filepath string, data map[string]*daily) error {

	buffer, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, buffer, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *history) load(filepath string) (map[string]*daily, error) {

	if !isFileExist(filepath) {
		return nil, fmt.Errorf("load failed")
	}
	buffer, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	data := make(map[string]*daily, 10000)
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		v.maxPrice, v.maxPercent = calcMaxmaxPriceAndPercent(v.Close[0])
	}
	return data, nil
}

func (p *history) queryLastCloseDay(cookies []*http.Cookie, query int64) (int64, error) {

	remoteURL := "https://stock.xueqiu.com/v5/stock/chart/kline.json"
	values := url.Values{}
	values.Add("symbol", "SZ300202")
	values.Add("begin", strconv.FormatInt(query*1000, 10))
	values.Add("period", "day")
	values.Add("type", "before")
	values.Add("count", "-1")
	body, err := fetch(remoteURL, values, cookies)
	if err != nil {
		return 0, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return 0, err
	}

	data := m["data"].(map[string]interface{})
	item := data["item"].([]interface{})
	if len(item) != 1 {
		return 0, fmt.Errorf("data error")
	}

	array := item[0].([]interface{})
	miliSecond, _ := array[0].(float64)
	close := time.Unix((int64)(miliSecond/1000), 0)
	return time.Date(close.Year(), close.Month(), close.Day(), 0, 0, 0, 0, time.Local).Unix(), nil
}

func (p *history) unixSecToFilePath(unixSec int64) string {
	tm := time.Unix(unixSec, 0)
	return fmt.Sprintf("daily/%d-%d-%d.json", tm.Year(), tm.Month(), tm.Day())
}

func (p *history) update() {

	tmNow := time.Now().Local()
	query := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), 0, 0, 0, 0, time.Local).Unix()
	if tmNow.Hour() < 15 || (tmNow.Hour() == 15 && tmNow.Minute() < 1) {
		query -= 24 * 60 * 60
	}

	if query != p.query || p.due == 0 || p.count <= 0 || len(p.data) < p.count {

		resp, err := http.Get("https://xueqiu.com/")
		if err != nil {
			return
		}
		cookies := resp.Cookies()
		resp.Body.Close()

		if p.query != query {
			due, err := p.queryLastCloseDay(cookies, query)
			if err != nil {
				return
			}
			p.query = query
			if p.due != due {
				p.due = due
				p.count = 0
				p.mutex.Lock()
				for k := range p.data {
					delete(p.data, k)
				}
				p.mutex.Unlock()
			}
		}

		var data map[string]*daily
		filepath := p.unixSecToFilePath(p.due)

		data, err = p.load(filepath)
		if err == nil {
			p.count = len(data)
			p.mutex.Lock()
			p.data = data
			p.mutex.Unlock()
			return
		}

		symbols, err := p.getAllStockSymbol()
		if err != nil {
			return
		}
		p.count = len(symbols)

		downList := make([]string, 0, len(symbols))
		p.mutex.Lock()
		for k := range symbols {
			_, ok := p.data[k]
			if !ok {
				downList = append(downList, k)
			}
		}
		p.mutex.Unlock()

		data, err = p.download(query, downList, cookies)
		if err != nil {
			return
		}
		p.mutex.Lock()
		for k, v := range data {
			p.data[k] = v
		}
		if len(p.data) == p.count {
			p.save(filepath, p.data)
		}
		p.mutex.Unlock()
	}
}

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
