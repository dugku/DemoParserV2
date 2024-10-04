package main

import (
	"fmt"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var JsonCoords string
var found bool
var posData PositionData
var positionData PositionData

type RoundInfo struct {
	RoundNum       int
	TeamA          string
	TeamB          string
	EconA          int
	EconB          int
	ScoreA         int
	ScoreB         int
	TypeofbuyA     string
	TypeofbuyB     string
	SurvivorsA     []string
	SurvivorsB     []string
	BombPlanted    bool
	PlayerPlanted  string
	RoundEndReason string
	Sidewon        string
	RoundKills     map[int]roundKill
}

func (p *parser) stateController(e events.RoundStart) {

	p.state.RoundonGoing = true
	p.state.round++

	round := RoundInfo{}
	//gotta append those rounds
	p.Match.Round = append(p.Match.Round, round)
}

func (p *parser) MatchStartHandler(e events.MatchStartedChanged) {
	gs := p.parser.GameState()

	/*
		Have to make this part of the program concurrent, part of the reason is because,
		since the demoinfocs lib is concurrent and is parsing all of these events at once.
		Since that is just how to game works and so fourth. Plus part of the reason why I made
		a new parser just for this is because I already tried to add concurrency to the other one
		but everything went to shit for about 5 hours... So yeah.
	*/

	wg.Add(1)
	go func() {
		defer wg.Done()
		if e.NewIsStarted {
			ActivePlayers := gs.Participants().Playing()

			p.getActivePlayers(ActivePlayers)

			coordDir := "C:\\Users\\Mike\\Desktop\\DemoParserV2\\mapCoords"

			p.Match.TeamOne.Side = common.TeamCounterTerrorists
			p.Match.TeamTwo.Side = common.TeamTerrorists

			p.Match.Map = p.parser.Header().MapName

			paths, err := os.ReadDir(coordDir)
			check(err)

			found = false

			for _, entry := range paths {
				maName := strings.Split(entry.Name(), ".")

				if maName[0] == p.Match.Map {
					JsonCoords = filepath.Join(coordDir, entry.Name())
					found = true
				}
			}
		}
		if found {
			posData = jsonLoader(JsonCoords, positionData)
		} else {
			fmt.Println("We are fucked bro")
		}
	}()
	wg.Wait()
}

func (p *parser) TeamSideSwitch(e events.TeamSideSwitch) {
	p.Match.TeamOne.Side = common.TeamTerrorists
	p.Match.TeamTwo.Side = common.TeamCounterTerrorists
}

func (p *parser) ScoreUpdater(e events.ScoreUpdated) {
	TeamAScore, TeamAName := p.checkSide(p.Match.TeamOne.Side)
	TeamBScore, TeamBName := p.checkSide(p.Match.TeamTwo.Side)

	if p.state.round > 0 && p.state.round <= len(p.Match.Round) {
		p.Match.Round[p.state.round-1].RoundNum = p.state.round
		p.Match.Round[p.state.round-1].ScoreA = TeamAScore
		p.Match.Round[p.state.round-1].ScoreB = TeamBScore
		p.Match.Round[p.state.round-1].TeamA = TeamAName
		p.Match.Round[p.state.round-1].TeamB = TeamBName
	}
}

func (p *parser) checkSide(team common.Team) (TeamScore int, TeamName string) {

	gs := p.parser.GameState()

	if team == common.TeamCounterTerrorists {
		return gs.TeamCounterTerrorists().Score(), gs.TeamCounterTerrorists().ClanName()
	}

	if team == common.TeamTerrorists {
		return gs.TeamTerrorists().Score(), gs.TeamTerrorists().ClanName()
	}

	return 0, " "
}

func (p *parser) RoundEcon(e events.RoundFreezetimeEnd) {
	TeamAEcon, TeamABuy := p.CheckEcon(p.Match.TeamOne.Side)
	TeamBEcon, TeamBBuy := p.CheckEcon(p.Match.TeamTwo.Side)

	if p.state.round > 0 && p.state.round <= len(p.Match.Round) {
		p.Match.Round[p.state.round-1].EconA = TeamAEcon
		p.Match.Round[p.state.round-1].EconB = TeamBEcon
		p.Match.Round[p.state.round-1].TypeofbuyA = TeamABuy
		p.Match.Round[p.state.round-1].TypeofbuyB = TeamBBuy
	}
}

func (p *parser) CheckEcon(team common.Team) (Econ int, TypeBuy string) {
	gs := p.parser.GameState()

	FullBuy := 20000
	HalfBuy := 10000
	SemiEco := 5000

	if team == common.TeamCounterTerrorists {
		equipmentVal := gs.TeamCounterTerrorists().CurrentEquipmentValue()

		return equipmentVal, assessBuyType(equipmentVal, FullBuy, HalfBuy, SemiEco)
	}

	if team == common.TeamTerrorists {
		equipmentVal := gs.TeamTerrorists().CurrentEquipmentValue()

		return equipmentVal, assessBuyType(equipmentVal, FullBuy, HalfBuy, SemiEco)
	}

	return 0, " "
}

func assessBuyType(Value, Full, Half, SemiEco int) string {
	switch {
	case Value >= Full:
		return "Fullbuy"
	case Value >= Half && Value < Full:
		return "Halfbuy"
	case Value >= SemiEco && Value < Half:
		return "ForceBuy"
	default:
		return "Eco"
	}
}

// set bombplanted = true and get player name.
func (p *parser) Bombplanted(e events.BombPlanted) {
	if p.state.round > 0 && p.state.round <= len(p.Match.Round) {
		roundInfo := &p.Match.Round[p.state.round-1]

		roundInfo.BombPlanted = true
		roundInfo.PlayerPlanted = e.Player.Name
	}
}

func (p *parser) ComplexRoundEndStuff(e events.RoundEnd) {
	ReasonsMap := map[int]string{
		1: "TargetBombed",
		7: "BombDefused",
		8: "CTWin",
		9: "TWin",
	}

	WinnerMap := map[int]string{
		2: "Terrorists",
		3: "Counter Terrorists",
	}

	Reason := e.Reason
	SideWon := e.Winner

	if p.state.round > 0 && p.state.round <= len(p.Match.Round) {
		roundInfo := &p.Match.Round[p.state.round-1]

		roundInfo.RoundEndReason = ReasonsMap[int(Reason)]
		roundInfo.Sidewon = WinnerMap[int(SideWon)]

		p.PlayersSurvived()
		p.TradeLogic()
		p.SideKillLogic()
		p.PositionGetter()
	}
}

func (p *parser) PlayersSurvived() {
	TeamAPlayers := p.CheckTeamSur(p.Match.TeamOne.Side)
	TeamBPlayers := p.CheckTeamSur(p.Match.TeamTwo.Side)

	for _, v := range TeamAPlayers {
		roundInfo := &p.Match.Round[p.state.round-1]
		if v.IsAlive() {
			playerId := v.SteamID64
			roundInfo.SurvivorsA = append(roundInfo.SurvivorsA, v.String())

			playerStat, exists := p.Match.players[int64(playerId)]

			if !exists {
				return
			}
			playerStat.RoundSurvived++
			p.Match.players[int64(playerId)] = playerStat
		}
	}

	for _, v := range TeamBPlayers {
		roundInfo := &p.Match.Round[p.state.round-1]
		if v.IsAlive() {
			playerId := v.SteamID64
			roundInfo.SurvivorsB = append(roundInfo.SurvivorsB, v.String())

			playerStat, exists := p.Match.players[int64(playerId)]
			if !exists {
				return
			}

			playerStat.RoundSurvived++

			p.Match.players[int64(playerId)] = playerStat
		}
	}
}

func (p *parser) CheckTeamSur(team common.Team) []*common.Player {
	gs := p.parser.GameState()

	if team == common.TeamCounterTerrorists {
		return gs.TeamCounterTerrorists().Members()
	}
	if team == common.TeamTerrorists {
		return gs.TeamTerrorists().Members()
	}
	return nil
}
func (p *parser) TradeLogic() {

	roundInfo := &p.Match.Round[p.state.round-1]

	for key, _ := range roundInfo.RoundKills {

		if key+1 < len(roundInfo.RoundKills) {
			nextValue := roundInfo.RoundKills[key+1]

			if roundInfo.RoundKills[key].Killer == nextValue.Victim && ((nextValue.TimeOfKill - roundInfo.RoundKills[key].TimeOfKill) < (5*time.Second) {
				TradeKillId := nextValue.KillerId
				TradeVictId := nextValue.VictimId

				playerStat, exists := p.Match.players[int64(TradeKillId)]

				if !exists {
					continue
				}

				playerStat.TradeKills++

				p.Match.players[int64(TradeKillId)] = playerStat

				VictStat, exists := p.Match.players[int64(TradeVictId)]

				if !exists {
					return
				}

				VictStat.TradeKills++

				p.Match.players[int64(TradeVictId)] = VictStat
			}
		}
	}
}
func (p *parser) PositionGetter() {}
