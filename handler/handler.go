package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

		data, err := io.ReadAll(file)
		if err != nil {
			io.WriteString(w, "read file failed")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 保存文件存储流到本地目录
		// 添加实际的文件处理逻辑
		file, head, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()
		newFile, err := os.Create("./tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf(" Failed to save data inio file err:%s\n", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}
