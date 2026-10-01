package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// db 全局变量
var db *sql.DB

type user struct {
	id   int
	name string
	age  int
}

func sqlInit() error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_demo"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open mysql failed: %w", err)
	}

	// 必须 Ping 一次，确认真正连通
	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("ping mysql failed: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return nil
}

func queryRowDemo() {
	sqlstr := "select id, name, age from user where id=?"
	var u user

	//  QueryRow后 一定要scan 释放掉
	err := db.QueryRow(sqlstr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan faild, err:%v\n", err)
		return
	}

	fmt.Printf("id:%d name:%s age%d\n", u.id, u.name, u.age)
}

func queryMultiRowDemo() {
	sqlstr := "select id, name, age from user where id > ?"

	rows, err := db.Query(sqlstr, 0)
	if err != nil {
		fmt.Printf("quwey faild, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("quwey faild, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age%d\n", u.id, u.name, u.age)
	}
	if err = rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return
	}

}

func insertRowDemo() {
	sqlStr := "insert into user(name, age) values(?, ?)"
	ret, err := db.Exec(sqlStr, "edsion", 30)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	var theId int64
	theId, err = ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastedinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theId)
}

func updateRowDemo() {
	sqlStr := "update user set age = ? where id = ?"
	ret, err := db.Exec(sqlStr, 36, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err= %v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d.\n", n)
}

func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 4)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err= %v\n", err)
		return
	}
	fmt.Printf("delete, affected rows:%d.\n", n)
}

// MYSQL语句预处理
func prepareQeryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("Prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("Query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(
			&u.id, &u.name, &u.age,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age%d\n", u.id, u.name, u.age)
	}
	if err = rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return
	}

}

// 预处理插入示例
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("小王子", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	_, err = stmt.Exec("沙河娜扎", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}

// 事务操作示例
// 像取钱时，需要进行两次数据库的更新，必须要都更新或者都不更新，就需要事务 
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}

func main() {
	if err := sqlInit(); err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("mysql connected")
	//queryRowDemo()
	//insertRowDemo()
	queryMultiRowDemo()
	//updateRowDemo()
	deleteRowDemo()
	queryMultiRowDemo()
}
