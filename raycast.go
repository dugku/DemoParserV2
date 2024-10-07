package main

import (
	"encoding/json"
	"math"
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
	jsonfile, err := os.ReadFile(file)
	check(err)

	e := json.Unmarshal(jsonfile, &data)
	check(e)

	return data
}

func raycast(victx, victy float64, edges []vector, name string) (bool, string) {
	count := 0
	where := ""
	tolerance := 1e-7

	for i := 0; i < len(edges); i++ {
		curr := edges[i]
		next := edges[(i+1)%len(edges)]

		if curr.Y == next.Y {
			continue
		}

		if victy < math.Min(curr.Y, next.Y) || victy > math.Max(curr.Y, next.Y) {
			continue
		}

		xInterept := (victy-curr.Y)*(next.X-curr.X)/(victy-curr.Y) + curr.X

		if victx < xInterept+tolerance {
			if where == "" {
				where = name
			}
			count++
		}

	}

	return count%2 == 1, where
}
