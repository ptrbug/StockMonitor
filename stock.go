package main

type stockHistory struct {
	close [5]float64
}

type stockToday struct {
	symbol             string
	percent            float64
	current            float64
	currentYearPercent float64
	name               string
}

type stock struct {
	stockToday
	stockHistory
}
