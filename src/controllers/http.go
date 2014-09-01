package controllers

import (
	//	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/log"
	"libs/lua"
	"models"
	"net/http"
	_ "net/http/pprof"
	"protodata"
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
		RoleData, err = models.Role.Role(OrderModel.Uid)
		if err != nil {
			return
		}
	}

	if err = OrderModel.Confirm(RoleData); err != nil {
		log.Warn("PAY FAIL order_id: %s , %v", orderId, err)
	}

	if conn != nil {
		*conn.Request.CmdId = 10121
		*conn.Request.CmdIndex = 10121
		conn.Send(StatusOK, &protodata.PaySuccessResponse{Role: roleProto(RoleData)})
	}
}
