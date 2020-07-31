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

	history := newHisory()
	go checkHistory(history)

	lp := newLimitUp()
	acc := newAccelerate()

	fmt.Print("start...\n")
	isWaitForOpen := true
	topPercentCount := 500
	for {
		tmNow := time.Now()
		weekday := tmNow.Weekday()
		hour := tmNow.Hour()
		minute := tmNow.Minute()

		opening := false
		hm := hour*60 + minute
		if weekday >= 1 && weekday <= 5 {
			if (hm > (9*60+25) && hm < (11*60+30)) || (hm > (13*60-1) && hm < (15*60)) {
				opening = true
			}
		}

		if opening {
			if isWaitForOpen {
				lp.reset()
				acc.reset()
				isWaitForOpen = false
				fmt.Print("working...\n")
			}
			realtime, err := getTopPercent(topPercentCount)
			if err != nil {
				fmt.Printf("getTopPercnet(%d) failed\n", topPercentCount)
			}
			lp.update(tmNow, history, realtime)
			acc.update(realtime)
		} else {
			if !isWaitForOpen {
				fmt.Print("Waiting...\n")
			}
			isWaitForOpen = true
		}

		<-time.After(time.Second * 5)
	}

}
