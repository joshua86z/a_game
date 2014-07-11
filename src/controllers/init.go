package controllers

import (
	"libs/log"
	"protodata"
	"reflect"
	"runtime"
	_ "models"
)

// 游戏主逻辑
var (
	login *Login
)

func init() {

	login = &Login{}

	handlers = getHandlerMap()
	handlerNames = make(map[int32]string)
	for index, handler := range handlers {
		handlerNames[index] = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	}

	log.Info("Program Run !")
}

// handlers
var (
	handlers     map[int32]func(int64, *protodata.CommandRequest) (string, error)
	handlerNames map[int32]string
)

// get command handlers for set handlers
func getHandlerMap() map[int32]func(int64, *protodata.CommandRequest) (string, error) {

	return map[int32]func(int64, *protodata.CommandRequest) (string, error){
		10000: login.Login,
	}
}
