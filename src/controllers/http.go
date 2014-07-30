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

	var RoleModel *models.RoleModel
	if ok && conn.Role != nil {
		RoleModel = conn.Role
	} else {
		RoleModel = models.NewRoleModel(OrderModel.Uid)
	}

	err = OrderModel.Confirm(RoleModel)
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

func httpAddDiamond(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Critical("RPC Panic: %v", err)
		}
	}()

	uidStr, pwd, num := r.FormValue("uid"), r.FormValue("pwd"), r.FormValue("num")
	if uidStr == "" || num == "" || pwd != "joshua" {
		return
	}

	uid := int64(models.Atoi(uidStr))

	conn, ok := getPlayer(uid)

	var RoleModel *models.RoleModel
	if ok && conn.Role != nil {
		RoleModel = conn.Role
	} else {
		RoleModel = models.NewRoleModel(uid)
	}

	RoleModel.AddDiamond(models.Atoi(num), models.FINANCE_ADMIN, "")
}
