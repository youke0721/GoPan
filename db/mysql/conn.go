package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init() 函数是 Go 语言中的一个特殊函数,它会在程序启动时自动执行
func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(192.168.222.128:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	//init() 函数是 Go 语言中的一个特殊函数,它会在程序启动时自动执行,一行代码三行erro
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn:返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
