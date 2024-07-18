package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeChatMessage represents the structure of a WeChat message
type WeChatMessage struct {
	ToUserName   string  `json:"ToUserName"`
	FromUserName string  `json:"FromUserName"`
	CreateTime   float64 `json:"CreateTime"`
	MsgType      string  `json:"MsgType"`
	Content      string  `json:"Content"`
	MsgId        float64 `json:"MsgId"`
}

// WeChatResponse represents the structure of a WeChat response
type WeChatResponse struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
}

// WeChatMsgHandler handles WeChat messages
func WeChatMsgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg WeChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	response := WeChatResponse{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "您发送的消息是: " + msg.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received message: %+v\n", msg)
}
