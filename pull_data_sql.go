package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	custSql = `select 
    cust_cn 客户名称,
    case cust_status
    when 0 then '无状态'
    when 1 then '跟进中'
    when 2 then '已签约'
    when 3 then '黑名单'
    when 4 then '已撤户'
    when 5 then '服务中'
    else '其他' end 状态
  from cust_info where cust_cn in (`

	clueSql = `
  select 
    clue_name 客户名称,
    case clue_status
    when 1 then '未跟进'
    when 2 then '跟进中'
    when 3 then '转客户' 
    when 4 then '失效信息'
    when 5 then '签约中'
    when 6 then '已签约'
    else '其他' end 线索状态
  from clue_info where clue_name in (`

	mobileSql = `
  select 
    clue_name 客户名称,
    contacts_mobile 手机号,
    case clue_status
    when 1 then '未跟进'
    when 2 then '跟进中'
    when 3 then '转客户' 
    when 4 then '失效信息'
    when 5 then '签约中'
    when 6 then '已签约'
    else '其他' end 线索状态
  from clue_info where contacts_mobile in (`
)

/**
** 协助客服拉CRM相关数据SQL
**/
func main() {
	var file string
	var sqlType string
	var dataType string

	flag.StringVar(&file, "file", "./data.txt", "文件路径")
	flag.StringVar(&sqlType, "sql", "cust", "sql类型：cust/cule/mobile")
	flag.StringVar(&dataType, "type", "int", "数据类型：int/string")

	if file == "" {
		fmt.Println("文件路径不能空")
	}
	flag.Parse()

	fmt.Println(fmt.Sprintf("文件名称：%s, 数据类型：%s, sql类型：%s", file, dataType, sqlType))

	if sqlType != "cust" && sqlType != "clue" && sqlType != "mobile" {
		fmt.Println("sql类型必须是：cust/cule/mobile")
		return
	}

	inputFile, inputError := os.Open(file)
	if inputError != nil {
		fmt.Println("文件读取失败")
		return
	}
	defer inputFile.Close()

	fileWriter, err := os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开文件失败")
	}
	defer fileWriter.Close()

	inputReader := bufio.NewReader(inputFile)

	if sqlType == "cust" {
		fileWriter.WriteString(custSql + "\n")
	} else if sqlType == "clue" {
		fileWriter.WriteString(clueSql + "\n")
	} else {
		fileWriter.WriteString(mobileSql + "\n")
	}

	for {
		lineString, err := inputReader.ReadString('\n')
		lineString = strings.TrimSpace(lineString)

		if dataType == "int" {
			if err != io.EOF {
				lineString = lineString + "," + "\n"
			} else {
				lineString = lineString + ");\n"
			}
		} else {
			if err != io.EOF {
				lineString = "'" + lineString + "'," + "\n"
			} else {
				lineString = "'" + lineString + "'" + ");\n"
			}
		}
		fileWriter.WriteString(lineString)

		if err == io.EOF {
			fmt.Println("File read End!")
			break
		}

	}
}
