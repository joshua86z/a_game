package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"math/rand"
	"models"
	"protodata"
	"time"
)

func (this *Connect) BuyDiamondRequest() error {

	request := &protodata.BuyDiamondRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	index := int(request.GetProductIndex())

	payCenter := models.ConfigPayCenterList()
	product := payCenter[index-1]

	now := time.Now()
	orderId := fmt.Sprintf("101%s%05d", now.Format("200601021504"), rand.Intn(100000))

	user, err := models.User.User(this.Uid)
	if err != nil {
		return this.Send(lineNum(), err)
	}

	order := &models.OrderData{
		OrderId:    orderId,
		PlatId:     user.PlatId,
		Uid:        this.Uid,
		Money:      product.Money,
		Diamond:    product.Diamond,
		CreateTime: now.Unix()}

	if err := models.DB().Insert(order); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.BuyDiamondResponse{OrderId: proto.String(orderId)})
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
