package main

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
