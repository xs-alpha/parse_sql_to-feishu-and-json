package service

// author：xiaosheng

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"xiaosheng/moudle"
	"xiaosheng/tools"
	"xiaosheng/util"
)

var (
	FeishuStringList = make([]moudle.StructModel, 0)
	Wg               = sync.WaitGroup{}
	FeishuStringJson = make([]string, 0)
)

type Model = moudle.StructModel

var FileStruct = new(moudle.FileStruct)
var configObj = new(moudle.IniModel)

type Operate struct {
	str string
}

func init() {
	FileStruct.JsonFileName = tools.JsonFileName
	FileStruct.XLSXFileName = tools.XLSXFileName
	FileStruct.DirPath = tools.DirPath
	FileStruct.ConfigPath = tools.ConfigPath
	FileStruct.FeishuParseFile = tools.FeishuParseFile
}

func (op *Operate) del(trim string) *Operate {
	op.str = strings.Replace(op.str, trim, "", -1)
	op.str = strings.TrimSpace(op.str)
	return op
}

func GetSql() {
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
			fmt.Println("[*]->:一共有", totLine, "行内容")
			rangeArray()
			Wg.Add(1)
			go writeToFile()
			Wg.Add(1)
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

	//如果含有表名或者
	if strings.Contains(lower, "table") {
		split := strings.Split(lower, " ")
		for _, item := range split {
			if strings.Contains(item, "`") {
				// 获取表名
				tableName := strings.Replace(item, "`", "", -1)
				jsonFileName := tableName + ".json"
				xlsxFileName := tableName + ".xlsx"
				FileStruct.JsonFileName = jsonFileName
				FileStruct.XLSXFileName = xlsxFileName
			}
		}
		fmt.Println(split)
		return
	}

	// 判断是否包含基本数据类型
	flag := false
	for _, item := range moudle.SqlKeyWord {
		if strings.Contains(strings.ToLower(content), item) {
			flag = true
		}
	}
	if !flag {
		return
	}

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
	if configObj.HupName {
		m.HumpName()
	}
	FeishuStringList = append(FeishuStringList, *m)
}

// writeToFile 写json
func writeToFile() {
	filename := path.Join(FileStruct.DirPath, FileStruct.JsonFileName)
	create := util.FileCreate(filename, configObj.NewJsonAndSqlFile)
	if create != "" {
		filename = create
	}
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	Wg.Done()
}

func writeToExcel() error {
	activeSheetName := tools.ActiveSheetName
	fileNamePath := path.Join(FileStruct.DirPath, FileStruct.XLSXFileName)
	exists, err := util.FileExists(fileNamePath)
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
	Wg.Done()
	return nil
}

func rangeArray() {
	fmt.Println("---------------------")
	for _, item := range FeishuStringList {
		fmt.Println(item)
	}
}

func Init() {
	// 初始化配置文件
	// 如果配置文件不存在，创建后往配置文件写
	configFileName := FileStruct.ConfigPath + tools.IniConfigFileName
	util.PathFileExists(FileStruct.ConfigPath, false)
	exists, _ := util.FileExists(configFileName)
	if !exists {
		f, err := os.Create(configFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(tools.IniConfig)
	}
	// 初始化sql文件
	util.FileCreate(tools.SqlName, configObj.NewJsonAndSqlFile)

	err := ini.MapTo(configObj, path.Join(tools.ConfigPath, tools.IniConfigFileName))
	if err != nil {
		log.Fatal(err)
	}
	coverConfig()

	util.PathFileExists(FileStruct.DirPath, false)
	filePattern := strings.Split(FileStruct.FeishuParseFile, ".")
	if len(filePattern) > 1 {
		FileStruct.FeiShuParseFileResult = filePattern[0] + "result." + filePattern[1]
	} else {
		FileStruct.FeiShuParseFileResult = filePattern[0] + "result"
	}
	Wg.Done()
}

func coverConfig() {
	if len(configObj.OutPutDir) != 0 {
		FileStruct.DirPath = configObj.OutPutDir
	}
	if len(configObj.FeishuParseFile) != 0 {
		FileStruct.FeishuParseFile = configObj.FeishuParseFile
	}
}

func ParseFeishu() {
	f, err := os.Open(FileStruct.FeishuParseFile)
	defer f.Close()
	if err != nil {
		fmt.Println("err:----", err)
	}
	reader := bufio.NewReader(f)
	totLine := 0
	for {
		content, isPrefix, err := reader.ReadLine()

		//当单行的内容超过缓冲区时，isPrefix会被置为真；否则为false；
		if !isPrefix {
			totLine++
		}
		cont := string(content)
		if len(strings.TrimSpace(cont)) > 0 {
			if configObj.FeishuJsonHupName {
				moudle.HumpNameInJson(&cont)
			}
			eachLine := "\"" + cont + "\":\"\""
			FeishuStringJson = append(FeishuStringJson, eachLine)
		}
		if err == io.EOF {
			fmt.Println("[*]->:一共解析有", totLine, "行内容")
			break
		}
	}
	writeToJson()
}

func writeToJson() {
	f, _ := os.OpenFile(path.Join(FileStruct.DirPath, FileStruct.FeiShuParseFileResult), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.WriteString("{\n")
	for index, item := range FeishuStringJson {
		insertStr := strings.TrimSpace(item)
		if index != len(FeishuStringJson)-1 {
			insertStr += ","
		}
		_, err := f.WriteString("\t" + insertStr + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	f.WriteString("}\n")
}
