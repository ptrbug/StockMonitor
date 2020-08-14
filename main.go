package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

func checkHistory(history *history) {

	for {
		history.update()
		time.Sleep(time.Second * 5)
	}
}

func main() {

	//ch dir
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	//set restart
	setAutoStart(true)

	//load config
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("config.json reading error", err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Unmarshal config.json file error", err)
	}

	//check config
	isValidRisen := true
	for i := 0; i < len(config.JiaShu.Risen); i++ {
		if config.JiaShu.Risen[i] <= 0 {
			isValidRisen = false
			break
		}
		for j := i + 1; j < len(config.JiaShu.Risen); j++ {
			if config.JiaShu.Risen[i] > config.JiaShu.Risen[j] {
				isValidRisen = false
				break
			}
		}
		if !isValidRisen {
			break
		}
	}

	if !isValidRisen {
		fmt.Printf("config.json { \"JiaShu\" : {\"Risen\"} }配置无效,使用默认配置\n")
		for i := 0; i < len(config.JiaShu.Risen); i++ {
			config.JiaShu.Risen[i] = 1 + math.Sqrt(float64(i+1))/math.Sqrt(float64(maxPrecentRecord))*3
			fmt.Printf("%0.2f ", config.JiaShu.Risen[i])
			if i%10 == 9 {
				fmt.Printf("\n")
			}
		}
	}

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
				fmt.Print("开盘中,等待开盘...\n")
			}
			realtime, err := getTopPercent(topPercentCount)
			if err != nil {
				fmt.Printf("getTopPercnet(%d) failed\n", topPercentCount)
			}
			if config.QueKou.On {
				lp.update(tmStart, history, realtime)
			}

			if config.JiaShu.On {
				acc.update(tmStart, realtime)
			}

		} else {
			if isLastOpening || isFristUpdate {
				fmt.Print("休盘中,等待开盘...\n")
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
