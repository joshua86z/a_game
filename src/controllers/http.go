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
	"time"
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

	OrderData, err := models.Order.Order(orderId)
	if err != nil {
		return
	}

	conn := playerMap.Get(OrderData.Uid)

	var RoleData *models.RoleData
	if conn != nil && conn.Role != nil {
		RoleData = conn.Role
	} else {
		RoleData, err = models.Role.Role(OrderData.Uid)
		if err != nil {
			return
		}
	}

	// ----- 启动事务 -----
	if OrderData.Status != 1 {
		log.Warn("Order status not is '1' ORDERID: %s", orderId)
		return
	}

	Transaction, err := models.DB().Begin()
	if err != nil {
		log.Warn("%v", err)
		return
	}

	OrderData.Status = 2

	affected_rows, err := Transaction.Update(OrderData)
	if err != nil || affected_rows != 1 {
		Transaction.Rollback()
		log.Warn("%v", err)
		return
	}

	oldDiamond := RoleData.Diamond
	RoleData.Diamond += OrderData.Diamond
	RoleData.UnixTime = time.Now().Unix()

	_, err = Transaction.Update(RoleData)
	if err != nil {
		RoleData.Diamond = oldDiamond
		Transaction.Rollback()
		log.Warn("%v", err)
		return
	}

	// ----- 提交 -----
	if Transaction.Commit() != nil {
		RoleData.Diamond = oldDiamond
		Transaction.Rollback()
		log.Warn("%v", err)
		return
	}

	models.InsertSubDiamondFinanceLog(OrderData.Uid, models.FINANCE_BUY_DIAMOND, oldDiamond, RoleData.Diamond, fmt.Sprintf("orderId: %s", orderId))

	if conn != nil {
		*conn.Request.CmdId = 10121
		*conn.Request.CmdIndex = 10121
		conn.Send(StatusOK, &protodata.PaySuccessResponse{Role: roleProto(RoleData)})
	}
}
