package main

// author：xiaosheng

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"xiaosheng/tools"
)

var (
	FeishuStringList = make([]tools.StructModel, 0)
	wg               = sync.WaitGroup{}
)

type Model = tools.StructModel

type Operate struct {
	str string
}

func (op *Operate) del(trim string) *Operate {
	op.str = strings.Replace(op.str, trim, "", -1)
	op.str = strings.TrimSpace(op.str)
	return op
}

func getSql() {
	f, err := os.Open(tools.SqlName)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(f)

	totLine := 0
	for {
		content, isPrefix, err := reader.ReadLine()

		//当单行的内容超过缓冲区时，isPrefix会被置为真；否则为false；
		if !isPrefix {
			totLine++
		}

		parse(string(content))
		if err == io.EOF {
			fmt.Println("一共有", totLine, "行内容")
			rangeArray()
			wg.Add(1)
			go writeToFile()
			wg.Add(1)
			go writeToExcel()
			break
		}

	}
}

func parse(content string) {
	pre := strings.ToLower(content)
	lower := strings.Replace(pre, "\n", "", -1)
	lower = strings.Replace(lower, "\r", "", -1)
	lower = strings.Replace(lower, "\t", "", -1)

	// 如果为空，返回
	if len(strings.TrimSpace(lower)) == 0 || strings.Contains(content, ";") {
		return
	}

	// 判断是否包含基本数据类型
	flag := false
	for _, item := range tools.SqlKeyWord {
		if strings.Contains(strings.ToLower(content), item) {
			flag = true
		}
	}
	if !flag {
		return
	}

	//如果含有表名或者
	if strings.Contains(lower, "table") {
		return
	} else {
		splitStringList := strings.Split(lower, " ")
		m := new(Model)
		m.IsNecessary = "否"
		originTypeIndex := 0
		for index, eachItem := range splitStringList {
			if len(eachItem) == 0 || strings.EqualFold(eachItem, " ") {
				originTypeIndex++
				continue
			}
			// 未后续操作做铺垫
			o := new(Operate)
			(*o).str = eachItem

			if strings.HasPrefix(eachItem, "`") && strings.HasSuffix(eachItem, "`") {
				// todo:去掉``
				o.del("`")
				m.Alias = o.str
				originTypeIndex = index + 1
			}
			//if strings.HasSuffix(eachItem, "'")||strings.HasSuffix(eachItem, ",")||IsChineseChar(eachItem)
			if index == originTypeIndex {
				o.del("'").del("\"")
				m.OriginType = o.str
			}
			m.IsNecessary = "否"
			if index == len(splitStringList)-1 {
				// 去掉‘’或者“”“”获取注释
				o.del("'").del("\"").del(",")
				m.NameNote = o.str
				m.Name = o.str
				//m.Type = "number"
			}

		}
		m.JudgeType()
		m.DealWithName()
		FeishuStringList = append(FeishuStringList, *m)
	}
	fmt.Println(FeishuStringList)
}

func writeToFile() {
	f, err := os.Create("json.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("{\n")
	for i, m := range FeishuStringList {
		insertStr := "\"" + m.Alias + "\"" + ": " + "\"\""
		if i != len(FeishuStringList)-1 {
			insertStr += ","
		}
		_, err := f.WriteString("\t" + insertStr + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	f.WriteString("}\n")
	wg.Done()
}

func writeToExcel() error {
	activeSheetName := tools.ActiveSheetName
	fileNamePath := path.Join(tools.DirPath, tools.XLSXFileName)
	exists, err := tools.FileExists(fileNamePath)
	rowNum := 0
	lastLineNum := 0
	var f *excelize.File
	// 创建excel
	if !exists || err != nil {
		f = excelize.NewFile()
		// Create a new sheet.
		index := f.NewSheet(activeSheetName)

		f.SetActiveSheet(index)
		tableInfo := map[string]string{
			"A1": "参数名",
			"B1": "变量",
			"C1": "类型",
			"D1": "必填",
			"E1": "描述",
		}
		for k, v := range tableInfo {
			f.SetCellValue(activeSheetName, k, v)
		}
	} else { // 追加写入excel
		f, _ = excelize.OpenFile(fileNamePath)
		rows, _ := f.GetRows(activeSheetName)
		lastLineNum = len(rows) //找到最后一行
	}
	// Set
	for index, list := range FeishuStringList {
		if exists || err != nil {
			//如果不存在从第2行写入
			rowNum = index + 2
		} else {
			//否则从文件内容尾行写入
			rowNum = lastLineNum + index + 1
		}
		f.SetCellValue(activeSheetName, fmt.Sprintf("A%d", rowNum), list.Name)
		f.SetCellValue(activeSheetName, fmt.Sprintf("B%d", rowNum), list.Alias)
		f.SetCellValue(activeSheetName, fmt.Sprintf("C%d", rowNum), list.Type)
		f.SetCellValue(activeSheetName, fmt.Sprintf("D%d", rowNum), list.IsNecessary)
		f.SetCellValue(activeSheetName, fmt.Sprintf("E%d", rowNum), list.NameNote)
	}
	// Save
	if err := f.SaveAs(fileNamePath); err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("save file failed, path:(%s)", fileNamePath))
	}
	wg.Done()
	return nil
}

func rangeArray() {
	fmt.Println("---------------------")
	for _, item := range FeishuStringList {
		fmt.Println(item)
	}
}

func main() {
	/**
	按行赋值

	func (f *File) SetSheetRow(sheet, cell string, slice interface{}) error

	根据给定的工作表名称、起始坐标和 slice 类型引用按行赋值。此功能是并发安全的。例如，在名为 Sheet1 的工作表第 6 行上，以 B6 单元格作为起始坐标按行赋值：

	err := f.SetSheetRow("Sheet1", "B6", &[]interface{}{"1", nil, 2})

	*/
	getSql()
	wg.Wait()
}
