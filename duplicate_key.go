package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/**
 * 对类似csv文件处理，找出首列key的重复，然后value整行文本
 * abc,123124
 * xyz,798038
 * abc,hello
 */
func main() {
	var file string
	flag.StringVar(&file, "file", "./data.txt", "文件路径")

	if file == "" {
		fmt.Println("文件路径不能空")
	}
	flag.Parse()

	fmt.Println(file)

	inputFile, inputError := os.Open(file)
	if inputError != nil {
		fmt.Println("文件读取失败")
		return
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)

	var slice1 []string = make([]string, 5) // 首列的key数组
	dataMap := make(map[string][]string)
	for {
		lineString, err := inputReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		s := strings.ReplaceAll(lineString, "\r\n", "")
		s2 := strings.Split(s, ",")
		slice1 = append(slice1, s2[0])
		dataMap[s2[0]] = append(dataMap[s2[0]], lineString)
	}

	// 写入文件
	writeFile(&dataMap)

}

/**
 * 将重复的数据写入到文件中
 */
func writeFile(dataMap *map[string][]string) {
	fileWriter, err := os.OpenFile("./data.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开文件失败")
	}
	defer fileWriter.Close()

	b, _ := json.Marshal(*dataMap)
	fmt.Println(string(b))

	for _, value := range *dataMap {
		if len(value) > 1 {
			for i := 0; i < len(value); i++ {
				fileWriter.WriteString(value[i])
			}
		}
	}
}
