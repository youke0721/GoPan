package handler

import (
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
	}
}
