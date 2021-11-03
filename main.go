package main

import (
	"aco/graph"
	"fmt"
)

func main() {
	w := graph.WorldMapGenerate()
	w.AntColony = []graph.Ant{graph.Ant{
		NowAt:   0,
		Visited: []int{0},
		Length:  0,
		World:   &w,
	}}
	fmt.Println(w.AntColony[0].GetNext())
}
