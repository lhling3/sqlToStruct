package model

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strings"
)

// 数据库字段类型转换成Go的相应数据类型
var DBTypeToStructType = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int8",
	"smallint":           "int16",
	"mediumint":          "int32",
	"bigint":             "int64",
	"int unsigned":       "uint",
	"integer unsigned":   "uint",
	"tinyint unsigned":   "uint8",
	"smallint unsigned":  "uint16",
	"mediumint unsigned": "uint32",
	"bigint unsigned":    "uint64",
	"bit":                "byte",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

var DBSkipColName = map[string]string{
	"is_deleted":  "is_deleted",
	"update_time": "update_time",
	"update_user": "update_user",
	"create_time": "create_time",
	"create_user": "create_user",
}

func SqlStrToGo(createTabStr string, pkgName string) (string, error) {
	statement, err := sqlparser.ParseStrictDDL(createTabStr)
	if err != nil {
		return "", err
	}
	staDDL, ok := statement.(*sqlparser.DDL)
	if !ok {
		return "", fmt.Errorf("input sql is not a create statment")
	}
	tableName := staDDL.NewName.Name.String()
	// convert to Go struct
	structStr, err := staToGoStruct(staDDL, tableName, pkgName)
	if err != nil {
		return "", err
	}
	//生成指定 methods
	methodStr, err := staToMethods(staDDL, tableName)
	res := fmt.Sprintf(structStr + "\n\n" + methodStr)
	return res, nil
}

func staToMethods(staDDL *sqlparser.DDL, tableName string) (string, error) {
	builder := strings.Builder{}
	structName := snakeCaseToCamel(tableName)
	//tableNameMethod
	funcName := fmt.Sprintf("func (m *%s)TableName() string{\n", structName)
	builder.WriteString(funcName)
	reStr := fmt.Sprintf("return \"%s\"\n", tableName)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")

	var inParam string = "db *gorm.DB"
	var outParam string = "err error"
	var returnHead string = "helpers.WrapError("
	//AddMethod
	funcName = fmt.Sprintf("func (m *%s)Add(%s) (%s){\n", structName, inParam, outParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sdb.Table(m.TableName()).Create(m).Error)\n", returnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//UpdateMethod
	funcName = fmt.Sprintf("func (m *%s)Update(%s) (%s){\n", structName, inParam, outParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sdb.Table(m.TableName()).Updates(m).Error)\n", returnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//FirstMethod
	inParam = "ctx context.Context"
	funcName = fmt.Sprintf("func (m *%s)First(%s) (%s){\n", structName, inParam, outParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).First(m).Error)\n", returnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//ListsMethod
	funcName = fmt.Sprintf("func (m *%s)Lists(%s,lists *[]%s) (%s){\n", structName, inParam, structName, outParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Order(\"id desc\").Find(lists).Error)\n", returnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//CountMethod
	funcName = fmt.Sprintf("func (m *%s)Count(%s, count *int64) (%s){\n", structName, inParam, outParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Count(count).Error)\n", returnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")

	return builder.String(), nil
}

func staToGoStruct(staDDL *sqlparser.DDL, tableName string, pkgName string) (string, error) {
	builder := strings.Builder{}
	header := fmt.Sprintf("package %s\n", pkgName)
	// import time package
	headerPkg := "import (\n" +
		"\t\"time\"\n" +
		")\n\n"
	//用于判断是否需要引入time包
	importTime := false

	structName := snakeCaseToCamel(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	for _, col := range staDDL.TableSpec.Columns {
		columnType := col.Type.Type
		if col.Type.Unsigned {
			columnType += " unsigned"
		}

		goType := DBTypeToStructType[columnType]
		if goType == "time.Time" {
			importTime = true
		}

		field := snakeCaseToCamel(col.Name.String())
		tags := col.Type.Comment
		if tags == nil {
			builder.WriteString(fmt.Sprintf("\t%s\t\t%s\t\t\t\t`json:\"%s\" gorm:\"column:%s\"` \n",
				field, goType, col.Name.String(), col.Name.String()))
		} else {
			builder.WriteString(fmt.Sprintf("\t%s\t\t%s\t\t\t\t`json:\"%s\" gorm:\"column:%s\"` \t\t\t\t//%s \n",
				field, goType, col.Name.String(), col.Name.String(), string(tags.Val)))
		}
	}
	builder.WriteString("}\n")

	if importTime {
		return header + headerPkg + builder.String(), nil
	}
	return header + builder.String(), nil
}

// 将数据库表字段的下划线命名改成驼峰命名
func snakeCaseToCamel(str string) string {
	builder := strings.Builder{}
	index := 0
	if str[0] >= 'a' && str[0] <= 'z' {
		builder.WriteByte(str[0] - ('a' - 'A'))
		index = 1
	}
	for i := index; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			if str[i+1] >= 'a' && str[i+1] <= 'z' {
				builder.WriteByte(str[i+1] - ('a' - 'A'))
				i++
				continue
			}
		}
		builder.WriteByte(str[i])
	}
	return builder.String()
}
