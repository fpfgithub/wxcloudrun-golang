package service

import (
	"io/ioutil"
	"net/http"
)

func Wconfig2CkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 将请求体转成字符串传递给Wconfig2Ck
	cookieJson, err := Wconfig2Ck(string(body))
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(cookieJson) == 0 {
		http.Error(w, "Api Failed to convert wconfig to cookie", http.StatusInternalServerError)
		return
	}
	// 直接返回 cookieJson 的内容作为响应
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(cookieJson))
}
