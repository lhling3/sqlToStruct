package main

import (
	"fmt"
	"gocode/sqlToStruct/src/model"
	"os"
	"os/exec"
)

func main() {
	//sql语句所在文件路径
	inputSql := "src/input/input.sql"
	//输出的go文件所在包
	pkgName := "entity"
	//输出go文件名
	filename := "src/entity/output.go"
	//读取建表的sql语句文件
	sqlData, err := os.ReadFile(inputSql)
	if err != nil {
		fmt.Println("读取文件出错")
	}
	sqlStr := string(sqlData)
	//解析sql语句为go struct string
	goStr, err := model.SqlStrToGo(sqlStr, pkgName)
	if err != nil {
		fmt.Println("解析失败，不能生成对应struct")
	}
	//根据输出文件名创建文件并写入保存
	fIn, err := os.Create(filename)
	if err != nil {
		fmt.Println("创建文件出错")
	}
	_, err1 := fIn.WriteString(goStr)
	if err1 != nil {
		fmt.Println("写入文件出错！")
	}
	// 使用gofmt 整理一下输入文件格式
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("go fmt:%s", err.Error())
	}
}
