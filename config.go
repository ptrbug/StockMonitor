package main

type cQueKou struct {
	On    bool
	MaxLB int
}

type cJiaShu struct {
	On    bool
	Risen [60]float64
}

type cConfig struct {
	QueKou cQueKou
	JiaShu cJiaShu
}

var config cConfig
