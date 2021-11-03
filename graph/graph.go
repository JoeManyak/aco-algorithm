package graph

import (
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

func (w *WorldMap) ReinitializeMap() {
	//переініціалізація карти
}

func (a *Ant) FindNewWay() {
	start := rand.Intn(MaxCities)
	a.NowAt = start
	a.Length = 0
	a.Visited = []int{}
	for len(a.Visited) != MaxCities {
		var next int
		if a.Common {
			next = a.GetNext()
		} else {
			next = a.GetNextWild()
		}
		a.Length += a.World.DistMap[a.NowAt][next][0]
		a.World.PheromoneMap[]
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
	Cities    [MaxCities]City
	DistMap   [MaxCities][MaxCities][2]float64
	AntColony []Ant
	Lmin      float64
	PheromoneMap map[int]map[int][]Ant
}

type City struct {
	Id    int
	X     float64
	Y     float64
	Name  string
	World *WorldMap
}

const MaxCities = 4
const MaxCord = 4

const p = 0.7

func (w *WorldMap) RenewPheromone() {
	for i := 0; i < MaxCities-1; i++ {
		for j := i + 1; j < MaxCities; j++ {
			newPheromone := 0.0
			for _, v := range w.AntColony {
				newPheromone += w.Lmin / v.Length
			}
			w.DistMap[i][j][1] = (1-p)*w.DistMap[i][j][1] + newPheromone
		}
	}
}

func (w *WorldMap) CalculateGreedy() float64 {
	distance := 0.0
	cities := make([]int, 0)
	id := 0
	for len(cities) != MaxCities {
		cities = append(cities, id)
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
	}
	distance += w.DistMap[cities[0]][cities[len(cities)-1]][0]
	return distance
}

func WorldMapGenerate() WorldMap {
	rand.Seed(time.Now().UnixNano())
	w := WorldMap{[MaxCities]City{},
		[MaxCities][MaxCities][2]float64{}, []Ant{}, 0}
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
