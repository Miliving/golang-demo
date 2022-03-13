package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	selectTagBegin string = "<select"
	selectTagEnd   string = "</select>"

	updateTagBegin string = "<update"
	updateTagEnd   string = "</update>"

	insertTagBegin string = "<insert"
	insertTagEnd   string = "</insert>"

	deleteTagBegin string = "<delete"
	deleteTagEnd   string = "</delete>"
)

/**
 ** 获取指定目录下的所有文件,包含子目录下的文件
 */
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".xml")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

/**
 * 找出指定目录下Mybatis的xml文件中指定表名的sql语句
 *
 */
func processFile(filePath string) (selectSql []string, updateSql []string, deleteSql []string, insertSql []string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件读取错误")
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)

	var content string
	var selectFlag = false
	var updateFlag = false
	var deleteFlag = false
	var insertFlag = false

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		s := string(a)

		// select语句处理
		if strings.Contains(s, selectTagBegin) {
			selectFlag = true // 说明select标签开始了
		}

		if strings.Contains(s, selectTagEnd) {
			content += s
			content += "\n"
			selectFlag = false // 说明select标签结束了
			selectSql = append(selectSql, content)
			content = ""
		}
		if selectFlag {
			content += s
			content += "\n"
		}

		// update 语句处理
		if strings.Contains(s, updateTagBegin) {
			updateFlag = true // 说明select标签开始了
		}

		if strings.Contains(s, updateTagEnd) {
			content += s
			content += "\n"
			updateFlag = false // 说明select标签结束了
			updateSql = append(updateSql, content)
			content = ""
		}
		if updateFlag {
			content += s
			content += "\n"
		}
		// delete 语句处理
		if strings.Contains(s, deleteTagBegin) {
			deleteFlag = true // 说明select标签开始了
		}

		if strings.Contains(s, deleteTagEnd) {
			content += s
			content += "\n"
			deleteFlag = false // 说明select标签结束了
			deleteSql = append(deleteSql, content)
			content = ""
		}
		if deleteFlag {
			content += s
			content += "\n"
		}
		// insert 语句处理
		if strings.Contains(s, insertTagBegin) {
			insertFlag = true // 说明select标签开始了
		}

		if strings.Contains(s, insertTagEnd) {
			content += s
			content += "\n"
			insertFlag = false // 说明select标签结束了
			insertSql = append(insertSql, content)
			content = ""
		}
		if insertFlag {
			content += s
			content += "\n"
		}
	}
	return
}

func main() {

	xfiles, _ := GetAllFiles("D:/JoyoProject/smarthr-financing/smarthr-finance/src/main/resources/mapper")
	for _, file := range xfiles {

		file := strings.ReplaceAll(file, "\\", "/")
		// var printFileNameflag = false
		fmt.Printf("当前处理的文件为：%s\n", file)

		selectSql, _, _, _ := processFile(file)

		for _, v := range selectSql {
			if strings.Contains(v, "product_bill_detail ") {
				// printFileNameflag = true
				fmt.Println(v)
			}
		}
		// if printFileNameflag {
		// 	fmt.Println(file)
		// }
	}

	// var file string = "D:/JoyoProject/smarthr-financing/smarthr-finance/src/main/resources/mapper/trade/FincBankAcctTradeDetailMapper.xml"

	// selectSql, _, _, _ := processFile(file)

	// for _, v := range selectSql {
	// 	if strings.Contains(v, "bank_acct_trade_detail") {
	// 		fmt.Println(v)
	// 	}
	// }

}
