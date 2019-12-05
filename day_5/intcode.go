package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("intcodeInput.txt")
	check(err)
	intcodeStr := strings.Split(string(file), ",")
	numIntcode := make([]int, len(intcodeStr))
	for i, v := range intcodeStr {
		numIntcode[i], err = strconv.Atoi(v)
		check(err)
	}
	executeIntcode(numIntcode)
}

func check(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
}

func executeIntcode(code []int) []int {
	currIndex := 0
	var e error
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
				} else {
					return true
				}
			}
			check(e)

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
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Please enter the input number")
			in, e := reader.ReadString('\n')
			check(e)
			val, e := strconv.Atoi(strings.TrimSpace(in))
			check(e)
			code[p1] = val
			currIndex += 2
		case 4:
			p1 := code[currIndex+1]
			if ifParamPosMode(0) {
				p1 = code[p1]
			}
			fmt.Println(p1)
			currIndex += 2
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

	return code
}
