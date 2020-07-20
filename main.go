package main

import (
	"fmt"
	"time"
)

func main() {

	lp := newLimitUp()
	acc := newAccelerate()

	fmt.Print("启动\n")
	needReset := true
	topPercentCount := 500
	for {
		tmNow := time.Now()
		weekday := tmNow.Local().Weekday()
		hour := tmNow.Local().Hour()
		minute := tmNow.Local().Minute()

		request := false
		hm := hour*60 + minute
		if weekday >= 1 && weekday <= 5 {
			if (hm > (9*60+25) && hm < (11*60+30)) || (hm > (13*60-1) && hm < (15*60)) {
				request = true
			} else {
				needReset = true
			}
		}

		request = true
		if request {
			if needReset {
				lp.reset()
				acc.reset()
				needReset = false
				fmt.Print("开盘采集中\n")
			}
			curs, err := getTopPercent(topPercentCount)
			if err != nil {
				fmt.Printf("getTopPercnet(%d) failed\n", topPercentCount)
			}
			lp.update(curs)
			acc.update(curs)
		}

		<-time.After(time.Second * 5)
	}

}
