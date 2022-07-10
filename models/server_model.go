package models

import (
	"errors"
	"strings"
)

type Server struct {
	Ip   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

func StringToServer(str string) (server *Server, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {

		return nil, errors.New("err")
	}

	server = &Server{
		Ip:   list[0],
		Port: list[1],
	}

	return
}
