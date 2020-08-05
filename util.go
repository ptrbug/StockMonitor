package main

import (
	"fmt"
	"math"
	"time"
)

func isSymbolMatch(symbol string, name string) bool {
	if len(symbol) < 5 {
		return false
	}
	if symbol[2:5] == "688" || symbol[2:5] == "787" || symbol[2:5] == "789" {
		return false
	}
	if symbol[2:5] == "171" {
		return false
	}
	if name[0:1] == "N" {
		return false
	}
	return true
}

func calcLimitUpPriceAndlimitUpPercent(lastClosePrice float64) (float64, float64) {
	n10 := math.Pow10(2)
	limitUpPrice := math.Trunc((lastClosePrice*1.1+0.5/n10)*n10) / n10
	limitUpPercent := math.Trunc((limitUpPrice*100/lastClosePrice+0.5/n10)*n10) / n10
	return limitUpPrice, limitUpPercent
}

func timeToString(tm time.Time) string {
	return fmt.Sprintf("%d:%d:%d", tm.Hour(), tm.Minute(), tm.Second())
}

func calcLimitUpContinueCount(close []float64) int {
	count := 0
	if len(close) >= 2 {
		for i := 0; i <= len(close)-2; i++ {
			limitUpPrice, _ := calcLimitUpPriceAndlimitUpPercent(close[i+1])
			if close[i] != limitUpPrice {
				break
			}
			count++
		}
	}
	return count
}
