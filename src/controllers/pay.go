package controllers

import (
	"models"
	"protodata"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

type Pay struct {

}

func (this *Login) List(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	request := &protodata.LoginRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return ReturnStr(commandRequest, 19, ""), fmt.Errorf("%v", err)
	}

	payCenterList := models.ConfigPayCenterList()

	var result []protodata.CoinProductData
	for index, val := range payCenterList {
		var temp protodata.CoinProductData
		temp.ProductIndex = proto.Int32(int32(index) + 1)
		temp.ProductName = proto.String(val.Name)
		temp.ProductDesc = proto.String(val.Name)
		temp.PriceDiamond = proto.Int32(int32(val.Diamond))
		temp.ProductCoin = proto.Int32(int32(val.Rmb))
		result = append(result, temp)
	}

	fmt.Println(result)

	response := &protodata.LoginResponse{}
	return ReturnStr(commandRequest, 1, response), nil
}


//
//type gamePay struct{}
//
//var Pay gamePay
//
//func (this gamePay) ConfirmOrder(orderId string) error {
//
//	log.Debug("GamePay call. %v", orderId)
//
//	payOrder, err := models.GetOrderById(orderId)
//	if err != nil {
//		return err
//	}
//
//	if err = models.OrderConfirm(orderId); err != nil {
//		return err
//	}
//
//	var firstGift bool // 判断首充礼包
//
//	if count, err := models.OrderCount(payOrder.Unique); err != nil {
//		return err
//	} else if count == 1 {
//		firstGift = true
//	}
//
//	if firstGift {
//		models.InsertItem(payOrder.Unique, configs.ConfigItemById(120), 1)
//
//		// 更新道具类的成就
//		models.Achievement(payOrder.Unique, 12, 13, 24)
//	}
//
//	return nil
//}
