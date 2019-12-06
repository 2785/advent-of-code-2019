package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(getFuel(100756, 0))
	file, err := ioutil.ReadFile("rocket_parts.txt")
	check(err)
	partsMass := strings.Split(string(file), "\n")
	fmt.Print(partsMass)
	partsFuelRequirement := make([]int, len(partsMass))
	for i, v := range partsMass {
		massInt, e := strconv.Atoi(strings.TrimSpace(v))
		check(e)
		partsFuelRequirement[i] = getFuel(massInt, 0)
	}
	totalFuelRequirement := 0
	for _, v := range partsFuelRequirement {
		totalFuelRequirement += v
	}

	fmt.Println("\nTotal fuel requirement = ", totalFuelRequirement)
}

func check(e error) {
	if e != nil {
		fmt.Println("\nError encountered: ", e.Error())
	}
}

func getFuel(remainingMass, totalFuel int) (totalFuelRequirement int) {
	fuelNeeded := remainingMass/3 - 2
	if fuelNeeded >= 1 {
		return getFuel(fuelNeeded, totalFuel+fuelNeeded)
	} else {
		return totalFuel
	}
}
