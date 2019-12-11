package main

import (
	"fmt"
	"io/ioutil"

	"strings"

	"strconv"
)

type layer struct {
	pixelData  map[int][]int
	valueCount map[int]int
}

func main() {
	// read and process image into numeric format
	rawImg, e := ioutil.ReadFile("theImage.txt")
	if e != nil {
		panic(e)
	}
	imgStrArr := strings.Split(string(rawImg), "")
	imgData := make([]int, len(imgStrArr))
	for i, v := range imgStrArr {
		imgData[i], _ = strconv.Atoi(strings.TrimSpace(v))
	}

	width, height := 25, 6
	layerSize := width * height
	if len(imgData)%layerSize != 0 {
		panic("Input image does not have integer layers")
	}

	layerCount := len(imgData) / layerSize
	layers := make(map[int]layer, layerCount)

	layerWithLowestZeroCount, minZeroCount := 0, layerSize

	for i := 0; i < layerCount; i++ {
		layerData := imgData[i*layerSize : (i+1)*layerSize]
		valueCount := make(map[int]int)
		for _, v := range layerData {
			_, ok := valueCount[v]
			if ok {
				valueCount[v]++
			} else {
				valueCount[v] = 1
			}
		}
		zeroCount, ok := valueCount[0]
		if !ok {
			fmt.Println("Layer ", i, " does not have a zero entry")
			zeroCount = 0
		}
		if zeroCount < minZeroCount {
			minZeroCount = zeroCount
			layerWithLowestZeroCount = i
		}
		pixelData := make(map[int][]int, height)
		for row := 0; row < height; row++ {
			pixelData[row] = layerData[row*width : (row+1)*width]
		}
		layers[i] = layer{pixelData: pixelData, valueCount: valueCount}
	}

	fmt.Println("Layer with the lowest zero count: ", layerWithLowestZeroCount,
		"\n#1 x #2 = ", layers[layerWithLowestZeroCount].valueCount[1]*layers[layerWithLowestZeroCount].valueCount[2])

	finalMsg := make([]int, layerSize)

	toString := func(intArr []int) string {
		arr := make([]string, len(intArr))
		for i, v := range intArr {
			switch v {
			case 0:
				arr[i] = " "
			case 1:
				arr[i] = "#"
			default:
				arr[i] = " "
			}
			// arr[i] = strconv.Itoa(v)
		}
		return strings.Join(arr, "")
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			val := 2
			for k := 0; k < layerCount; k++ {
				if layers[k].pixelData[i][j] != 2 {
					val = layers[k].pixelData[i][j]
					break
				}
			}
			finalMsg[i*width+j] = val
		}

		fmt.Println(toString(finalMsg[i*width : (i+1)*width]))
	}
	// fmt.Println(toString(finalMsg))
	// for i := 0; i < height; i++ {}
}
