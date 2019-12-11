package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type amplifierState struct {
	mem       []int
	currIndex int
	done      bool
}

func main() {
	file, err := ioutil.ReadFile("intcodeInput.txt")
	check(err)
	intcodeStr := strings.Split(string(file), ",")
	numIntcode := make([]int, len(intcodeStr))
	for i, v := range intcodeStr {
		numIntcode[i], err = strconv.Atoi(v)
		check(err)
	}

	amplifierInput := func(phase int, input int) func(int) (int, error) {
		return func(callCount int) (int, error) {
			if callCount == 0 {
				return phase, nil
			} else if callCount == 1 {
				return input, nil
			}
			return 0, errors.New("Input called more than 2 times")
		}
	}

	integerInput := func(num int) func(int) (int, error) {
		return func(callCount int) (int, error) {
			return num, nil
		}
	}

	// Part 1

	maxOutput, maxOutputPhase := 0, []int{}
	for _, seq := range permutations([]int{0, 1, 2, 3, 4}) {
		currVal := 0
		for _, phase := range seq {
			dupeIntcode := make([]int, len(numIntcode))
			for i, v := range numIntcode {
				dupeIntcode[i] = v
			}
			_, o, _, _, _ := executeIntcode(0, dupeIntcode, amplifierInput(phase, currVal))

			currVal = o
		}
		if currVal > maxOutput {
			maxOutput, maxOutputPhase = currVal, seq
		}
	}
	fmt.Print("Maximum Engine Output: ", maxOutput, "\nPhase Arrangement: ", maxOutputPhase)

	// Part 2

	maxOutput, maxOutputPhase = 0, []int{}
	for _, seq := range permutations([]int{5, 6, 7, 8, 9}) {
		endVal := 0
		currVal := 0
		amplifiers := make([]amplifierState, 5)
		for i, v := range seq {
			dupeIntcode := make([]int, len(numIntcode))
			for i, v := range numIntcode {
				dupeIntcode[i] = v
			}
			ampIndex, o, mem, _, _ := executeIntcode(0, dupeIntcode, amplifierInput(v, currVal))
			amplifiers[i] = amplifierState{mem: mem, currIndex: ampIndex, done: false}
			currVal = o
		}
		endVal = currVal
		finished := false

		for !finished {
			for i, v := range amplifiers {
				ampIndex, o, mem, done, e := executeIntcode(v.currIndex, v.mem, integerInput(currVal))
				if e != nil {
					panic(e)
				}
				if done {
					amplifiers[i].done = true
					finished = true
					break
				}
				currVal = o
				if i == len(amplifiers)-1 {
					endVal = currVal
				}
				amplifiers[i].mem = mem
				amplifiers[i].currIndex = ampIndex
			}
		}
		if endVal > maxOutput {
			maxOutput, maxOutputPhase = endVal, seq
		}
	}
	fmt.Print("Maximum Engine Output: ", maxOutput, "\nPhase Arrangement: ", maxOutputPhase)
}

func check(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
}

func executeIntcode(currIndex int, code []int, getInput func(callCount int) (int, error)) (index int, out int, mem []int, done bool, e error) {
	inputCallCount := 0
	for code[currIndex] != 99 {
		numericOpcode := code[currIndex]
		getTwoParam := func() (int, int) {
			return code[code[currIndex+1]], code[code[currIndex+2]]
		}
		ifParamPosMode := func(index int) bool { return true }
		if len(strconv.Itoa(numericOpcode)) > 1 {
			opcodeWithMode := strings.Split(strconv.Itoa(code[currIndex]), "")
			opcode := opcodeWithMode[len(opcodeWithMode)-2:]
			poscode := opcodeWithMode[:len(opcodeWithMode)-2]
			numericPosCode := make([]int, len(poscode))
			for i, v := range poscode {
				numericPosCode[i], _ = strconv.Atoi(v)
			}
			numericOpcode, e = strconv.Atoi(strings.Join(opcode, ""))
			ifParamPosMode = func(index int) bool {
				if index >= 0 && index < len(poscode) {
					return numericPosCode[len(poscode)-index-1] == 0
				}
				return true

			}
			if e != nil {
				return 0, 0, nil, false, e
			}

			getTwoParam = func() (int, int) {
				p1 := code[currIndex+1]
				if ifParamPosMode(0) {
					p1 = code[p1]
				}
				p2 := code[currIndex+2]
				if ifParamPosMode(1) {
					p2 = code[p2]
				}
				return p1, p2
			}
		}

		switch numericOpcode {
		case 1:
			p1, p2 := getTwoParam()
			code[code[currIndex+3]] = p1 + p2
			currIndex += 4
		case 2:
			p1, p2 := getTwoParam()
			code[code[currIndex+3]] = p1 * p2
			currIndex += 4
		case 3:
			p1 := code[currIndex+1]
			code[p1], e = getInput(inputCallCount)
			if e != nil {
				panic(e)
			}
			inputCallCount++
			currIndex += 2
		case 4:
			p1 := code[currIndex+1]
			if ifParamPosMode(0) {
				p1 = code[p1]
			}
			// fmt.Println(p1)
			currIndex += 2
			return currIndex, p1, code, false, e
		case 5:
			p1, p2 := getTwoParam()
			if p1 != 0 {
				currIndex = p2
			} else {
				currIndex += 3
			}
		case 6:
			p1, p2 := getTwoParam()
			if p1 == 0 {
				currIndex = p2
			} else {
				currIndex += 3
			}
		case 7:
			p1, p2 := getTwoParam()
			code[code[currIndex+3]] = func() int {
				if p1 < p2 {
					return 1
				}
				return 0
			}()
			currIndex += 4
		case 8:
			p1, p2 := getTwoParam()
			code[code[currIndex+3]] = func() int {
				if p1 == p2 {
					return 1
				}
				return 0
			}()
			currIndex += 4
		}
	}

	return 0, 0, nil, true, e
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
