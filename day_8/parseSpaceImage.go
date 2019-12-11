package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type layer struct {
	pixelData  map[int][]int
	valueCount map[int]int
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// read and process image into numeric format
	rawImg, e := ioutil.ReadFile("theImage.txt")
	checkError(e)
	imgStrArr := strings.Split(string(rawImg), "")
	imgData := make([]int, len(imgStrArr))
	for i, v := range imgStrArr {
		imgData[i], _ = strconv.Atoi(strings.TrimSpace(v))
	}
	// Process image data into layers, the layer number count thing is printed in the function parselayers
	layers, e := parseLayers(imgData, 25, 6)
	checkError(e)
	// Make the image to a string format
	renderedImage, e := makeImg(layers, map[int]string{0: " ", 1: "#", 2: " "}, 2)
	checkError(e)
	fmt.Println(renderedImage)
}

func parseLayers(imgData []int, width, height int) (out []layer, e error) {
	layerSize := width * height
	if len(imgData)%layerSize != 0 {
		panic("Input image does not have integer layers")
	}
	layerCount := len(imgData) / layerSize
	layers := make([]layer, layerCount)

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
	return layers, nil
}

func makeImg(layers []layer, pixelMap map[int]string, transparentVal int) (out string, e error) {
	if len(layers) == 0 {
		return "", errors.New("Empty image data received")
	}
	width, height, layerCount := len(layers[0].pixelData[0]), len(layers[0].pixelData), len(layers)
	renderedImg := ""
	toString := func(intArr []int) string {
		arr := make([]string, len(intArr))
		for i, v := range intArr {
			arr[i] = pixelMap[v]
		}
		return strings.Join(arr, "")
	}
	for i := 0; i < height; i++ {
		row := make([]int, width)
		for j := 0; j < width; j++ {
			val := transparentVal
			for k := 0; k < layerCount; k++ {
				if layers[k].pixelData[i][j] != transparentVal {
					val = layers[k].pixelData[i][j]
					break
				}
			}
			row[j] = val
		}
		renderedImg += toString(row) + "\n"
	}
	return renderedImg, nil
}
