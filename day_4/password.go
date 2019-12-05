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
	elapsed := time.Since(start)
	fmt.Println("Numbers of possible passwords: ", len(qualifyingNumbers))
	fmt.Printf("Computation done in %s", elapsed)
}
