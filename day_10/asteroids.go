package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strings"
)

type coordinate struct {
	x int
	y int
}

type asteroid struct {
	coordinate          coordinate
	lineOfSightCount    int
	orderOfVaporization int
}

func main() {
	// testMapStr := ".#..##.###...#######\n##.############..##.\n.#.######.########.#\n.###.#######.####.#.\n#####.##.#.##.###.##\n..#####..#.#########\n####################\n#.####....###.#.#.##\n##.#################\n#####.##.###..####..\n..######..##.#######\n####.##.####...##..#\n.#####..#.######.###\n##...#.##########...\n#.##########.#######\n.####.#.###.###.#.##\n....##.##.###..#####\n.#.#.###########.###\n#.#.#.#####.####.###\n###.##.####.##.#..##"

	// Read map file
	inputMapStr, e := ioutil.ReadFile("map.txt")
	if e != nil {
		panic(e)
	}
	// Parse map into coordinates
	parsedAsteroids := func(asteroidMap string) []asteroid {
		re := regexp.MustCompile(`#`)
		totalAsteroidsCount := len(re.FindAllString(asteroidMap, -1))
		asteroids := make([]asteroid, totalAsteroidsCount)
		rows := strings.Split(asteroidMap, "\n")
		asteroidIndex := 0
		for rowNumber, row := range rows {
			allAsteroids := re.FindAllStringIndex(row, -1)
			for _, bounds := range allAsteroids {
				asteroids[asteroidIndex] = asteroid{
					coordinate:          coordinate{x: bounds[0], y: rowNumber},
					lineOfSightCount:    0,
					orderOfVaporization: 0}
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
	station := parsedAsteroids[bestAsteroid]
	fmt.Println(station)

	// Vaporize 'em all
	countToBeVaporized := len(parsedAsteroids) - 1
	numberVaporized := 0
	vaporizeAsteroids := func() {
		losMap := make(map[float64]int)
		for i, target := range parsedAsteroids {
			if target.orderOfVaporization != 0 || i == bestAsteroid {
				continue
			}
			var dx, dy float64 = float64(target.coordinate.x - station.coordinate.x),
				float64(target.coordinate.y - station.coordinate.y)
			arclength := math.Sqrt(dx*dx + dy*dy)
			angle := math.Acos(dx / arclength)
			if dy < 0 {
				angle = 2*math.Pi - angle
			}
			// transform to make clockwise starting from top
			angle = angle + math.Pi/2
			if angle >= 2*math.Pi {
				angle = angle - 2*math.Pi
			}
			angle = round(angle, 4)
			losMap[angle] = i
		}
		keysArr := make([]float64, len(losMap))
		ind := 0
		for i := range losMap {
			keysArr[ind] = i
			ind++
		}
		sort.Float64s(keysArr)
		for _, v := range keysArr {
			parsedAsteroids[losMap[v]].orderOfVaporization = numberVaporized + 1
			numberVaporized++
		}
	}
	for numberVaporized < countToBeVaporized {
		vaporizeAsteroids()
	}
	for _, v := range parsedAsteroids {
		if v.orderOfVaporization == 200 {
			fmt.Print(v)
		}
	}

}

func round(num float64, precision int) float64 {
	multiple := math.Pow10(precision)
	return math.Round(num*multiple) / multiple
}
