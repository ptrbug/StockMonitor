package main

import (
	"fmt"
	"time"
)

const maxPrecentRecord = 60

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

func (p *speed) isSpeedUp(percent float64) (bool, int, float64) {
	sz := p.size()
	for i := 0; i < sz; i++ {
		lastPercent := p.prev(-i)
		diff := percent - lastPercent
		if diff >= config.JiaShu.Risen[i] {
			return true, i + 1, diff
		}
	}
	return false, 0, 0
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
			isUp, dt, value := exist.isSpeedUp(v.percent)
			if isUp {
				fmt.Printf("%s %s %s 加速上涨 涨幅:%v 现价:%v (%d秒涨%.2f)\n",
					timeToString(tmNow), v.name, v.symbol, v.percent, v.current, dt, value)
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
