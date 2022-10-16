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
