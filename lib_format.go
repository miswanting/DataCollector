package main

import (
	"strings"
)

var dbVersion = "v171007[2]"

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
		switch {
		// F-7000 FL Spectrophotometer
		case string(cache.Data[i].Data)[0:6] == "Sample":
			fl(cache, i)
		// UV
		case string(cache.Data[i].Data)[0:4] == "nm;A":
			uv(cache, i)
		}
	}
}
func arrange() [][]string {
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

// 数据库

func stdDP(cache Cache, i int, prefix string, midfix string, suffix string) { // 标准处理流程
	var dataPiece DataPiece
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
			if lines[j] == suffix { // 需设定(可选)：数据后一排
				flag = false
			} else {
				tmp := strings.Split(lines[j], midfix) // 需设定：数据间分隔符
				dataPair.X, dataPair.Y = tmp[0], tmp[1]
				dataPiece.Data = append(dataPiece.Data, dataPair)
			}
		} else {
			if lines[j] == prefix { // 需设定：数据前一排
				flag = true
			}
		}
	}
	dataList.Data = append(dataList.Data, dataPiece)
}

func fl(cache Cache, i int) { // 荧光（适用标准处理流程）
	stdDP(cache, i, "nm\tData", "\t", "")
}
func uv(cache Cache, i int) { // 紫外（适用标准处理流程）
	stdDP(cache, i, ";0.00 sec", ";", "")
}
