package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type coordinate struct {
	x int
	y int
}

type intcodeComputerState struct {
	mem          []int64
	indexPointer int
	relativeBase int
}

func initializeIntcodeComputer(numericIntcode []int64) intcodeComputerState {
	return intcodeComputerState{
		mem:          numericIntcode,
		indexPointer: 0,
		relativeBase: 0,
	}
}

type direction int

func main() {
	strIntcode, e := ioutil.ReadFile("intcodeInput.txt")
	check(e)
	intcodeStrArr := strings.Split(string(strIntcode), ",")
	numIntcode := make([]int64, len(intcodeStrArr))
	for i, v := range intcodeStrArr {
		num, e := strconv.Atoi(v)
		check(e)
		numIntcode[i] = int64(num)
	}

	integerInput := func(num int) func(int) (int, error) {
		return func(callCount int) (int, error) {
			if callCount == 0 {
				return num, nil
			}
			return 0, errors.New("Called more than once")
		}
	}

	getNewDirection := func(turn int, currDir direction) direction {
		newDir := func() direction {
			switch turn {
			case 0:
				return currDir - 1
			case 1:
				return currDir + 1
			default:
				panic("Unknown direction")
			}
		}()
		if newDir == 4 {
			newDir = 0
		} else if newDir == -1 {
			newDir = 3
		}
		return newDir
	}

	currCoordinate := coordinate{x: 0, y: 0}
	var currDirection direction = 0
	path := make(map[coordinate]int)
	path[coordinate{x: 0, y: 0}] = 1 // Starting condition for part 2. Remove for part 1
	currentComputerState := initializeIntcodeComputer(numIntcode)

	for {
		color, ok := path[currCoordinate]
		if !ok {
			color = 0
		}
		out, newState, fin, e := executeIntcode(currentComputerState, integerInput(color))
		check(e)
		currentComputerState = newState
		if len(out) != 2 {
			panic("Did not get 2 outputs")
		}
		colorToPaint, turnCommand := out[0], out[1]
		path[currCoordinate] = int(colorToPaint)
		// fmt.Println(currCoordinate, colorToPaint)
		newDirection := getNewDirection(int(turnCommand), currDirection)
		currDirection = newDirection
		currCoordinate = func() coordinate {
			switch newDirection {
			case 0:
				return coordinate{currCoordinate.x, currCoordinate.y + 1}
			case 1:
				return coordinate{currCoordinate.x + 1, currCoordinate.y}
			case 2:
				return coordinate{currCoordinate.x, currCoordinate.y - 1}
			case 3:
				return coordinate{currCoordinate.x - 1, currCoordinate.y}
			default:
				panic("Unknown direction received")
			}
		}()
		if fin {
			break
		}
	}
	fmt.Println(len(path)) // output for part 1

	left, right, top, bottom := 0, 0, 0, 0
	for i := range path {
		if i.x < left {
			left = i.x
		}
		if i.x > right {
			right = i.x
		}
		if i.y > top {
			top = i.y
		}
		if i.y < bottom {
			bottom = i.y
		}
	}
	width, height := right-left+1, top-bottom+1
	frame := make([][]string, height)
	for i := range frame {
		frame[i] = make([]string, width)
		for j := 0; j < width; j++ {
			x, y := left+j, top-i
			point := coordinate{x: x, y: y}
			color, ok := path[point]
			if !ok {
				color = 0
			}
			frame[i][j] = func() string {
				switch color {
				case 0:
					return "."
				case 1:
					return "#"
				default:
					panic("Unknown color")
				}
			}()
		}
		fmt.Println(strings.Join(frame[i], ""))
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

func executeIntcode(
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
