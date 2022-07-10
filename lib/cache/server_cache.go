package cache

import (
	"encoding/json"
	"fmt"
	"gowebsocket/lib/redislib"
	"gowebsocket/models"
	"strconv"
)

const (
	serversHashKey       = "acc:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60        // key过期时间
	serversHashTimeout   = 3 * 60             // 超时时间
)

func getServersHashKey() (key string) {
	key = serversHashKey

	return
}

func GetServerAll(currentTime uint64) (servers []*models.Server, err error) {
	servers = make([]*models.Server, 0)
	key := getServersHashKey()

	redisClient := redislib.GetClient()

	val, err := redisClient.Do("hGetAll", key).Result()

	valByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valByte))

	serverMap, err := redisClient.HGetAll(key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)

		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		// 超时
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}

		server, err := models.StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
