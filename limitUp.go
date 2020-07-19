package main

import "fmt"

type limitUp struct {
	stocks map[string]*current
	times  int
}

func newLimitUp() *limitUp {
	return &limitUp{stocks: make(map[string]*current, 1000), times: 0}
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
			exist.times++
			if exist.current > 9.90 {
				if v.current < exist.current {
					delete(p.stocks, v.symbol)
					fmt.Printf("%s %s 打开缺口\n", v.name, v.symbol)
				}
			}
		} else {
			v.times = p.times
			p.stocks[v.symbol] = v
		}
	}

	for k, v := range p.stocks {
		if v.times != p.times {
			if v.current > 9.90 {
				fmt.Printf("%s %s 打开缺口\n", v.name, v.symbol)
			}
			delete(p.stocks, k)
		}
	}

}
