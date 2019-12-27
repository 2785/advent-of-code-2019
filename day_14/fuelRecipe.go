package main

import (
	"fmt"
	"io/ioutil"
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
	fulfilled := false
	state := make(map[string]fulfillmentState)
	state["FUEL"] = fulfillmentState{required: 1, fulfilled: 0}
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
	fmt.Println("Ore required: ", state["ORE"].required)
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
