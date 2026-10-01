package main

import (
	"fmt"
	//"net/http"

	"github.com/fsnotify/fsnotify"
	//"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	//"github.com/spf13/pflag"
)

type Config struct {
	Version     string `mapstructure:"version"`
	mysqlConfig `mapstructure:"mysql"`
}

type mysqlConfig struct {
	Port   int    `mapstructure:"port"`
	Host   string `mapstructure:"host"`
	Dbname string `mapstructure:"dbname"`
}

func main() {
	//建立默认值
	viper.SetDefault("fileUrl", "./")

	//读取配置文件
	//viper.SetConfigFile("./config.yaml")  // 指定配置文件路径
	viper.SetConfigName("config")         // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")           // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")              // 还可以在工作目录中查找配置

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// //写入配置文件
	// viper.WriteConfig() // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径
	// viper.SafeWriteConfig()
	// viper.WriteConfigAs("/path/to/my/.config")
	// viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
	// viper.SafeWriteConfigAs("/path/to/my/.other_config")

	//监控并重新读取配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	var c Config
	if err:= viper.Unmarshal(&c); err!= nil{
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return
	}
	
	//%v：按默认格式打印变量。
	//%#v：按照 Go 代码/Go 语法形式打印变量的详细表示。
	fmt.Printf("c:%#v\n", c)

	// r := gin.Default()
	// r.GET("version", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, viper.GetString("version"))
	// })
	// r.Run()

	//		// flags: 命令行参数
	//		pflag.Int("port", 8080, "server port")
	//		//把启动程序时写在命令行里的参数真正解析出来
	//	    pflag.Parse()
	//		//把命令行里的 port 和 Viper 里的 port 绑定起来
	//	    viper.BindPFlag("port", pflag.Lookup("port"))
	//	    fmt.Println(viper.GetInt("port"))
}
