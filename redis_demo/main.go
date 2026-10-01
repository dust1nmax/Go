package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func redisExamole() {
	err := rdb.Set("Socre", 100, 0).Err()
	if err != nil {
		fmt.Printf("set socre failed, err%v\n", err)
		return
	}

	val, err := rdb.Get("Socre").Result()
	if err != nil {
		fmt.Printf("get socre failed, err%v\n", err)
		return
	}
	fmt.Println("Socre", val)

	val2, err := rdb.Get("name").Result()
	if err == redis.Nil {
		fmt.Printf("name not exist")
	} else if err != nil {
		fmt.Printf("get name failed, err%v\n", err)
		return
	} else {
		fmt.Println("Socre", val2)
	}

}

func hgetDemo() {
	v, err := rdb.HGetAll("user").Result()
	if err != nil {
		fmt.Printf("hgetall failed, err:%v\n", err)
		return
	}

	fmt.Println(v)

	v2 := rdb.HMGet("user", "name", "age").Val()
	fmt.Println(v2)

	v3 := rdb.HGet("user", "age").Val()
	fmt.Println(v3)
}

func redisExamole2() {
	zsetKey := "language_rank"
	languages := []redis.Z{
		redis.Z{Score: 90.0, Member: "golang"},
		redis.Z{Score: 98.0, Member: "java"},
		redis.Z{Score: 95.0, Member: "python"},
		redis.Z{Score: 97.0, Member: "javaScript"},
		redis.Z{Score: 99.0, Member: "C/C++"},
	}

	//ZADD
	num, err := rdb.ZAdd(zsetKey, languages...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Printf("Zadd %d success.\n", num)

	// 把Golang的分数加10
	newScore, err := rdb.ZIncrBy(zsetKey, 10.0, "golang").Result()
	if err != nil {
		fmt.Printf("ZIncrBy failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	//取分数最高的三个
	ret, err := rdb.ZRevRangeWithScores(zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("ZRevRangeWithScores failed, err:%v\n", err)
		return
	}
	for _, z := range ret{
		fmt.Println(z.Member, z.Score)
	}

	// 取95-100 的
	op := redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(zsetKey, op).Result()
	if err != nil {
		fmt.Printf("ZRangeByScoreWithScores failed, err:%v\n", err)
		return
	}
	for _, z := range ret{
		fmt.Println(z.Member, z.Score)
	}
}

func watchDemo(){
	//监视watch_count的值 并在值不变的前提下+1
	key := "watch_count"
	err := rdb.Watch(
		func(tx *redis.Tx) error {
			n, err := tx.Get(key).Int()
			if err != nil && err != redis.Nil{
				return err
			}
 
			_, err = tx.Pipelined(
				func(p redis.Pipeliner) error {
					//业务逻辑
					p.Set(key, n+1, 0,)
					return nil
				})
			return err
		}, key)
	if err != nil{
		fmt.Printf("tx exec failed, err:%v\n", err)
		return
	}

	fmt.Printf("tx exec success")

}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init failed err:%v\n", err)
		return
	}
	fmt.Printf("init succsess ...")
	defer rdb.Close()

	//redisExamole()
	//hgetDemo()
	redisExamole2()
}
