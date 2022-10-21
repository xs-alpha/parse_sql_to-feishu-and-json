package main

import (
	"fmt"
	"os"
	"xiaosheng/service"
)

func main() {
	// 初始化配置文件，判断有没有配置文件和sql文件

	// 用配置文件中的配置覆盖默认配置

	fmt.Println("[*]->:\t正在进行初始化配置文件和sql文件")
	service.Wg.Add(1)
	go service.Init()

	var excute = 1

	for {
		fmt.Println("\n请输入要进行的操作： 1：根据sql生成json和飞书表单  2： 根据飞书表格字段名生成json 3:退出")
		fmt.Print("[*]->:")
		fmt.Scanln(&excute)
		if excute == 1 {
			service.GetSql()
		} else if excute == 2 {
			service.Wg.Add(1)
			go func() {
				stat, _ := os.Stat(service.FileStruct.FeishuParseFile)
				if stat == nil {
					f, err := os.Create(service.FileStruct.FeishuParseFile)
					if err != nil {
						fmt.Println("err occured when create file: ", err)
					}
					f.Close()
				}

				service.Wg.Done()
			}()
			fmt.Println("[*]->:在填充完当前目录下的解析文件后按回车")
			fmt.Print("[*]->:")
			enter := ""
			fmt.Scanln(&enter)
			service.ParseFeishu()
		} else {
			break
		}
	}

	service.Wg.Wait()
	// 根据数据库字段名自动生成json, 添加cmd args处理或者从用户端输入
}
