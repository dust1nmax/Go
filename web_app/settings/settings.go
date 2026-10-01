package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量 用来保存程序的所以配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level       string `mapstructure:"level"`
	FileName    string `mapstructure:"filename"`
	Max_size    int    `mapstructure:"max_size"`
	Max_age     int    `mapstructure:"max_age"`
	Max_backups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	Dbname        string `mapstructure:"dbname"`
	MMaxOpenConns int    `mapstructure:"maxopenconns"`
	MaxIdleConns  int    `mapstructure:"maxidleconns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolsize"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml") //专用与从远程获取配置信息
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
		return
	}

	// 将读取到的配置信息反序列化到Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
		}
	})
	return
}
