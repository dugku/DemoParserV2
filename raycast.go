package main

import (
	"encoding/json"
	"os"
)

type vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Position struct {
	Name   string   `json:"name"`
	Points []vector `json:"points"`
}

type PositionData struct {
	Positions []Position `json:"positions"`
}

func jsonLoader(file string, data PositionData) PositionData {
	jsonfile, err := os.Open(file)
	check(err)

	e := json.Unmarshal(jsonfile, &data)
	check(e)

	return data
}
