package main

import (
	"fmt"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func main() {
	digitLookup := make(map[string]int, 10)
	for i := 0; i <= 9; i++ {
		digitLookup[strconv.Itoa(i)] = i
	}
	low, hi := 172851, 675869
	qualifyingNumbers := []int{}
	start := time.Now()
	for i := low; i <= hi; i++ {
		byteArr := []byte(strconv.Itoa(i))
		noDecrease := func(arr []byte) bool {
			for i := 1; i < len(byteArr); i++ {
				if byteArr[i] < byteArr[i-1] {
					return false
				}
			}
			return true
		}(byteArr)
		doubleNumber := func(arr []byte) bool {
			for i := 1; i < len(byteArr); i++ {
				if byteArr[i] == byteArr[i-1] {
					return true
				}
			}
			return false
		}(byteArr)
		if noDecrease && doubleNumber {
			qualifyingNumbers = append(qualifyingNumbers, i)
		}
	}

	additionalRule := []int{}
	for _, v := range qualifyingNumbers {
		byteArr := []byte(strconv.Itoa(v))
		checkTwoPosEqual := func(first int, second int) bool {
			inRange := func(i int) bool {
				return i >= 0 && i < len(byteArr)
			}
			if !(inRange(first) && inRange((second))) {
				return false
			} else {
				return byteArr[first] == byteArr[second]
			}
		}

		ruleCheck := func() bool {
			for i := 1; i < len(byteArr); i++ {
				if checkTwoPosEqual(i, i-1) && !checkTwoPosEqual(i, i+1) && !checkTwoPosEqual(i, i-2) {
					return true
				}
			}
			return false
		}()

		if ruleCheck {
			additionalRule = append(additionalRule, v)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Numbers of possible passwords: ", len(qualifyingNumbers))
	fmt.Printf("Computation done in %s", elapsed)
	fmt.Println("\nNumber of passwords satisfying additional Rule", len(additionalRule))
}
