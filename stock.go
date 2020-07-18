package main

type stockHistory struct {
	close [5]float32
}

type stockToday struct {
	symbol             string
	percent            float32
	current            float32
	currentYearPercent float32
	name               string
}

type stock struct {
	stockToday
	stockHistory
}
