package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	checkDuplication := func(state []object) (x, y, z bool) {
		x, y, z = true, true, true
		for i, v := range moonsInitialPosition {
			if v.pos.x != state[i].pos.x || v.vel.x != state[i].vel.x {
				x = false
			}
			if v.pos.y != state[i].pos.y || v.vel.y != state[i].vel.y {
				y = false
			}
			if v.pos.z != state[i].pos.z || v.vel.z != state[i].vel.z {
				z = false
			}
		}
		return
	}

	tMax := 1000000
	logStep := 100000
	states := moonsInitialPosition
	start := time.Now()
	duplicationTimes := coord{x: 0, y: 0, z: 0}
	for t := 1; t <= tMax; t++ {
		if (t-1)%logStep == 0 {
			elapsed := time.Since(start)
			fmt.Println("; took", elapsed)
			start = time.Now()
			fmt.Print("Computing times t = ", t, " to ", t+logStep-1)
		}
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
		x, y, z := checkDuplication(newStates)
		if x && duplicationTimes.x == 0 {
			fmt.Println("\nDuplicated x found: t = ", t)
			duplicationTimes.x = t
		}
		if y && duplicationTimes.y == 0 {
			fmt.Println("\nDuplicated y found: t = ", t)
			duplicationTimes.y = t
		}
		if z && duplicationTimes.z == 0 {
			fmt.Println("\nDuplicated z found: t = ", t)
			duplicationTimes.z = t
		}
		if duplicationTimes.x != 0 && duplicationTimes.y != 0 && duplicationTimes.z != 0 {
			dup := lcm(duplicationTimes.x, duplicationTimes.y, duplicationTimes.z)
			fmt.Println("Duplication Steps = ", dup)
			return
		}
		states = newStates
	}
	fmt.Println("\nNo duplication found until t =", tMax)
	// totalEnergy := 0
	// for _, v := range states {
	// 	totalEnergy += v.totalEnergy
	// }
	// fmt.Println(totalEnergy)
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

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
