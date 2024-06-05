package main

import (
	"Gopan/handler"
	"fmt"
	"net/http"
)

func main() {

	// 设定路由规则
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//当用户请求以 /static/ 开头的 URL 时,服务器会去除 /static/ 前缀,然后在 ./static 目录下查找并返回对应的静态文件。
	//用户请求 http://localhost:8080/static/images/logo.png，服务器会在 ./static 目录下查找 images/logo.png 文件并返回。

	//文件
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	//用户
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SigninHandler)
	http.HandleFunc("/user/info", handler.UserInfoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err:%s", err.Error())
	}
}
