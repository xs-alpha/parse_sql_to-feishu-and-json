package tools

// author：xiaosheng

import "strings"

var SqlKeyWord = []string{"int", "varchar", "date", "timestamp", "bigint", "tinyint", "bool", "double", "float", "decimal", "char", "text", "enum", "bit", "set", "binary", "blob"}

type StructModel struct {
	Name        string
	Alias       string
	Type        string
	IsNecessary string
	NameNote    string
	OriginType  string
}

type IniModel struct {
	SqlConfig      `ini:"sql"`
	FileNameConfig `ini:"filename"`
}

type SqlConfig struct {
	HupName           bool `ini:"hupName"`
	NewJsonAndSqlFile bool `ini:"newJsonAndSqlFile"`
	TableName         bool `ini:"tableName"`
}

type FileNameConfig struct {
	OutPutDir       string `ini:"outPutDir"`
	FeishuParseFile string `ini:"feishuParseFile"`
}

type FileStruct struct {
	JsonFileName          string
	XLSXFileName          string
	DirPath               string
	ConfigPath            string
	FeishuParseFile       string
	FeiShuParseFileResult string
}

func (t *StructModel) JudgeType() {
	item := t.OriginType
	if FieldsInclude(item, "varchar") {
		t.Type = "string"
	}
	if FieldsInclude(item, "int") {
		t.Type = "integer"
	}
	if FieldsInclude(item, "bigint") {
		t.Type = "number"
	}
	if FieldsInclude(item, "tinyint") {
		t.Type = "integer"
	}
	if FieldsInclude(item, "timestamp") {
		t.Type = "timestamp"
	}
	if FieldsInclude(item, "date") {
		t.Type = "date"
	}
	if FieldsInclude(item, "bool") {
		t.Type = "Boolean"
	}
	if FieldsInclude(item, "double") {
		t.Type = "float"
	}
	if FieldsInclude(item, "float") {
		t.Type = "float"
	}
	if FieldsInclude(item, "decimal") {
		t.Type = "float"
	}
	if FieldsInclude(item, "text") {
		t.Type = "string"
	}
	if FieldsInclude(item, "enum") {
		t.Type = "enum"
	}
	if FieldsInclude(item, "binary") {
		t.Type = "binary"
	}
}

func FieldsInclude(oriStr string, destStr string) bool {
	if strings.Contains(strings.ToLower(oriStr), destStr) {
		return true
	}
	return false
}

func (t *StructModel) DealWithName() {
	if strings.Contains(t.Name, "(") || strings.Contains(t.Name, "（") {
		splitStr := strings.Split(t.Name, "(")
		if strings.EqualFold(splitStr[0], t.Name) {
			splitStr = strings.Split(t.Name, "（")
		}
		t.Name = splitStr[0]
	}
}

// HumpName 驼峰处理
func (t *StructModel) HumpName() {
	if strings.Contains(t.Alias, "_") {
		split := strings.Split(t.Alias, "_")
		allWorld := ""
		for index, item := range split {
			if index == 0 {
				allWorld += item
			}
			if index != 0 && len(item) != 0 {
				firstLetter := strings.ToUpper(string(item[0]))
				otherLetter := item[1:]
				eachWorld := firstLetter + otherLetter
				allWorld += eachWorld
			}
		}
		t.Alias = allWorld
	}
}
