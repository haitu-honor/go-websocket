package websocket

import (
	"fmt"
	"gowebsocket/models"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultAppId = 101 // 默认平台Id
)

var (
	clientManager = NewClientManager()                    // 管理者
	appIds        = []uint32{defaultAppId, 102, 103, 104} // 全部的平台

	serverIp   string
	serverPort string
)

func GetAppIds() []uint32 {

	return appIds
}

// 启动服务监听
func StartWebSocket() {
	http.HandleFunc("/acc", wsPage)
	http.ListenAndServe(":8089", nil)
}

// 协议升级
func wsPage(w http.ResponseWriter, req *http.Request) {
	// 升级连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
			return true
		},
	}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	fmt.Println("websocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()
	// 用户连接事件
	clientManager.Register <- client
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}

	return
}

func InAppIds(appId uint32) (inAppId bool) {

	for _, value := range appIds {
		if value == appId {
			inAppId = true

			return
		}
	}

	return
}
