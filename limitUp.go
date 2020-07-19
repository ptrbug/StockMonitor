package main

type stockSpeed struct {
	current
}

type limitUp struct {
	stocks map[string]current
}

func (p *limitUp) update(curs []current) {

}
