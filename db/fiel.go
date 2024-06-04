package db

import (
	mydb "Gopan/db/mysql"
	"log"
)

// OnFileUploadFinished : 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		//以提高查询效率,SQL 语句只需要编译一次,执行预备好的 SQL 语句，防止sql注入
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		//记录严重错误信息并立即终止程序执行。
		log.Fatal("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	//接受一个可变长度的参数列表,这些参数将被用来替换 SQL 语句中的占位符(通常是问号 ?)。
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		//log.Println() 用于记录一些可恢复的错误或者一些有助于调试的信息。它不会像 log.Fatal() 那样导致程序立即退出,而是让程序继续运行,尽管可能会有一些副作用。
		log.Println(err.Error())
		return false
	}
	//获取了受影响的行数 (rf) 以及是否有错误发生 (err).
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Fatalf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	//如果在上传过程中发生了其他错误,则返回 false 表示上传失败。
	return false
}

//额外：if nil == err: 如果没有错误发生,则进入下一步检查。if rf <= 0: 如果受影响的行数小于等于 0,意味着没有新的数据被插入。这种情况下,可能是因为文件之前已经被上传过了。
