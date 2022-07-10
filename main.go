package main

import (
	"fmt"
	"net/http"

	"gowebsocket/lib/redislib"
	"gowebsocket/routers"
	"gowebsocket/servers/websocket"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	initConfig()
	// 初始化redis
	initRedis()

	router := gin.Default()
	// 初始化路由
	routers.InitRouter(router)
	routers.WebsocketInit()

	go websocket.StartWebSocket()

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath("./") // 添加搜索路径
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// fmt.Println("config app:", viper.Get("app"))
	// fmt.Println("config redis:", viper.Get("redis"))
}

func initRedis() {
	redislib.ExampleNewClient()
}
