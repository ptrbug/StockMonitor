package main

import (
	"fmt"
	"time"
)

const maxPrecentRecord = 12

var risen [12]float64 = [12]float64{2.0, 2.2, 2.4, 2.6, 2.8, 3.0, 3.2, 3.4, 3.6, 3.8, 3.9, 4.0}

type speed struct {
	index    int
	count    int
	percents [maxPrecentRecord]float64
}

func (p *speed) push(percent float64) {
	p.count++
	p.index++
	if p.index >= maxPrecentRecord {
		p.index = p.index % maxPrecentRecord
	}
	p.percents[p.index] = percent
}

func (p *speed) size() int {
	if p.count > maxPrecentRecord {
		return maxPrecentRecord
	}
	return p.count
}

func (p *speed) prev(offset int) float64 {
	index := p.index + offset
	if index < 0 {
		index = (index % maxPrecentRecord) + maxPrecentRecord
	}
	return p.percents[index]
}

func (p *speed) isSpeedUp(percent float64) bool {
	sz := p.size()
	for i := 0; i < maxPrecentRecord; i++ {
		if sz > i {
			lastPercent := p.prev(-i)
			diff := percent - lastPercent
			if diff > risen[i] {
				return true
			}
		}
	}
	return false
}

type currentspeed struct {
	*realtime
	speed
}

type accelerate struct {
	reals map[string]*currentspeed
	times int
}

func newAccelerate() *accelerate {
	return &accelerate{reals: make(map[string]*currentspeed, 1000), times: 0}
}

func (p *accelerate) reset() {
	for k := range p.reals {
		delete(p.reals, k)
	}
	p.times = 0
}

func (p *accelerate) update(tmNow time.Time, reals []*realtime) {

	p.times++
	for _, v := range reals {
		exist, ok := p.reals[v.symbol]
		if ok {
			v.flag = p.times
			exist.realtime = v
			if exist.isSpeedUp(v.percent) {
				fmt.Printf("%s %s %s 加速上涨 涨幅:%v 现价:%v \n", timeToString(tmNow), v.name, v.symbol, v.percent, v.current)
			}
			exist.push(v.percent)

		} else {
			v.flag = p.times
			cur := &currentspeed{realtime: v}
			cur.push(v.percent)
			p.reals[v.symbol] = cur
		}
	}

	for k, v := range p.reals {
		if v.flag != p.times {
			delete(p.reals, k)
		}
	}
}
