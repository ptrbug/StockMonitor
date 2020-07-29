package main

import "fmt"

type stock struct {
	*current
	*daily
	maxPercent float64
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

func (p *limitUp) update(curs []*current) {

	p.times++

	for _, v := range curs {
		exist, ok := p.stocks[v.symbol]
		if ok {
			if exist.percent == exist.maxPercent {
				if v.percent < exist.percent {
					fmt.Printf("%s %s 打开缺口 涨幅:%v 现价:%v \n", v.name, v.symbol, v.percent, v.current)
				}
			}
			v.flag = p.times
			p.stocks[v.symbol].current = v
		} else {
			var day *daily
			var maxPercent float64
			s := &stock{current: v, daily: day, maxPercent: maxPercent}
			v.flag = p.times
			p.stocks[v.symbol] = s
		}

	}

	for k, v := range p.stocks {
		if v.flag != p.times {
			if v.percent > 9.90 {
				fmt.Printf("%s %s 打开缺口\n", v.name, v.symbol)
			}
			delete(p.stocks, k)
		}
	}

}
