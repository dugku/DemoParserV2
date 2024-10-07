package main

import (
	"encoding/json"
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
	Players map[int64]playerstats
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
	p.parser.RegisterEventHandler(p.RoundEcon)
	p.parser.RegisterEventHandler(p.Bombplanted)
	p.parser.RegisterEventHandler(p.playergetter)
	p.parser.RegisterEventHandler(p.killHandler)
	p.parser.RegisterEventHandler(p.ComplexRoundEndStuff)
	p.parser.RegisterEventHandler(p.GetPresRoundKill)

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

	outputFileName := fmt.Sprintf("%s-%s-%d.json", data.Match.TeamOne.Name, data.Match.TeamTwo.Name, 1)

	jsonData, err2 := json.MarshalIndent(data, "", " ")
	check(err2)

	err3 := os.WriteFile(outputFileName, jsonData, 0644)
	check(err3)

	printData(data)
}

func printData(data *parser) {
	// Iterate over the players map and print the SteamID and stats
	fmt.Println(len(data.Match.Players))
	for steamID, stats := range data.Match.Players {
		fmt.Printf("SteamID: %d\n", steamID)
		fmt.Printf("Player Stats: %+v\n", stats)
	}
}

func get_demos() {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
