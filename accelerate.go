package main

type speed struct {
	current
	index    int
	count    int
	percents [60]float64
}

type accelerate struct {
	stocks map[string]*speed
	times  int
}

func newAccelerate() *accelerate {
	return &accelerate{stocks: make(map[string]*speed, 1000), times: 0}
}

func (p *accelerate) reset() {
	for k := range p.stocks {
		delete(p.stocks, k)
	}
	p.times = 0
}

func (p *accelerate) update(curs []*current) {

	p.times++
	for _, v := range curs {
		exist, ok := p.stocks[v.symbol]
		if ok {
			v.flag = p.times
			exist.current = *v
		} else {
			v.flag = p.times
			p.stocks[v.symbol] = &speed{current: *v}
		}
	}

	for k, v := range p.stocks {
		if v.flag != p.times {
			delete(p.stocks, k)
		}
	}
}
