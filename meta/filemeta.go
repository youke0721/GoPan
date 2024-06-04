package meta

import (
	mydb "Gopan/db"
)

// FileMeta: 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

//文件的元信息，包括文件的 SHA1 值、文件名、文件大小、存储位置和上传时间等字段。

// fileMetas 是一个全局变量，用于存储文件元信息。它是一个 map 类型，键是文件的 SHA1 值，值是对应的 FileMeta 结构体。
// 使用 SHA1 值作为键的原因主要是确保键的唯一性和文件的唯一性
var fileMetas map[string]FileMeta

// init() 函数在包被导入时自动执行，这里用于初始化 fileMetas 变量，确保在使用时已经被正确初始化。
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 函数用于新增或更新文件的元信息。它接受一个 FileMeta 结构体作为参数，然后将其存储到 fileMetas map 中，以文件的 SHA1 值作为键。
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMeta: 新增/更新文件元信息到数据库
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMeta 函数用于通过文件的 SHA1 值获取对应的文件元信息对象。它接受一个 SHA1 值作为参数，并从 fileMetas map 中查找对应的元信息，然后返回。
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB : 从MySQL获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) { //对象的指针才能表示为nil
	tfile, err := mydb.GetFileMeta(fileSha1)
	if tfile == nil || err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash, //需要秋木耳数据类型的转换
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil
}

//TableFile结构体:
//这个结构体可能是直接与数据库表tbl_file对应的结构体,用于存储从数据库查询到的原始数据。
//它的字段可能与数据库表的列名一一对应,比如FileHash、FileAddr、FileName和FileSize。
//FileMeta结构体:
//这个结构体可能代表了应用程序层面的文件元信息,字段名可能与TableFile有所不同,更贴近业务逻辑。
//比如这里的FileSha1、FileName、FileSize和Location。

// 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
