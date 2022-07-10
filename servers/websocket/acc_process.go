package websocket

import (
	"encoding/json"
	"fmt"
	"gowebsocket/common"
	"gowebsocket/models"
	"sync"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value
	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	value, ok = handlers[key]
	return
}

// 处理数据
func ProcessData(client *Client, message []byte) {
	fmt.Println("处理数据", client.Addr, message)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("处理数据 stop", err)
		}
	}()

	request := &models.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		client.SendMsg([]byte("数据不合法"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		client.SendMsg([]byte("处理数据失败"))
		return
	}

	seq := request.Seq
	cmd := request.Cmd

	var (
		code uint32
		msg  string
		data interface{}
	)

	// request
	fmt.Println("acc_request", cmd, client.Addr)
	// 采用 map 注册的方式
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		fmt.Println("处理数据 路由不存在", client.Addr, "cmd", cmd)
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := models.NewResponseHead(seq, cmd, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		return
	}
	client.SendMsg(headByte)
	fmt.Println("acc_response send", client.Addr, client.AppId, client.UserId, "cmd", cmd, "code", code)

	return
}
