package controllers

import (
	"fmt"
	"libs/log"
	"libs/lua"
	"models"
	"net/http"
	_ "net/http/pprof"
)

func init() {

	go func() {
		Lua, _ := lua.NewLua("conf/app.lua")
		port := Lua.GetInt("httpPort")
		Lua.Close()
		http.HandleFunc("/confirm", httpConfirm)
//		http.HandleFunc("/adddiamond", httpAddDiamond)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}

func getPlayer(uid int64) (*Connect, bool) {
	playLock.Lock()
	conn, ok := playerMap[uid]
	playLock.Unlock()
	return conn, ok
}

func httpConfirm(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Critical("RPC Panic: %v", err)
		}
	}()

	orderId := r.FormValue("order_id")
	if orderId == "" {
		return
	}

	OrderModel, err := models.NewOrderModel(r.FormValue("order_id"))
	if err != nil {
		return
	}

	conn, ok := getPlayer(OrderModel.Uid)

	var RoleData *models.RoleData
	if ok && conn.Role != nil {
		RoleData = conn.Role
	} else {
		RoleData = models.Role.Role(OrderModel.Uid)
	}

	err = OrderModel.Confirm(RoleData)
	if ok {
		var code int = StatusOK
		var result interface{}
		if err != nil {
			code = lineNum()
			result = err
		}
		conn.Send(code, result)
	}
}
