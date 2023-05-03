package service

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"znews/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Socket upgrade faild",
		}).Error(err.Error())
		return
	}
	defer conn.Close()

	// 设置读取超时时间为 60 秒
	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	//读取客户端发送的第一条消息，获取客户端的用户名
	_, p, err := conn.ReadMessage()
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Socket readMessage faild",
		}).Error(err.Error())
		return
	}
	token := string(p)

	mc, err := middleware.ParseToken(token)
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Socket parseToken faild",
		}).Error(err.Error())
		return
	}
	username := mc.Account

	// 存储客户端连接
	if clients[username] == nil {
		clients[username] = conn
	} else {
		// test
		// username = "test44@gmail.com"
		// clients[username] = conn

		cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "連線已存在")
		if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
			// handle error
		}
		conn.Close()
		return
	}

	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()

		// 判断是否为超时错误
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			conn.Close()
			delete(clients, username) // 删除断开连接的客户端
			break
		}

		if err != nil {
			middleware.Logger().WithFields(logrus.Fields{
				"title": "Socket readMessage faild",
			}).Error(err.Error())
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
			middleware.Logger().WithFields(logrus.Fields{
				"title": "Socket client not found",
			}).Error(target)
			continue
		}

		// 发送消息给目标客户端
		err = targetConn.WriteMessage(messageType, []byte(message))
		if err != nil {
			middleware.Logger().WithFields(logrus.Fields{
				"title": "Socket writeMessage faild",
			}).Error(err.Error())
			continue
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	}
}
