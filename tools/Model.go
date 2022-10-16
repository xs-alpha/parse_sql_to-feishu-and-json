package tools

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
}

func FieldsInclude(oriStr string, destStr string) bool {
	if strings.Contains(strings.ToLower(oriStr), destStr) {
		return true
	}
	return false
}
