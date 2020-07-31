package main

import (
	"fmt"
	"time"
)

type stock struct {
	*realtime
	*daily
}

type limitUp struct {
	stocks map[string]*stock
	times  int
}

func newLimitUp() *limitUp {
	return &limitUp{stocks: make(map[string]*stock, 1000), times: 0}
}

func (p *limitUp) reset() {
	for k := range p.stocks {
		delete(p.stocks, k)
	}
	p.times = 0
}

func (p *limitUp) update(tmNow time.Time, history *history, reals []*realtime) {

	p.times++

	strTime := fmt.Sprintf("%d:%d:%d", tmNow.Hour(), tmNow.Minute(), tmNow.Second())

	for _, v := range reals {
		exist, ok := p.stocks[v.symbol]
		if ok {
			if exist.current == exist.maxPrice {
				if v.current < exist.maxPrice {
					fmt.Printf("%s %s %s 打开缺口 涨幅:%v 现价:%v \n", strTime, v.name, v.symbol, v.percent, v.current)
				}
			}
			v.flag = p.times
			p.stocks[v.symbol].realtime = v
		} else {

			var daily *daily
			history.mutex.Lock()
			daily, _ = history.data[v.symbol]
			history.mutex.Unlock()

			if daily != nil {
				s := &stock{realtime: v, daily: daily}
				v.flag = p.times
				p.stocks[v.symbol] = s
			}
		}
	}

	for k, v := range p.stocks {
		if v.flag != p.times {
			if v.current == v.maxPrice {
				fmt.Printf("%s %s %s 打开缺口\n", strTime, v.name, v.symbol)
			}
			delete(p.stocks, k)
		}
	}

}
