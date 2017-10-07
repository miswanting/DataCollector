package main

import (
	"fmt"
	"strings"
)

// var matrix [][]string

type DataPair struct {
	X string
	Y string
}
type DataPiece struct {
	Data []DataPair
}

type DataList struct {
	Data []DataPiece
}

var dataList DataList

// Collect ...
func Collect(cache Cache) [][]string {
	dissolve(cache)
	return arrange()
}
func dissolve(cache Cache) {
	for i := 0; i < len(cache.Data); i++ {
		var dataPiece DataPiece
		switch {
		// F-7000 FL Spectrophotometer
		case string(cache.Data[i].Data)[0:6] == "Sample":
			var dataPair DataPair
			tmpName := strings.Split(cache.Data[i].Name, "\\")
			tmpNameLen := len(tmpName)
			tmpNameName := tmpName[tmpNameLen-1]
			tmpName = strings.Split(tmpNameName, ".")
			tmpNameLen = len(tmpName)
			tmpNameName = strings.Join(tmpName[0:tmpNameLen-1], ".")
			dataPair.X, dataPair.Y = "Length", tmpNameName
			dataPiece.Data = append(dataPiece.Data, dataPair)
			lines := strings.Split(string(cache.Data[i].Data), "\r\n")
			var flag = false
			for j := 0; j < len(lines); j++ {
				if flag {
					if lines[j] == "" {
						flag = false
					} else {
						tmp := strings.Split(lines[j], "\t")
						dataPair.X, dataPair.Y = tmp[0], tmp[1]
						dataPiece.Data = append(dataPiece.Data, dataPair)
					}
				} else {
					if lines[j] == "nm\tData" {
						flag = true
					}
				}
			}
		}
		dataList.Data = append(dataList.Data, dataPiece)
	}
}
func arrange() [][]string {
	fmt.Println(&dataList.Data[0].Data[0].X)
	fmt.Println(&dataList.Data[0].Data[1].X)
	x := len(dataList.Data) + 1
	y := len(dataList.Data[0].Data)
	// var matrix [][]string
	matrix := make([][]string, y)
	for i := range matrix {
		matrix[i] = make([]string, x)
	}
	for i := 0; i < y; i++ {
		matrix[i][0] = dataList.Data[0].Data[i].X
	}
	for i := 0; i < x-1; i++ {
		for j := 0; j < y; j++ {
			matrix[j][i+1] = dataList.Data[i].Data[j].Y
		}
	}
	return matrix
}
