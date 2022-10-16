package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	SqlKeyWord       = []string{"int", "varchar", "date", "timestamp", "bigint", "tinyint", "bool", "double", "float", "decimal", "char", "text", "enum", "bit", "set", "binary", "blob", "null"}
	FeishuStringList = make([]string, 0)
)

type Model struct {
	Name        string
	Alias       string
	IsNecessary string
	NameNote    string
}

type Operate struct {
	str string
}

func (op *Operate) del(trim string) *Operate {
	op.str = strings.Replace(op.str, trim, "", -1)
	op.str = strings.TrimSpace(op.str)
	return op
}

func getSql() {
	f, err := os.Open("./parse.sql")
	defer f.Close()
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	//func NewReader(rd io.Reader) *Reader
	reader := bufio.NewReader(f)

	totLine := 0
	for {
		//func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
		content, isPrefix, err := reader.ReadLine()

		//fmt.Println(string(content), isPrefix, err)

		//当单行的内容超过缓冲区时，isPrefix会被置为真；否则为false；
		if !isPrefix {
			totLine++
		}

		parse(string(content))
		if err == io.EOF {
			fmt.Println("一共有", totLine, "行内容")
			rangeArray()
			writeToFile()
			break
		}

	}
}

// IsChineseChar 判断是否是中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
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
	for _, item := range SqlKeyWord {
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
		for index, eachItem := range splitStringList {
			if len(eachItem) == 0 {
				continue
			}
			// 未后续操作做铺垫
			o := new(Operate)
			(*o).str = eachItem

			if strings.HasPrefix(eachItem, "`") && strings.HasSuffix(eachItem, "`") {
				// todo:去掉``
				o.del("`")
				m.Alias = o.str
			}
			//if strings.HasSuffix(eachItem, "'")||strings.HasSuffix(eachItem, ",")||IsChineseChar(eachItem)
			if index == len(splitStringList)-1 {
				// 去掉‘’或者“”“”获取注释
				o.del("'").del("\"").del(",")
				m.NameNote = o.str
				m.Name = o.str
			}

		}
		insertStr := m.Name + "\n" + m.Alias + "\n" + m.IsNecessary + "\n" + m.NameNote
		FeishuStringList = append(FeishuStringList, insertStr)
	}
	fmt.Println(FeishuStringList)
}

func writeToFile() {
	f, err := os.Create("afile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, line := range FeishuStringList {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeToExcel() error {
	activeSheetName := "test"
	fileNamePath := "./test.xlsx"
	//_ = tools.PathDirCreate("./test.xlxs") //不存在创建目录
	//activeSheetName := "Sheet1"
	//fileNamePath := path.Join(dirPath, fileName)
	//exists, err := tools.PathFileExists(fileNamePath) //判断文件是否存在创建目录
	exists, err := os.Stat(fileNamePath)
	rowNum := 0
	lastLineNum := 0
	var f *excelize.File
	// 创建excel
	if exists != nil || err != nil {
		f = excelize.NewFile()
		// Create a new sheet.
		index := f.NewSheet(activeSheetName)

		// Set active sheet of the workbook.
		f.SetActiveSheet(index)
		// Set tabletitle value of a cell.
		tableInfo := map[string]string{
			"A1": "Id",
			"B1": "Filename",
			"C1": "Product",
			"D1": "Fofaquery",
		}
		for k, v := range tableInfo {
			f.SetCellValue(activeSheetName, k, v)
		}
	} else { // 追加写入excel
		f, _ = excelize.OpenFile(fileNamePath)
		rows, _ := f.GetRows(activeSheetName)
		lastLineNum = len(rows) //找到最后一行
	}
	// Set table content value of a cell.
	for index, list := range FeishuStringList {
		if exists != nil || err != nil {
			//如果不存在从第2行写入
			rowNum = index + 2
		} else {
			//否则从文件内容尾行写入
			rowNum = lastLineNum + index + 1
		}
		f.SetCellValue(activeSheetName, fmt.Sprintf("A%d", rowNum), list.Id)
		f.SetCellValue(activeSheetName, fmt.Sprintf("B%d", rowNum), list.Filename)
		f.SetCellValue(activeSheetName, fmt.Sprintf("C%d", rowNum), list.Product)
		f.SetCellValue(activeSheetName, fmt.Sprintf("D%d", rowNum), list.Fofaquery)
	}
	// Save spreadsheet by the given path.  static/downloads/Book1.xlsx
	if err := f.SaveAs(fileNamePath); err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("save file failed, path:(%s)", fileNamePath))
	}
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
}
