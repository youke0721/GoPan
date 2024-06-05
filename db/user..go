package db

import (
	//使用别名的好处是,当你需要引用包中的成员时
	//可以直接使用别名,而不需要写出完整的包路径,从而使代码更加简洁明了。
	mydb "Gopan/db/mysql"
	"fmt"
)

// UserSignup : 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	//获取执行 SQL 语句后受影响的行数,存储在 rowsAffected 变量中。
	//如果执行过程中没有错误,且受影响的行数大于 0,则说明插入成功,返回 true。
	//使用准备好的 SQL 语句执行插入操作,并将结果存储在 ret 变量中
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
