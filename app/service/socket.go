package service

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[string]*websocket.Conn)

func Socket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	// 设置读取超时时间为 60 秒
	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	var username string

	// 读取客户端发送的第一条消息，获取客户端的用户名
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println("ReadMessage:", err)
		return
	}
	username = string(p)

	// 存储客户端连接
	if clients[username] == nil {
		clients[username] = conn
	}

	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()

		// 判断是否为超时错误
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("WebSocket connection closed due to timeout.")
			conn.Close()
			delete(clients, username) // 删除断开连接的客户端
			break
		}

		if err != nil {
			log.Println("ReadMessage:", err)
			delete(clients, username) // 删除断开连接的客户端
			return
		}

		// 解析客户端发送的消息，提取目标客户端和消息内容
		var target string
		var message string
		fmt.Sscanf(string(p), "%s %s", &target, &message)

		// 查找目标客户端的连接
		targetConn, ok := clients[target]
		if !ok {
			log.Println("Client not found:", target)
			continue
		}

		// 发送消息给目标客户端
		err = targetConn.WriteMessage(messageType, []byte(message))
		if err != nil {
			log.Println("WriteMessage:", err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	}
}
