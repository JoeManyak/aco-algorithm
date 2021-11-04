package graph

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const alpha = 2
const beta = 4

type Ant struct {
	NowAt   int
	Visited []int
	Length  float64
	World   *WorldMap
	Common  bool
}

func (w *WorldMap) Solve(common, wild, time int) {
	w.AntColony = []*Ant{}
	for i := 0; i < common; i++ {
		w.AntColony = append(w.AntColony, &Ant{
			World:  w,
			Common: true,
		})
	}
	for i := 0; i < wild; i++ {
		w.AntColony = append(w.AntColony, &Ant{
			World:  w,
			Common: false,
		})
	}
	for i := 0; i < time; i++ {
		for _, ant := range w.AntColony {
			ant.FindNewWay()
		}
		w.RenewPheromone()
		fmt.Println("i:", i, ":", w.Lmin, w.Shortest())
	}
	fmt.Println(w.Lmin, w.Shortest())
}

func (w *WorldMap) Shortest() float64 {
	shortest := 0.0
	shortestId := -1
	for i, v := range w.AntColony {
		if shortestId == -1 || v.Length < shortest {
			shortestId = i
			shortest = v.Length
		}
	}
	return shortest
}

func (a *Ant) FindNewWay() {
	start := rand.Intn(MaxCities)
	a.NowAt = start
	a.Length = 0
	a.Visited = []int{start}
	for len(a.Visited) != MaxCities {
		var next int
		if a.Common {
			next = a.GetNext()
		} else {
			next = a.GetNextWild()
		}
		a.Length += a.World.DistMap[a.NowAt][next][0]
		if a.World.PheromoneMap[a.NowAt] == nil {
			a.World.PheromoneMap[a.NowAt] = map[int][]*Ant{}
		}
		if a.World.PheromoneMap[a.NowAt][next] == nil {
			a.World.PheromoneMap[a.NowAt][next] = []*Ant{}
		}
		a.World.PheromoneMap[a.NowAt][next] = append(a.World.PheromoneMap[a.NowAt][next], a)
		if a.World.PheromoneMap[next] == nil {
			a.World.PheromoneMap[next] = map[int][]*Ant{}
		}
		a.World.PheromoneMap[next][a.NowAt] = a.World.PheromoneMap[a.NowAt][next]
		a.Visited = append(a.Visited, next)
		a.NowAt = next
	}
	a.Length += a.World.DistMap[a.NowAt][start][0]
}

func (a *Ant) GetNextWild() int {
	newList := make([]int, 0, 0)
	for _, v := range a.World.Cities {
		b := false
		for _, v2 := range a.Visited {
			if v2 == v.Id {
				b = true
			}
		}
		if b {
			continue
		}
		newList = append(newList, v.Id)
	}
	random := rand.Intn(len(newList))
	return newList[random]
}

func (a *Ant) GetNext() int {
	CityChance := [MaxCities]float64{}
	WholeChances := 0.0
	for i := 0; i < MaxCities; i++ {
		skip := false
		for _, v := range a.Visited {
			if v == i {
				skip = true
			}
		}
		if skip {
			CityChance[i] = 0
			continue
		}
		CityChance[i] = math.Pow(a.World.DistMap[a.NowAt][i][1], alpha) *
			math.Pow(1/(a.World.DistMap[a.NowAt][i][0]), beta)
		WholeChances += CityChance[i]
	}
	randomized := rand.Float64() * WholeChances
	for i, v := range CityChance {
		randomized -= v
		if randomized <= 0 {
			return i
		}
	}
	return MaxCities - 1
}

type WorldMap struct {
	Cities       [MaxCities]City
	DistMap      [MaxCities][MaxCities][2]float64
	AntColony    []*Ant
	Lmin         float64
	PheromoneMap map[int]map[int][]*Ant
}

type City struct {
	Id    int
	X     float64
	Y     float64
	Name  string
	World *WorldMap
}

const MaxCities = 200
const MaxCord = 40

const p = 0.7

func (w *WorldMap) RenewPheromone() {
	for i := 0; i < MaxCities-1; i++ {
		for j := i + 1; j < MaxCities; j++ {
			newPheromone := 0.0
			if w.PheromoneMap[i][j] != nil {
				for _, v := range w.PheromoneMap[i][j] {
					newPheromone += w.Lmin / v.Length
				}
			}
			w.DistMap[i][j][1] = (1-p)*w.DistMap[i][j][1] + newPheromone
			w.DistMap[j][i][1] = w.DistMap[i][j][1]
		}
	}
	w.PheromoneMap = map[int]map[int][]*Ant{}
}

func (w *WorldMap) CalculateGreedy() float64 {
	distance := 0.0
	cities := []int{0}
	id := 0
	for len(cities) != MaxCities {
		nextId := -1
		minVal := 0.0
		for i := 0; i < MaxCities; i++ {
			skip := false
			for _, v := range cities {
				if i == v {
					skip = true
				}
			}
			if skip {
				continue
			}
			if nextId == -1 || w.DistMap[id][i][0] < minVal {
				nextId = i
				minVal = w.DistMap[id][i][0]
			}
		}
		id = nextId
		distance += minVal
		cities = append(cities, id)
	}
	distance += w.DistMap[cities[0]][cities[len(cities)-1]][0]
	return distance
}

func WorldMapGenerate() WorldMap {
	rand.Seed(time.Now().UnixNano())
	w := WorldMap{
		Cities:       [MaxCities]City{},
		DistMap:      [MaxCities][MaxCities][2]float64{},
		AntColony:    []*Ant{},
		Lmin:         0,
		PheromoneMap: map[int]map[int][]*Ant{},
	}
	for i := 0; i < MaxCities; i++ {
		w.Cities[i] = w.cityGenerate(i, "City"+strconv.Itoa(i))
	}
	for i := 0; i < MaxCities-1; i++ {
		for j := i + 1; j < MaxCities; j++ {
			dist := w.Cities[i].CalculateDistance(w.Cities[j])
			w.DistMap[i][j] = [2]float64{dist, 1}
			w.DistMap[j][i] = [2]float64{dist, 1}
		}
	}
	return w
}

func (c City) CalculateDistance(to City) float64 {
	x := c.X - to.X
	y := c.Y - to.Y
	return math.Sqrt(x*x + y*y)
}

func (w *WorldMap) cityGenerate(id int, name string) City {
	return City{
		Id:    id,
		X:     generateCord(),
		Y:     generateCord(),
		Name:  name,
		World: w,
	}
}

func generateCord() float64 {
	return float64(rand.Intn(MaxCord*100)) / 100
}
