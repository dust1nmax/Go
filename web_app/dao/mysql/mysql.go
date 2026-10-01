package mysql

import (
	"Go/web_app/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
	cfg.User,
	cfg.Password,
	cfg.Host,
	cfg.Port,
	cfg.Dbname,
)
	// 也可以使用MustConnect连接不成功就panic
	// Open + ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed ", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.MaxOpenConns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.MaxIdleConns"))
	return
}

func Close(){
	_ = db.Close()
}