package main

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var (
	playerMap = make(map[string]int)
)

/*
	I forgot how this exactly works but in essence it gets
	the kills for the player had for that round, and put it into a playermap
	i think it has something to do with manipulation on the library level
	don't know how i even solved this.
*/

func (p *parser) GetPresRoundKill(e events.RoundFreezetimeEnd) {

	TeamOneMem := p.parser.GameState().TeamCounterTerrorists().Members()
	TeamTwoMem := p.parser.GameState().TeamTerrorists().Members()

	p.printThis(TeamOneMem)
	p.printThis(TeamTwoMem)

}

func (p *parser) printThis(c []*common.Player) {

	for _, v := range c {
		playerMap[v.Name] = v.Kills()
	}

}

// going to experiment with per round stats at a later date.
