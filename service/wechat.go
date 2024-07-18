package service

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// WeChatMessage represents the structure of a WeChat message
type WeChatMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        int64  `xml:"MsgId"`
}

// WeChatResponse represents the structure of a WeChat response
type WeChatResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

// WeChatMsgHandler handles WeChat messages
func WeChatMsgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg WeChatMessage
	if err := xml.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to parse XML", http.StatusBadRequest)
		return
	}

	response := WeChatResponse{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "您发送的消息是: " + msg.Content,
	}

	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received message: %+v\n", msg)
}
