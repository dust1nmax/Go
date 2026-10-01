package redis

import (
	"Go/web_app/settings"
	"fmt"

	"github.com/go-redis/redis"
	//"github.com/spf13/viper"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:    fmt.Sprintf("%s:%d",
		cfg.Host,
		cfg.Port,
	),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return nil
}

func Close(){
	_ = rdb.Close()
}