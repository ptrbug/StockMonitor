package main

import (
	"fmt"
	"time"
)

func main() {

	lp := newLimitUp()
	acc := newAccelerate()

	fmt.Print("start...\n")
	isWorking := false
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
			}
		}

		if request {
			if !isWorking {
				lp.reset()
				acc.reset()
				isWorking = true
				fmt.Print("working...\n")
			}
			curs, err := getTopPercent(topPercentCount)
			if err != nil {
				fmt.Printf("getTopPercnet(%d) failed\n", topPercentCount)
			}
			lp.update(curs)
			acc.update(curs)
		} else {
			if isWorking {
				fmt.Print("Waiting...\n")
			}
			isWorking = false
		}

		<-time.After(time.Second * 5)
	}

}
