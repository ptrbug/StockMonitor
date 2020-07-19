package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	resp, err := http.Get("https://xueqiu.com/")
	if err != nil {
		return
	}
	cookies := resp.Cookies()
	for _, v := range cookies {
		fmt.Println(v.Name)
		fmt.Println(v.Value)
	}

	lp := newLimitUp()
	acc := newAccelerate()

	needReset := true
	topPercentCount := 500
	for {
		tmNow := time.Now()
		weekday := tmNow.Local().Weekday()
		hour := tmNow.Local().Hour()
		minute := tmNow.Local().Minute()

		request := false
		if weekday >= 1 && weekday <= 5 {
			if (hour >= 9 && hour <= 10) || (hour == 11 && minute < 30) || (hour >= 13 && hour < 15) {
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
