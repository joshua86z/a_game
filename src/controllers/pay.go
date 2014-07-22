package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) BuyDiamondRequest() error {

	request := &protodata.BuyDiamondRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	index := int(request.GetProductIndex())

	payCenter := models.ConfigPayCenterList()
	product := payCenter[index-1]

	response := &protodata.BuyDiamondResponse{
		OrderId: proto.String(models.InsertOrder(this.Uid, product.Rmb, product.Diamond)),
	}
	return this.Send(StatusOK, response)
}

func (this *Connect) BuyCoinRequest() error {

	request := &protodata.BuyCoinRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	index := int(request.GetProductIndex())
	configList := models.ConfigCoinDiamondList()
	needDiamond := configList[index-1].Diamond
	addCoin := configList[index-1].Coin

	if needDiamond > this.Role.Diamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	err := this.Role.DiamondIntoCoin(needDiamond, addCoin, fmt.Sprintf("index : %d", index))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.BuyCoinResponse{
		Role: roleProto(this.Role),
		Coin: proto.Int32(int32(addCoin)),
	}
	return this.Send(StatusOK, response)
}
