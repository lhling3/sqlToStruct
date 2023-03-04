package model

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"gocode/sqlToStruct/src/config"
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

func SqlStrToGo(createTabStr string, Conf *config.SqlToGoConfig) (string, error) {
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
	structStr, err := staToGoStruct(staDDL, tableName, Conf)
	if err != nil {
		return "", err
	}
	//生成指定 methods
	methodStr, err := staToMethods(staDDL, tableName, Conf)
	res := fmt.Sprintf(structStr + "\n\n" + methodStr)
	Conf.StructName = snakeCaseToCamel(tableName)
	return res, nil
}

func staToMethods(staDDL *sqlparser.DDL, tableName string, Conf *config.SqlToGoConfig) (string, error) {
	builder := strings.Builder{}
	structName := snakeCaseToCamel(tableName)
	//tableNameMethod
	funcName := fmt.Sprintf("func (m *%s)TableName() string{\n", structName)
	builder.WriteString(funcName)
	reStr := fmt.Sprintf("return \"%s\"\n", tableName)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")

	//AddMethod
	funcName = fmt.Sprintf("func (m *%s)Add(%s) (%s){\n", structName, Conf.InParam1, Conf.OutParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sdb.Table(m.TableName()).Create(m).Error)\n", Conf.ReturnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//UpdateMethod
	funcName = fmt.Sprintf("func (m *%s)Update(%s) (%s){\n", structName, Conf.InParam1, Conf.OutParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sdb.Table(m.TableName()).Updates(m).Error)\n", Conf.ReturnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//FirstMethod
	funcName = fmt.Sprintf("func (m *%s)First(%s) (%s){\n", structName, Conf.InParam2, Conf.OutParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).First(m).Error)\n", Conf.ReturnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//ListsMethod
	funcName = fmt.Sprintf("func (m *%s)Lists(%s,lists *[]%s) (%s){\n", structName, Conf.InParam2, structName, Conf.OutParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Order(\"id desc\").Find(lists).Error)\n", Conf.ReturnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")
	//CountMethod
	funcName = fmt.Sprintf("func (m *%s)Count(%s, count *int64) (%s){\n", structName, Conf.InParam2, Conf.OutParam)
	builder.WriteString(funcName)
	reStr = fmt.Sprintf("return %sglobal.DB(ctx, m.DbName()).Where(m).Table(m.TableName()).Scopes(m.scopes...).Count(count).Error)\n", Conf.ReturnHead)
	builder.WriteString(reStr)
	builder.WriteString("}\n\n")

	return builder.String(), nil
}

func staToGoStruct(staDDL *sqlparser.DDL, tableName string, Conf *config.SqlToGoConfig) (string, error) {
	builder := strings.Builder{}
	header := fmt.Sprintf("package %s\n", Conf.PkgName)
	// import time package
	headerPkg := "import (\n" +
		"\t\"time\"\n" +
		")\n\n"
	//用于判断是否需要引入time包
	importTime := false

	structName := snakeCaseToCamel(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	isBaseModel := structName == Conf.BaseModel
	if !isBaseModel {
		//写入数据库名
		builder.WriteString(fmt.Sprintf("\t%s\n", Conf.DBName))
		//写入BaseModel
		builder.WriteString(fmt.Sprintf("\t%s\n", Conf.BaseModel))
	}
	for _, col := range staDDL.TableSpec.Columns {
		field := snakeCaseToCamel(col.Name.String())
		//如果是写入BaseModel的字段，不写入
		if !isBaseModel && (field == "Id" || field == "CreateTime" || field == "UpdateTime") {
			continue
		}
		columnType := col.Type.Type
		if col.Type.Unsigned {
			columnType += " unsigned"
		}

		goType := DBTypeToStructType[columnType]
		if goType == "time.Time" {
			importTime = true
		}
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
