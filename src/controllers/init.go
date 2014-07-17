package controllers

import (
	"libs/log"
	"models"
	_ "models"
	"protodata"
	"reflect"
	"runtime"
)

// 游戏主逻辑
var (
	login   *Login
	role    *Role
	general *General
	item    *Item
	pay     *Pay
	mail    *Mail
)

func init() {

	handlers = getHandlerMap()
	handlerNames = make(map[int32]string)
	for index, handler := range handlers {
		handlerNames[index] = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	}

	log.Info("Program Run !")
}

// handlers
var (
	handlers     map[int32]func(*models.RoleModel, *protodata.CommandRequest) (protodata.StatusCode, interface{}, error)
	handlerNames map[int32]string
)

// get command handlers for set handlers
func getHandlerMap() map[int32]func(*models.RoleModel, *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	return map[int32]func(*models.RoleModel, *protodata.CommandRequest) (protodata.StatusCode, interface{}, error){
		10103: login.Login,
		10104: role.UserDataRequest,
		10105: role.SetRoleName,
		10106: role.RandomName,
		10107: role.BuyStaminaRequest,
		10108: pay.BuyCoinRequest,
		//10109			//补充钻石
		10110: general.LevelUp,
		10111: general.Buy,
		10112: mail.List,
		10113: mail.MailRewardRequest,
		10114: item.LevelUp,
		//10115			//战斗初始化
		//10116			//战斗结束
	}
}

func lineNum() protodata.StatusCode {
	_, _, line, ok := runtime.Caller(1)
	if ok {
		return protodata.StatusCode(line)
	}
	return -1
}
