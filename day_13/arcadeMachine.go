package main

import (
	"advent-of-code-2019/intcodecomputer"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"strings"

	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
	return
}

type point struct {
	x, y, id int
}

func main() {
	file, e := ioutil.ReadFile("intcodeInput.txt")
	check(e)
	codeStr := strings.Split(string(file), ",")
	numIntcode := make([]int64, len(codeStr))
	for i, v := range codeStr {
		num, e := strconv.Atoi(v)
		check(e)
		numIntcode[i] = int64(num)
	}
	initialComputerState := intcodecomputer.InitializeIntcodeComputer(numIntcode)
	nilInput := func(callCount int) (int, error) {
		return 0, errors.New("No input allowed")
	}

	out, currState, _, _ := intcodecomputer.ExecuteIntcode(initialComputerState, nilInput)
	// fmt.Println(len(out))
	readPoints := func(output []int64) (points []point, left, right, top, bottom int, ball, pad point, blockCount, score int) {
		points = make([]point, len(out)/3)
		left, right, top, bottom, blockCount, score = 0, 0, 0, 0, 0, 0
		for i := range points {
			x, y, id := int(output[i*3]), int(output[i*3+1]), int(output[i*3+2])
			if x == -1 && y == 0 {
				score = id
				continue
			}
			if id == 4 {
				ball = point{x: x, y: y, id: id}
			}
			if id == 3 {
				pad = point{x: x, y: y, id: id}
			}
			if id == 2 {
				blockCount++
			}
			if x < left {
				left = x
			}
			if x > right {
				right = x
			}
			if y > top {
				top = y
			}
			if y < bottom {
				bottom = y
			}
			points[i] = point{x: x, y: y, id: id}
		}
		return
	}
	currBoard, left, right, top, bottom, currBall, currPad, blockCount, currScore := readPoints(out)
	height, width := top-bottom+1, right-left+1

	getDirection := func(ball, pad point) int {
		out := 0
		if ball.x > pad.x {
			out = 1
		}
		if ball.x < pad.x {
			out = -1
		}
		return out
	}
	// emptyRow := make([]string, width)

	pointsToFrame := func(points []point) []string {
		frame := make([]string, height*width)
		for i := range frame {
			frame[i] = " "
		}

		for _, point := range points {
			symbol := func() string {
				switch point.id {
				case 0:
					return " "
				case 1:
					return "X"
				case 2:
					return "\u2592"
				case 3:
					return "_"
				case 4:
					return "\u25ef"
				default:
					panic("unknown tile")
				}
			}()
			frame[width*(point.y-bottom)+point.x-left] = symbol
		}
		return frame
	}

	updateFrame := func(frame []string, points []point) []string {
		for _, point := range points {
			symbol := func() string {
				switch point.id {
				case 0:
					return " "
				case 1:
					return "X"
				case 2:
					return "\u2592"
				case 3:
					return "_"
				case 4:
					return "\u25ef"
				default:
					panic("unknown tile")
				}
			}()
			frame[width*(point.y-bottom)+point.x-left] = symbol
		}
		return frame
	}

	printBoard := func(frame []string) {
		for i := 0; i < height; i++ {
			fmt.Println(strings.Join(frame[i*width:(i+1)*width], " "))
		}
		return
	}

	currFrame := pointsToFrame(currBoard)

	printBoard(currFrame)

	for blockCount > 0 {
		newScore := currScore
		input := func(callCount int) (int, error) {
			if callCount > 0 {
				return 0, errors.New("No Duplicated input allowed")
			}
			return getDirection(currBall, currPad), nil
		}
		out, currState, _, _ = intcodecomputer.ExecuteIntcode(currState, input)
		currBoard, _, _, _, _, currBall, currPad, _, newScore = readPoints(out)
		if newScore > currScore {
			currScore = newScore
		}
		currFrame = updateFrame(currFrame, currBoard)
		blockCount = func() int {
			count := 0
			for _, v := range currFrame {
				if v == "\u2592" {
					count++
				}
			}
			return count
		}()
		printBoard(currFrame)
		fmt.Println("Current Score: ", currScore)
		time.Sleep(200 * time.Millisecond)
	}

}
