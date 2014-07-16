package models

import (
	"fmt"
	"math/rand"
	"time"
)

type OrderModel struct {
	OrderId    string `db:"order_id"`
	Uid        int64  `db:"roles_unique"`
	Rmb        int    `db:"order_rmb"`
	Diamond    int    `db:"order_diamond"`
	Status     int    `db:"order_status"` //0创建订单1付款成功2划账到游戏
	CreateTime int64  `db:"order_createtime"`
	OkTime     int64  `db:"order_oktime"`
	Mark       string `db:"order_mark"`
}

func init() {
	DB().AddTableWithName(OrderModel{}, "pay_orders").SetKeys(false, "OrderId")
}

//// 判断充值次数
//func OrderCount(uniqueId int64) (int, error) {
//	count, err := PayDB().SelectInt("SELECT COUNT(order_id) FROM g_pay_orders WHERE server_id = ? AND roles_unique = ? AND order_status = 2 ", ServerId, uniqueId)
//	return int(count), err
//}

func InsertOrder(uid int64, rmb int, diamond int) string {

	now := time.Now()

	orderId := fmt.Sprintf("%3d%s%05d", 101, now.Format("200601021504"), rand.Intn(100000))

	order := OrderModel{
		OrderId:    orderId,
		Uid:        uid,
		Rmb:        rmb,
		Diamond:    diamond,
		CreateTime: now.Unix(),
	}

	if err := DB().Insert(&order); err != nil {
		DBError(err)
	}

	return orderId
}

//func (order *PayOrder) Create() error {
//	order.Status = 0
//	order.CreateTime = int(time.Now().Unix())
//	return dbMap.Insert(order)
//}

func NewOrderModel(orderId string) (OrderModel, error) {

	var OrderModel OrderModel

	if err := DB().SelectOne(&OrderModel, "SELECT * FROM pay_orders WHERE order_id = ?", orderId); err != nil {
		return OrderModel, err
	}

	return OrderModel, nil
}

//func getOrderLockKey(orderId string) string {
//	return fmt.Sprintf("ORDERLOCK_%s", orderId)
//}

func (this OrderModel) Confirm() error {

	if this.Status != 1 {
		return fmt.Errorf("Order status not is '1' ")
	}

	Transaction, err := DB().Begin()
	if err != nil {
		return err
	}

	this.Status = 2

	affected_rows, err := Transaction.Update(&this)
	if err != nil || affected_rows != 1 {
		Transaction.Rollback()
		return err
	}

	RoleModel := NewRoleModel(this.Uid)

	oldDiamond := RoleModel.Diamond
	RoleModel.Diamond += this.Diamond
	RoleModel.UnixTime = time.Now().Unix()

	_, err = Transaction.Update(RoleModel)
	if err != nil {
		Transaction.Rollback()
		return err
	}

	if Transaction.Commit() == nil {
		Transaction.Rollback()
		return err
	}

	InsertSubDiamondFinanceLog(this.Uid, A, oldDiamond, RoleModel.Diamond, "")

	return nil
}
