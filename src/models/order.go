package models

import (
//	"fmt"
//	"math/rand"
//	"time"
)

var Order OrderModel

type OrderModel struct {
}

type OrderData struct {
	OrderId    string `db:"order_id"`
	PlatId     int    `db:"plat_id"`
	Uid        int64  `db:"uid"`
	Money      int    `db:"order_money"`
	Diamond    int    `db:"order_diamond"`
	Status     int    `db:"order_status"` //0创建订单1付款成功2划账到游戏
	CreateTime int64  `db:"order_createtime"`
	OkTime     int64  `db:"order_oktime"`
	Mark       string `db:"order_mark"`
}

func init() {
	DB().AddTableWithName(OrderData{}, "pay_orders").SetKeys(false, "OrderId")
}

//// 判断充值次数
//func OrderCount(uniqueId int64) (int, error) {
//	count, err := PayDB().SelectInt("SELECT COUNT(order_id) FROM g_pay_orders WHERE server_id = ? AND roles_unique = ? AND order_status = 2 ", ServerId, uniqueId)
//	return int(count), err
//}

//func (this OrderModel) Insert(uid int64, rmb int, diamond int) (string, error) {
//
//	now := time.Now()
//	orderId := fmt.Sprintf("101%s%05d", now.Format("200601021504"), rand.Intn(100000))
//
//	order := OrderData{
//		OrderId:    orderId,
//		Uid:        uid,
//		Rmb:        rmb,
//		Diamond:    diamond,
//		CreateTime: now.Unix(),
//	}
//
//	return orderId, DB().Insert(&order)
//}

func (this OrderModel) Order(orderId string) (*OrderData, error) {
	OrderData := new(OrderData)
	err := DB().SelectOne(OrderData, "SELECT * FROM pay_orders WHERE order_id = ?", orderId)
	return OrderData, err
}

//func getOrderLockKey(orderId string) string {
//	return fmt.Sprintf("ORDERLOCK_%s", orderId)
//}
//
//func (this OrderModel) Confirm(RoleData *RoleData) error {
//
//	if this.Status != 1 {
//		return fmt.Errorf("Order status not is '1' ")
//	}
//
//	Transaction, err := DB().Begin()
//	if err != nil {
//		return err
//	}
//
//	this.Status = 2
//
//	affected_rows, err := Transaction.Update(&this)
//	if err != nil || affected_rows != 1 {
//		Transaction.Rollback()
//		return err
//	}
//
//	oldDiamond := RoleData.Diamond
//	RoleData.Diamond += this.Diamond
//	RoleData.UnixTime = time.Now().Unix()
//
//	_, err = Transaction.Update(RoleData)
//	if err != nil {
//		RoleData.Diamond -= this.Diamond
//		Transaction.Rollback()
//		return err
//	}
//
//	if Transaction.Commit() == nil {
//		return err
//	} else {
//		RoleData.Diamond -= this.Diamond
//	}
//
//	InsertSubDiamondFinanceLog(this.Uid, FINANCE_BUY_DIAMOND, oldDiamond, RoleData.Diamond, fmt.Sprintf("diamond: %d -> %d", oldDiamond, RoleData.Diamond))
//
//	return nil
//}
