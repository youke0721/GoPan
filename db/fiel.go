package db

import (
	mydb "Gopan/db/mysql"
	"database/sql"
	"fmt"
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

// TableFile : 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	//创建一个空的TableFile结构体实例。以接收后面的查询结果
	tfile := TableFile{}
	//使用之前准备好的SQL语句,并传入filehash参数。QueryRow方法会执行这条SQL语句,并返回一个*sql.Row对象,代表查询结果的第一行数据。
	err = stmt.QueryRow(filehash).Scan( //这一部分将查询结果中的各个列的值,依次赋值给tfile结构体的对应字段。
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		} else {
			//如果出现其他错误,打印错误信息并返回nil和错误。
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
	//指向TableFile结构体的指针   错误信息返回给调用者
}
