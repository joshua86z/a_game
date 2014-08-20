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
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
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

	conn := playerMap.Get(OrderModel.Uid)

	var RoleData *models.RoleData
	if conn != nil && conn.Role != nil {
		RoleData = conn.Role
	} else {
		RoleData = models.Role.Role(OrderModel.Uid)
	}

	err = OrderModel.Confirm(RoleData)
	if conn != nil {
		var code int = StatusOK
		var result interface{}
		if err != nil {
			code = lineNum()
			result = err
		}
		conn.Send(code, result)
	}
}
