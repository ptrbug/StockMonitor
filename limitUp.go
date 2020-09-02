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

	for _, v := range reals {
		exist, ok := p.stocks[v.symbol]
		if ok {
			if exist.current == exist.limitUpPrice {
				if v.current < exist.limitUpPrice {

					msg := fmt.Sprintf("%s %s %s 打开缺口 涨幅:%.2f 现价:%0.2f 连续涨停:%d\n",
						timeToString(tmNow), v.name, v.symbol, v.percent, v.current, exist.limitUpContinueCount)

					msgLog.Printf("%s", msg)
					fmt.Printf("%s", msg)
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
				if config.QueKou.MaxLB < 0 || daily.limitUpContinueCount <= config.QueKou.MaxLB {
					s := &stock{realtime: v, daily: daily}
					v.flag = p.times
					p.stocks[v.symbol] = s
				}
			}
		}
	}

	for k, v := range p.stocks {
		if v.flag != p.times {
			if v.current == v.limitUpPrice {
				msg := fmt.Sprintf("%s %s %s 打开缺口 连续涨停:%d\n",
					timeToString(tmNow), v.name, v.symbol, v.limitUpContinueCount)

				msgLog.Printf("%s", msg)
				fmt.Printf("%s", msg)
			}
			delete(p.stocks, k)
		}
	}

}
