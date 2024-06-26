package handler

import (
	"Gopan/meta"
	"Gopan/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的HTML页面
		file, err := os.Open("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "upload failed")
			return
		}
		defer file.Close()

		//读取文件并打印
		data, err := io.ReadAll(file)
		if err != nil {
			io.WriteString(w, "read file failed")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 保存文件存储流到本地目录
		// 添加实际的文件处理逻辑
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		//将文件保存目录参数化，参考create方法的参数改动
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-06 15:04:05"),
		}
		//创建保存文件,修改了保存位置错误的问题
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		//关闭文件句柄
		defer newFile.Close()
		//字节数表示文件的大小，因此将其赋值给 fileMeta.FileSize 字段，这样可以记录上传文件的大小。
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf(" Failed to save data inio file err:%s\n", err.Error())
			return
		}

		//这一行代码将文件指针（读/写位置）移到文件的起始位置，因为在之前将文件内容写入到文件中时，文件指针已经移到了文件的末尾。
		newFile.Seek(0, 0)
		//为了计算文件的 SHA1 值，需要重新将文件指针移到起始位置，以便重新读取文件的内容。
		//计算了上传文件的 SHA1 值
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取文件的哈希值
	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetMeta(filehash)
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	//使用文件哈希值 fsha1 从全局的文件元信息中获取对应的文件元信息 fm。
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	//读取文件的全部内容，并存储在 data 变量
	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//设置响应头的 Content-Type 为 application/octet-stream，表示这是一个二进制文件。
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)

}
