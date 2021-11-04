package main

import (
	"aco/graph"
)

func main() {
	w := graph.WorldMapGenerate()
	w.Lmin = w.CalculateGreedy()
	w.Solve(30, 15, 100)
}
