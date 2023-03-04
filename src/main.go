package main

import (
	_ "embed"
	"fmt"
	"gocode/sqlToStruct/src/config"
	"gocode/sqlToStruct/src/model"
	"os"
	"os/exec"
	"strings"
)

// 配置文件struct
var Conf = config.NewConfig()

//go:embed input/input.sql
var sqlData string

func main() {
	//读取建表的sql语句文件
	//sqlData, err := os.ReadFile(Conf.InPath)
	//if err != nil {
	//	fmt.Println("读取文件出错")
	//}
	sqls := string(sqlData)
	sqlArr := strings.Split(sqls, "CREATE")
	for i := 1; i < len(sqlArr); i++ {
		fmt.Println(i)
		sqlStr := fmt.Sprintf("CREATE %s", sqlArr[i])
		//解析sql语句为go struct string
		goStr, err := model.SqlStrToGo(sqlStr, Conf)
		if err != nil {
			fmt.Printf("第%v条记录解析失败，不能生成对应struct", i)
			continue
		}
		//根据输出文件名创建文件并写入保存
		filename := fmt.Sprintf("%s%s.go", Conf.OutPath, Conf.StructName)
		fIn, err := os.Create(filename)
		if err != nil {
			fmt.Println("创建文件出错")
			continue
		}
		_, err1 := fIn.WriteString(goStr)
		if err1 != nil {
			fmt.Println("写入文件出错！")
			continue
		}
		// 使用gofmt 整理一下输入文件格式
		cmd := exec.Command("gofmt", "-w", filename)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("go fmt:%s", err.Error())
		}
	}
}
