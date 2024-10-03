package main

import (
	"fmt"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"os"
)

type parser struct {
	parser demoinfocs.Parser
	state  parsingState
	Match  *MatchInfo
}

type parsingState struct {
	round        int
	RoundonGoing bool
	warmupkill   []events.Kill
}

type MatchInfo struct {
	Map     string
	TeamOne TeamA
	TeamTwo TeamB
	Round   []RoundInfo
	players map[int64]playerstats
}

type TeamA struct {
	Name string
	Side common.Team
}

type TeamB struct {
	Name string
	Side common.Team
}

func (p *parser) startParsing(fileDEM string) error {
	f, e := os.Open(fileDEM)
	check(e)
	defer f.Close()

	p.parser = demoinfocs.NewParser(f)
	defer p.parser.Close()

	p.Match = &MatchInfo{}

	p.parser.RegisterEventHandler(p.stateController)
	p.parser.RegisterEventHandler(p.MatchStartHandler)
	p.parser.RegisterEventHandler(p.TeamSideSwitch)
	p.parser.RegisterEventHandler(p.ScoreUpdater)

	e = p.parser.ParseToEnd()
	check(e)

	return nil
}

func main() {
	demodir := "C:\\Users\\iphon\\Desktop\\DemoParseV2\\gun5-vs-passion-ua-m1-dust2.dem"

	//have to do this or else it won't work because this is how golang works
	data := &parser{Match: &MatchInfo{Round: make([]RoundInfo, 0)}}

	err := data.startParsing(demodir)
	check(err)

	printData(data)
}

func printData(data *parser) {
	fmt.Println("here")
	for i, round := range data.Match.Round {
		fmt.Printf("Round #%d - %+v\n", i, round)
	}
}

func get_demos() {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
