package main

import (
	"fmt"
	"time"
)

func checkHistory(history *history) {

	for {
		history.update()
		time.Sleep(time.Second * 5)
	}
}

func main() {

	setAutoStart(true)

	history := newHisory()
	go checkHistory(history)

	lp := newLimitUp()
	acc := newAccelerate()

	fmt.Print("启动...\n")
	isLastOpening := false
	isFristUpdate := true
	topPercentCount := 500
	for {
		tmStart := time.Now()
		weekday := tmStart.Weekday()
		hour := tmStart.Hour()
		minute := tmStart.Minute()

		opening := false
		hm := hour*60 + minute
		if weekday >= 1 && weekday <= 5 {
			if (hm > (9*60+25) && hm < (11*60+30)) || (hm > (13*60-1) && hm < (15*60)) {
				opening = true
			}
		}

		if opening {
			if !isLastOpening {
				lp.reset()
				acc.reset()
				fmt.Print("开盘中...\n")
			}
			realtime, err := getTopPercent(topPercentCount)
			if err != nil {
				fmt.Printf("getTopPercnet(%d) failed\n", topPercentCount)
			}
			lp.update(tmStart, history, realtime)
			acc.update(tmStart, realtime)
		} else {
			if isLastOpening || isFristUpdate {
				fmt.Print("休盘中...\n")
			}
		}
		isLastOpening = opening
		isFristUpdate = false

		tmEnd := time.Now()
		duration := tmEnd.Sub(tmStart)
		if duration < time.Second {
			<-time.After(time.Second - duration)
		}

	}

}
