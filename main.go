package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	// http.HandleFunc("/find", service.FindHandler) // 新增 find 接口
	http.HandleFunc("/find", service.WeChatMsgHandler)
	// 增加一个接口/Wconfig2Ck 传入参数为jsonString 一个json字符串
	// 接收参数后调用service.Wconfig2Ck(jsonString) 会生成一个cookie json字符串
	http.HandleFunc("/Wconfig2Ck", service.Wconfig2CkHandler) // 新增接口
	log.Fatal(http.ListenAndServe(":80", nil))
}
