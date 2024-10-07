package main

func (p *parser) calcADR(dmg int) float64 {
	roundsPlayed := p.parser.GameState().TotalRoundsPlayed()

	adr := float64(dmg) / float64(roundsPlayed)

	return adr
}

func (p *parser) calcKDRatio(kills, deaths int) float64 {
	KD := float64(kills) / float64(deaths)

	return KD
}

func (p *parser) calcHSPercent(kills, headshots int) float64 {
	return float64(headshots) / float64(kills)
}

func (p *parser) calcKPR(kills int) float64 {
	return float64(kills) / float64(p.parser.GameState().TotalRoundsPlayed())
}

func (p *parser) calcDPR(deaths int) float64 {
	return float64(deaths) / float64(p.parser.GameState().TotalRoundsPlayed())
}

func (p *parser) calcAPR(assists int) float64 {
	return float64(assists) / float64(p.parser.GameState().TotalRoundsPlayed())
}

func (p *parser) calcImpact(avgKil, avgAssists float64) float64 {
	return (2.13*avgKil + 0.42*avgAssists - 0.41)
}
