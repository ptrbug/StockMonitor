package main

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


func (p *limitUp) update(history *history, curs []*current) {

	p.times++

	for _, v := range curs {
		exist, ok := p.stocks[v.symbol]
		if ok {
			if exist.percent == exist.maxPercent {
				if v.percent < exist.percent {

				}
			}
			v.flag = p.times
			p.stocks[v.symbol].current = v
		} else {

			var daily *daily
			history.mutex.Lock()
			daily, _ = history.data[v.symbol]
			history.mutex.Unlock()

			if daily != nil {
				if daily.Close[0] > 0 {
					var maxPercent float64
					s := &stock{current: v, daily: daily, maxPercent: maxPercent}
					v.flag = p.times
					p.stocks[v.symbol] = s
				}
			}
		}
	}

	for k, v := range p.stocks {
		if v.flag != p.times {
			if v.percent == v.maxPercent {

			}
			delete(p.stocks, k)
		}
	}

}
