package main

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"time"
)

type playerstats struct {
	ImpactPerRnd     float64
	UserName         string
	SteamID          uint64
	Kills            int
	Deaths           int
	Assists          int
	HS               int
	HeadPercent      float64
	ADR              float64
	KAST             float64
	KDRatio          float64
	Firstkill        int
	FirstDeath       int
	FKDiff           int
	Round2k          int
	Round3k          int
	Round4k          int
	Round5k          int
	Totaldmg         int
	TradeKills       int
	TradeDeath       int
	CTkills          int
	Tkills           int
	EffectiveFlashes int
	AvgflshDuration  float64
	WeaponKill       map[int]int
	AvgDist          float64
	TotalDist        float64
	FlashesThrown    int
	ClanName         string
	TotalUtilDmg     int
	AvgKillsRnd      float64
	AvgDeathsRnd     float64
	AvgAssistsRnd    float64
	RoundSurvived    int
	RoundTraded      int
}

type roundKill struct {
	TimeOfKill     time.Duration
	Killer         string
	Victim         string
	KillerId       int64
	VictimId       int64
	Assistor       string
	IsHeadshot     bool
	VictFlashed    bool
	KillerFlashed  bool
	KillerWeapon   int
	KillerTeam     common.Team
	KillerDmgTkn   int
	VictimTeam     common.Team
	VictRel        bool
	positionKilled victPos
	KillerPos      killerPos
	Opening        bool
}

type victPos struct {
	X, Y, Z float64
}

type killerPos struct {
	X, Y, Z float64
}

func (p *parser) getActivePlayers(c []*common.Player) {
	for _, player := range c {
		steamId := player.SteamID64

		if p.Match.players == nil {
			p.Match.players = make(map[int64]playerstats)
		}

		if _, exists := p.Match.players[int64(steamId)]; !exists {
			return
		} else {
			p.Match.players[int64(steamId)] = p.ThePlayer(player)
		}
	}
}

func (p *parser) ThePlayer(player *common.Player) playerstats {
	return playerstats{
		ImpactPerRnd:     0,
		UserName:         player.Name,
		SteamID:          player.SteamID64,
		Kills:            0,
		Deaths:           0,
		Assists:          0,
		HS:               0,
		HeadPercent:      0,
		ADR:              0,
		KAST:             0,
		KDRatio:          0,
		Firstkill:        0,
		FirstDeath:       0,
		FKDiff:           0,
		Round2k:          0,
		Round3k:          0,
		Round4k:          0,
		Round5k:          0,
		Totaldmg:         0,
		TradeKills:       0,
		TradeDeath:       0,
		CTkills:          0,
		Tkills:           0,
		EffectiveFlashes: 0,
		AvgflshDuration:  0,
		AvgDist:          0,
		TotalDist:        0,
		FlashesThrown:    0,
		ClanName:         "",
		TotalUtilDmg:     0,
		AvgKillsRnd:      0,
		AvgDeathsRnd:     0,
		AvgAssistsRnd:    0,
		RoundSurvived:    0,
		RoundTraded:      0,
		WeaponKill:       p.makeweapons(),
	}
}

func (p *parser) makeweapons() map[int]int {
	return make(map[int]int)
}

func (p *parser) playergetter(e events.RoundEnd) {
	gs := p.parser.GameState()

	TeamH := gs.TeamTerrorists().Members()
	TeamJ := gs.TeamTerrorists().Members()

	p.statsetter(TeamH)
	p.statsetter(TeamJ)

}

func (p *parser) statsetter(c []*common.Player) {

	//gs := p.parser.GameState()

	for i := range c {
		steamId := c[i].SteamID64
		playerStat, exists := p.Match.players[int64(steamId)]

		if !exists {
			continue
		}

		playerStat.Kills = c[i].Kills()
		playerStat.Assists = c[i].Assists()
		playerStat.Deaths = c[i].Deaths()
		playerStat.Totaldmg = c[i].TotalDamage()
		playerStat.TotalUtilDmg = c[i].UtilityDamage()
	}
}

func (p *parser) killHandler(e events.Kill) {

	opening := false

	if e.Killer == nil || e.Victim == nil {
		return
	}

	if p.parser.GameState().IsWarmupPeriod() {
		p.state.warmupkill = append(p.state.warmupkill, e)
	}

	var assistorName string
	if e.Assister != nil {
		assistorName = e.Assister.Name
	}

	if e.Killer.ActiveWeapon() == nil {
		return
	}

	if p.state.round > 0 && p.state.round <= len(p.Match.Round) {
		if p.Match.Round[p.state.round-1].RoundKills == nil {
			p.Match.Round[p.state.round-1].RoundKills = make(map[int]roundKill)
		}
		count := len(p.Match.Round[p.state.round-1].RoundKills) + 1

		if count == 1 {
			opening = true
		} else {
			opening = false
		}

		if _, exists := p.Match.Round[p.state.round-1].RoundKills[count]; exists {
			return
		} else {

			VictKilAt := victPos{
				X: e.Victim.Position().X,
				Y: e.Victim.Position().Y,
			}

			KillerAt := killerPos{
				X: e.Killer.Position().X,
				Y: e.Killer.Position().Y,
			}
			p.Match.Round[p.state.round-1].RoundKills[count] = roundKill{
				TimeOfKill:     p.parser.CurrentTime(),
				Killer:         e.Killer.Name,
				Victim:         e.Victim.Name,
				Assistor:       assistorName,
				KillerId:       int64(e.Killer.SteamID64),
				VictimId:       int64(e.Victim.SteamID64),
				KillerTeam:     e.Killer.Team,
				VictimTeam:     e.Victim.Team,
				IsHeadshot:     e.IsHeadshot,
				VictFlashed:    e.Victim.IsBlinded(),
				KillerFlashed:  e.Killer.IsBlinded(),
				KillerDmgTkn:   100 - e.Killer.Health(),
				VictRel:        e.Victim.IsReloading,
				positionKilled: VictKilAt,
				KillerPos:      KillerAt,
				Opening:        opening,
			}
			count++
		}
	}

}
