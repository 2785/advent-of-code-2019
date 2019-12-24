package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("boost.txt")
	check(err)
	intcodeStr := strings.Split(string(file), ",")
	numIntcode := make([]int64, len(intcodeStr))
	for i, v := range intcodeStr {
		numIntcode[i], err = strconv.ParseInt(v, 10, 64)
		check(err)
	}

	integerInput := func(num int) func(int) (int, error) {
		return func(callCount int) (int, error) {
			return num, nil
		}
	}

	// Part 1
	// _, out, _, _, _ := executeIntcode(0, numIntcode, integerInput(1), true)
	// Part 2
	_, out, _, _, _ := executeIntcode(0, numIntcode, integerInput(2), true)
	fmt.Print(out)
}

func check(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
}

func len64(arr []int64) int64 {
	return int64(len(arr))
}

func executeIntcode(
	currIndex int64,
	code []int64,
	getInput func(callCount int) (int, error),
	inputCallOnce bool) (index int64, output []int64, mem []int64, done bool, e error) {
	inputCallCount := 0
	var relativeBase int64 = 0
	for code[currIndex] != 99 {
		numericOpcode := code[currIndex]
		getByIndex := func(pos int64) int64 {
			if pos >= int64(len(code)) {
				return 0
			}
			return code[pos]
		}
		getParams := func(paramCount int) []int64 {
			out := make([]int64, paramCount)
			for i := 0; i < paramCount; i++ {
				out[i] = getByIndex(currIndex + 1 + int64(i))
			}
			return out
		}

		if len(strconv.FormatInt(numericOpcode, 10)) > 1 {
			opcodeWithMode := strings.Split(strconv.FormatInt(code[currIndex], 10), "")
			opcode := opcodeWithMode[len(opcodeWithMode)-2:]
			poscode := opcodeWithMode[:len(opcodeWithMode)-2]
			numericPosCode := make([]int64, len(poscode))
			for i, v := range poscode {
				numericPosCode[i], _ = strconv.ParseInt(v, 10, 64)
			}
			numericOpcode, e = strconv.ParseInt(strings.Join(opcode, ""), 10, 64)
			getParamPosition := func(paramNumber int) int64 {
				if !(paramNumber >= 0 && paramNumber < len(poscode)) ||
					numericPosCode[len(poscode)-paramNumber-1] == 0 {
					return getByIndex(currIndex + int64(paramNumber) + 1)
				}
				c := numericPosCode[len(poscode)-paramNumber-1]
				if c == 1 {
					return currIndex + int64(paramNumber) + 1
				} else if c == 2 {
					return relativeBase + getByIndex(currIndex+int64(paramNumber)+1)
				} else {
					panic("Unknown parameter mode")
				}
			}
			if e != nil {
				return 0, nil, nil, false, e
			}

			getParams = func(paramCount int) []int64 {
				out := make([]int64, paramCount)
				for i := 0; i < paramCount; i++ {
					out[i] = getParamPosition(i)
				}
				return out
			}

		}

		write := func(pos, val int64) []int64 {
			if pos >= len64(code) {
				newCode := append(code, make([]int64, pos-len64(code)+1)...)
				newCode[len(newCode)-1] = val
				return newCode
			} else {
				code[pos] = val
				return code
			}

		}

		switch numericOpcode {
		case 1:
			params := getParams(3)
			p1, p2, p3 := params[0], params[1], params[2]
			code = write(p3, getByIndex(p1)+getByIndex(p2))
			currIndex += 4
		case 2:
			params := getParams(3)
			p1, p2, p3 := params[0], params[1], params[2]
			code = write(p3, getByIndex(p1)*getByIndex(p2))
			currIndex += 4
		case 3:
			params := getParams(1)
			p1 := params[0]
			in, e := getInput(inputCallCount)
			if e != nil {
				return currIndex, output, code, false, nil
			}
			code = write(p1, int64(in))
			inputCallCount++
			currIndex += 2
		case 4:
			params := getParams(1)
			p1 := params[0]
			currIndex += 2
			output = append(output, getByIndex(p1))
		case 5:
			params := getParams(2)
			p1, p2 := params[0], params[1]
			if getByIndex(p1) != 0 {
				currIndex = getByIndex(p2)
			} else {
				currIndex += 3
			}
		case 6:
			params := getParams(2)
			p1, p2 := params[0], params[1]
			if getByIndex(p1) == 0 {
				currIndex = getByIndex(p2)
			} else {
				currIndex += 3
			}
		case 7:
			params := getParams(3)
			p1, p2, p3 := params[0], params[1], params[2]
			val := func() int64 {
				if getByIndex(p1) < getByIndex(p2) {
					return 1
				}
				return 0
			}()
			code = write(p3, val)
			currIndex += 4
		case 8:
			params := getParams(3)
			p1, p2, p3 := params[0], params[1], params[2]
			val := func() int64 {
				if getByIndex(p1) == getByIndex(p2) {
					return 1
				}
				return 0
			}()
			code = write(p3, val)
			currIndex += 4
		case 9:
			params := getParams(1)
			p1 := params[0]
			relativeBase += getByIndex(p1)
			currIndex += 2
		}
	}
	if code[currIndex] == 99 {
		return 0, output, nil, true, e
	}

	return currIndex, output, code, false, e
}
