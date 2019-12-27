package intcodecomputer

import (
	"fmt"
	"strconv"
	"strings"
)

type intcodeComputerState struct {
	mem          []int64
	indexPointer int
	relativeBase int
}

func InitializeIntcodeComputer(numericIntcode []int64) intcodeComputerState {
	return intcodeComputerState{
		mem:          numericIntcode,
		indexPointer: 0,
		relativeBase: 0,
	}
}

func check(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
}

func len64(arr []int64) int64 {
	return int64(len(arr))
}

func ExecuteIntcode(
	startingState intcodeComputerState, getInput func(callCount int) (int, error),
) (
	output []int64, finalState intcodeComputerState, done bool, e error,
) {

	code := startingState.mem
	currIndex := startingState.indexPointer
	relativeBase := startingState.relativeBase
	inputCallCount := 0
	constructNewState := func() intcodeComputerState {
		return intcodeComputerState{
			mem:          code,
			indexPointer: currIndex,
			relativeBase: relativeBase,
		}
	}
	for currIndex >= 0 {
		numericOpcode := code[currIndex]

		// define util functions
		getByIndex := func(pos int) int64 {
			if pos >= len(code) {
				return 0
			}
			return code[pos]
		}

		getParamPositions := func(paramCount int) []int {
			out := make([]int, paramCount)
			for i := 0; i < paramCount; i++ {
				out[i] = int(getByIndex(currIndex + 1 + i))
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

			check(e)
			getParamPositions = func(paramCount int) []int {
				getSingleParamPosition := func(paramNumber int) int {
					if !(paramNumber >= 0 && paramNumber < len(poscode)) ||
						numericPosCode[len(poscode)-paramNumber-1] == 0 {
						return int(getByIndex(currIndex + paramNumber + 1))
					}
					c := numericPosCode[len(poscode)-paramNumber-1]
					if c == 1 {
						return currIndex + paramNumber + 1
					} else if c == 2 {
						return relativeBase + int(getByIndex(currIndex+paramNumber+1))
					} else {
						panic("Unknown parameter mode")
					}
				}
				out := make([]int, paramCount)
				for i := 0; i < paramCount; i++ {
					out[i] = getSingleParamPosition(i)
				}
				return out
			}

		}

		write := func(pos int, val int64) []int64 {
			for pos >= len(code) {
				code = append(code, 0)
			}
			code[pos] = val
			return code
		}

		switch numericOpcode {
		case 1:
			params := getParamPositions(3)
			p1, p2, p3 := params[0], params[1], params[2]
			code = write(p3, getByIndex(p1)+getByIndex(p2))
			currIndex += 4
		case 2:
			params := getParamPositions(3)
			p1, p2, p3 := params[0], params[1], params[2]
			code = write(p3, getByIndex(p1)*getByIndex(p2))
			currIndex += 4
		case 3:
			params := getParamPositions(1)
			p1 := params[0]
			in, e := getInput(inputCallCount)
			if e != nil {
				return output, constructNewState(), false, nil
			}
			code = write(p1, int64(in))
			inputCallCount++
			currIndex += 2
		case 4:
			params := getParamPositions(1)
			p1 := params[0]
			currIndex += 2
			output = append(output, getByIndex(p1))
		case 5:
			params := getParamPositions(2)
			p1, p2 := params[0], params[1]
			if getByIndex(p1) != 0 {
				currIndex = int(getByIndex(p2))
			} else {
				currIndex += 3
			}
		case 6:
			params := getParamPositions(2)
			p1, p2 := params[0], params[1]
			if getByIndex(p1) == 0 {
				currIndex = int(getByIndex(p2))
			} else {
				currIndex += 3
			}
		case 7:
			params := getParamPositions(3)
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
			params := getParamPositions(3)
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
			params := getParamPositions(1)
			p1 := params[0]
			relativeBase += int(getByIndex(p1))
			currIndex += 2
		default:
			currIndex = -1
		}
	}
	return output, constructNewState(), true, e
}
