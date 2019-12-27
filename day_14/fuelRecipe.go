package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type recipeEntry struct {
	reactants map[string]int
	name      string
	prodCount int
}

type fulfillmentState struct {
	required, fulfilled int
}

func main() {
	// Read file
	file, e := ioutil.ReadFile("recipe.txt")
	check(e)
	lines := strings.Split(string(file), "\n")

	// Parse into recipe
	recipes := make(map[string]recipeEntry)
	for _, v := range lines {
		entry := lineToRecipeEntry(v)
		recipes[entry.name] = entry
	}

	// Get 'em fuel
	getOreForFuel := func(fuel int) int {
		fulfilled := false
		state := make(map[string]fulfillmentState)
		state["FUEL"] = fulfillmentState{required: fuel, fulfilled: 0}
		for !fulfilled {
			for i, v := range state {
				if v.required > v.fulfilled {
					if i == "ORE" {
						continue
					}
					need := v.required - v.fulfilled
					reactionCount := need / recipes[i].prodCount
					if need%recipes[i].prodCount != 0 {
						reactionCount++
					}
					for j, r := range recipes[i].reactants {
						val, ok := state[j]
						if !ok {
							val = fulfillmentState{required: 0, fulfilled: 0}
						}
						val.required += r * reactionCount
						state[j] = val
					}
					state[i] = fulfillmentState{required: v.required, fulfilled: v.fulfilled + reactionCount*recipes[i].prodCount}
				}
			}
			fulfilled = func() bool {
				out := true
				for i, v := range state {
					if i == "ORE" {
						continue
					}
					if v.required > v.fulfilled {
						out = false
						return out
					}
				}
				return out
			}()
		}
		return state["ORE"].required
	}

	fuelPerOre := getOreForFuel(1)

	oreCapacity := 1e12
	fuel := int(math.Round(oreCapacity / float64(fuelPerOre)))
	fmt.Println(fuel)
	ore := getOreForFuel(fuel)
	lo, hi := fuel, 2*fuel
	for float64(getOreForFuel(hi)) < oreCapacity {
		hi += lo
	}
	for hi-lo >= 100 {
		mid := (hi + lo) / 2
		if float64(getOreForFuel(mid)) < oreCapacity {
			lo = mid
		} else {
			hi = mid
		}
	}
	for i := 0; i < 110; i++ {
		ore = getOreForFuel(lo + 1)
		if float64(ore) > oreCapacity {
			break
		}
		lo++
	}
	fmt.Println("Max Fuel: ", lo)
}

func lineToRecipeEntry(line string) recipeEntry {
	parts := strings.Split(line, "=>")
	if len(parts) != 2 {
		panic("????")
	}
	r, p := strings.Split(parts[0], ","), parts[1]
	r = stringSliceMap(r, func(s string) string { return strings.TrimSpace(s) })
	reactants := make(map[string]int)
	for _, v := range r {
		split := strings.Split(v, " ")
		count, e := strconv.Atoi(split[0])
		check(e)
		reactants[split[1]] = count
	}
	split := strings.Split(strings.TrimSpace(p), " ")
	count, e := strconv.Atoi(split[0])
	check(e)
	return recipeEntry{
		reactants: reactants,
		prodCount: count,
		name:      split[1],
	}
}

func stringSliceMap(arr []string, f func(string) string) []string {
	out := make([]string, len(arr))
	for i, v := range arr {
		out[i] = f(v)
	}
	return out
}
