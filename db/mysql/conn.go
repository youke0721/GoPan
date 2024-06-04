package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init() 函数是 Go 语言中的一个特殊函数,它会在程序启动时自动执行
func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(192.168.222.128)/fileserver?charset=uft8")
	db.SetMaxOpenConns(1000)
	err := db.ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn:返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
