package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/**
 ** 协助客服拉数据SQL
 **/
func main() {
	var custSql string = `select 
	cust_cn 客户名称,
	case cust_status
	when 0 then '无状态'
		when 1 then '跟进中'
		when 2 then '已签约'
		when 3 then '黑名单'
		when 4 then '已撤户'
		when 5 then '服务中'
		else '其他'  end 状态,
		create_time as '创建时间'
	from cust_info where cust_cn in (`


	var clueSql string = "select" +
		"clue_name 客户名称,"+
		"case clue_status"+
		"when 1 then '未跟进'"+
		"when 2 then '跟进中'"+
		"when 3 then '转客户'"+
		"when 4 then '失效信息'"+
		"when 5 then '签约中'"+
		"when 6 then '已签约'"+
		"else '其他' end 线索状态,"+
		"create_time as '创建时间',"+
		"concat(t3.real_name,'(',t3.staff_code,')') as '当前跟进人'"+
	"FROM"+
		"clue_info t1"+
		"LEFT JOIN clue_track t2 ON t2.clue_id  = t1.id "+
		"left join `smarthr-privilege`.staff t3 on t3.id = t2.staff_id"+
	"where t1.clue_name in ("

	var mobileSql = "select " +
		"t1.clue_name 客户名称," +
		"t1.contacts_mobile 手机号," +
		"case t1.clue_status" +
		"when 1 then '未跟进'" +
		"when 2 then '跟进中'" +
		"when 3 then '转客户' " +
		"when 4 then '失效信息'" +
		"when 5 then '签约中'" +
		"when 6 then '已签约'" +
		"else '其他' end 线索状态," +
		"t1.create_time as '创建时间'," +
		"concat(t3.real_name,'(',t3.staff_code,')') as '当前跟进人'"+
	"FROM" +
		"clue_info t1" +
		"LEFT JOIN clue_track t2 ON t2.clue_id  = t1.id " +
		"left join `smarthr-privilege`.staff t3 on t3.id = t2.staff_id" +
	"where t1.contacts_mobile in ("


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
