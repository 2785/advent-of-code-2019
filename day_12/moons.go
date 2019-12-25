package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type coord struct {
	x, y, z int
}

type object struct {
	pos, vel    coord
	totalEnergy int
}

func main() {
	// Parse Input
	file, _ := ioutil.ReadFile("initialMoonPosition.txt")
	re := regexp.MustCompile(`[<>xyz=]`)
	fileStr := re.ReplaceAllString(string(file), "")
	fileStrArr := strings.Split(fileStr, "\n")
	moonCount := len(fileStrArr)
	moonsInitialPosition := make([]object, moonCount)
	for i, v := range fileStrArr {
		posStr := strings.Split(v, ",")
		point := make([]int, 3)
		for j, s := range posStr {
			s = strings.TrimSpace(s)
			num, _ := strconv.Atoi(s)
			point[j] = num
		}
		iniPos := coord{x: point[0], y: point[1], z: point[2]}
		iniVel := coord{x: 0, y: 0, z: 0}
		moonsInitialPosition[i] = object{pos: iniPos, vel: iniVel, totalEnergy: getTotalEnergy(iniPos, iniVel)}
	}

	// Calculate energy

	tMax := 1000
	states := moonsInitialPosition
	for t := 1; t <= tMax; t++ {
		newStates := make([]object, moonCount)
		for k, v := range states {
			newStates[k] = v
		}
		for i, currObject := range states {
			for j, targetObject := range states {
				if j == i {
					continue
				}
				if targetObject.pos.x > currObject.pos.x {
					newStates[i].vel.x++
				} else if targetObject.pos.x < currObject.pos.x {
					newStates[i].vel.x--
				}
				if targetObject.pos.y > currObject.pos.y {
					newStates[i].vel.y++
				} else if targetObject.pos.y < currObject.pos.y {
					newStates[i].vel.y--
				}
				if targetObject.pos.z > currObject.pos.z {
					newStates[i].vel.z++
				} else if targetObject.pos.z < currObject.pos.z {
					newStates[i].vel.z--
				}
			}
			newStates[i].pos.x = newStates[i].pos.x + newStates[i].vel.x
			newStates[i].pos.y = newStates[i].pos.y + newStates[i].vel.y
			newStates[i].pos.z = newStates[i].pos.z + newStates[i].vel.z
			newStates[i].totalEnergy = getTotalEnergy(newStates[i].pos, newStates[i].vel)
		}
		states = newStates
	}
	totalEnergy := 0
	for _, v := range states {
		totalEnergy += v.totalEnergy
	}
	fmt.Println(totalEnergy)
}

func getTotalEnergy(pos, vel coord) int {
	pe, ke := intAbs(pos.x)+intAbs(pos.y)+intAbs(pos.z), intAbs(vel.x)+intAbs(vel.y)+intAbs(vel.z)
	return pe * ke
}

func intAbs(num int) int {
	if num >= 0 {
		return num
	}
	return 0 - num
}
