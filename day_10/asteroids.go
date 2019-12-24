package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"
)

type Coordinate struct {
	x int
	y int
}

type Asteroid struct {
	coordinate       Coordinate
	lineOfSightCount int
}

func main() {
	// testMapStr := "#.#...#.#.\n.###....#.\n.#....#...\n##.#.#.#.#\n....#.#.#.\n.##..###.#\n..#...##..\n..##....##\n......#...\n.####.###."

	// Read map file
	inputMapStr, e := ioutil.ReadFile("map.txt")
	if e != nil {
		panic(e)
	}
	// Parse map into coordinates
	parsedAsteroids := func(asteroidMap string) []Asteroid {
		re := regexp.MustCompile(`#`)
		totalAsteroidsCount := len(re.FindAllString(asteroidMap, -1))
		asteroids := make([]Asteroid, totalAsteroidsCount)
		rows := strings.Split(asteroidMap, "\n")
		asteroidIndex := 0
		for rowNumber, row := range rows {
			allAsteroids := re.FindAllStringIndex(row, -1)
			for _, bounds := range allAsteroids {
				asteroids[asteroidIndex] = Asteroid{
					coordinate:       Coordinate{x: bounds[0], y: rowNumber},
					lineOfSightCount: 0}
				asteroidIndex++
			}
		}
		return asteroids
	}(string(inputMapStr))
	// fmt.Print(parsedAsteroids)

	maxLosCount, bestAsteroid := 0, 0

	// Run through each coordinate, see how many asteroids it can see
	for i, curr := range parsedAsteroids {
		losMap := make(map[float64]struct{})
		for j, target := range parsedAsteroids {
			if j == i {
				continue
			}
			var dx, dy float64 = float64(target.coordinate.x - curr.coordinate.x),
				float64(target.coordinate.y - curr.coordinate.y)
			arclength := math.Sqrt(dx*dx + dy*dy)
			cosine := math.Acos(dx / arclength)
			if dy < 0 {
				cosine = cosine + math.Pi
			}
			cosine = round(cosine, 4)
			_, ok := losMap[cosine]
			if !ok {
				losMap[cosine] = struct{}{}
			}
		}
		parsedAsteroids[i].lineOfSightCount = len(losMap)
		if len(losMap) > maxLosCount {
			bestAsteroid = i
			maxLosCount = len(losMap)
		}
	}
	fmt.Print(parsedAsteroids[bestAsteroid])
}

func round(num float64, precision int) float64 {
	multiple := math.Pow10(precision)
	return math.Round(num*multiple) / multiple
}
