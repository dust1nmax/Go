package main

import (
	"flag"
	"fmt"
	"time"
)

// func main(){
// 	fmt.Println(os.Args)
// 	if len(os.Args) > 0{
// 		for index, arg := range os.Args{
// 			fmt.Printf("args[%d]=%v\n", index, arg)
// 		}
// 	}
// }

func main(){
	var name string
	var age int 
	var married bool
	var delay time.Duration
	flag.StringVar(&name, "name", "max", "姓名")
	flag.IntVar(&age, "age", 22, "年龄")
	flag.BoolVar(&married, "married", false, "婚否")
	flag.DurationVar(&delay, "delay", 0, "延迟的时间间隔")

	//解析命令行参数
	flag.Parse()
	fmt.Println(name, age, married, delay)
	//返回命令行参数后的其他参数
	fmt.Println(flag.Args())
	//返回命令行参数后的其他参数的个数
	fmt.Println(flag.NArg())
	//返回使用的命令行参数的个数
	fmt.Println(flag.NFlag())

	//使用例子
	//1.go build
	//2../flag_demo 
	//3../flag_demo -h
	
}