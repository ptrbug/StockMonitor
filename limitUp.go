package main

import "fmt"

type limitUp struct {
	stocks map[string]current
}

func newLimitUp() *limitUp {
	return &limitUp{stocks: make(map[string]current, 1000)}
}

func (p *limitUp) update(curs []current) {

	for _, v := range curs {
		if v.percent > 9.90 {
			exist, ok := p.stocks[v.symbol]
			if ok {
				if v.percent < exist.percent {
					fmt.Printf("%s %s 打开缺口\n", v.name, v.symbol)
				}
			} else {
				p.stocks[v.symbol] = v
			}

		} else {
			_, ok := p.stocks[v.symbol]
			if ok {
				delete(p.stocks, v.symbol)
				fmt.Printf("%s %s 打开缺口\n", v.name, v.symbol)
			}
		}
	}

}
