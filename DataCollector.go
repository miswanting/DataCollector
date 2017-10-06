package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var doc = `实验数据处理程序

*作者：
  何雨航 <453542772@qq.com>
*版本：
  v1.0.0-go
*介绍：
  该脚本可以将分散于多个文件的实验数据处理成方便导入到OriginLab等各种表格程序的格式。
*使用方法：
  1. 依次将实验数据文件按顺序，拖动到本程序文件上，程序将自动导入实验数据并存入名为"cache.json"的缓存文件中（请不要手动删除它）；
  2. 实验数据文件全部被程序读取之后，直接启动本程序。程序将自动清除缓存文件"cache.json"并将整理好的数据输出到程序旁边的名为："output.txt"的文件中；
  3. 打开并用“Ctrl+A”全选"output.txt"中的数据并粘贴到表格中。
*注意事项：
  1. 本程序支持的实验数据文件格式有限，无法适用于所有试验数据文件，若需要作者支持更多的格式，请联系作者；
  2. 由于不同版本的操作系统的参数机制不同，推荐将数据文件按照顺序一个一个地拖动到本程序文件上进行读取，而不要同时拖动多个，否则会有概率导致读取顺序不理想；
  3. 由于本程序会在旁边生成名为"cache.json"、"output.txt"文件，所以请不要把其他有价值的文件命名为这个并和本程序放在一起；
  4. 拖动操作是指类似于把文件框选并用鼠标左键按住拖动到文件夹的方式，把实验数据文件框选并用鼠标左键按住并拖动到本程序文件上；
  5. 本程序的功能于Windows 10上测试通过，其他版本并未测试；
  6. 请不要将不同种类的数据混合处理；
  7. 数据读取后的表头来源于导入文件的文件名，请认真命名并养成好习惯。`

var debug = true
var cacheFileExist = checkFileExist(getCurrentPath() + "cache.json")

// Item 通用
type Item struct {
	Name string
	Data []byte
}

// Cache 通用
type Cache struct {
	Data []*Item
}

var cache Cache

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if debug {
		fmt.Println("调试模式：", debug)
		fmt.Println("程序参数：", os.Args)
		fmt.Println("程序参数长度：", len(os.Args))
	}
	switch {
	case len(os.Args) > 1: // 输入文件
		if debug {
			fmt.Println("文件列表：", os.Args[1:])
		}
		if cacheFileExist {
			readCache()
		}
		readFile(os.Args[1:])
		writeCache()
	case len(os.Args) == 1 && cacheFileExist: // 输出
		readCache()
		analyse()
		writeOutput()
	case len(os.Args) == 1 && !cacheFileExist: // 显示文档
		fmt.Println(doc)
	}
	fmt.Println("程序结束")
	fmt.Scanln()
}

func readCache() {
	cacheByte, err := ioutil.ReadFile(getCurrentPath() + "cache.json")
	check(err)
	json.Unmarshal(cacheByte, &cache)
}
func readFile(filePathList []string) {
	// 按顺序读取单个文件
	for i := 0; i < len(filePathList); i++ {
		fileByte, err := ioutil.ReadFile(filePathList[i])
		check(err)
		os.Stdout.Write(fileByte)
		newItem := new(Item)
		newItem.Name = filePathList[i]
		newItem.Data = fileByte
		cache.Data = append(cache.Data, newItem)
	}
}
func writeCache() {
	cacheByte, err := json.Marshal(cache)
	check(err)
	ioutil.WriteFile(getCurrentPath()+"cache.json", cacheByte, 0666)

}
func analyse() {

}
func writeOutput() {
	os.Remove(getCurrentPath() + "cache.json")
}
func getCurrentPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	return string(path[0 : i+1])
}
func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}