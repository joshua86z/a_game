package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

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

	oldCoin := this.Role.Coin
	this.Role.Coin += addCoin

	err := this.Role.SubDiamond(needDiamond, models.BUY_COIN, fmt.Sprintf("index: %d , coin:%d -> %d , diamond:%d -> %d", index, oldCoin, this.Role.Coin, this.Role.Diamond, this.Role.Diamond-needDiamond))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.BuyCoinResponse{
		Role: roleProto(this.Role),
		Coin: proto.Int32(int32(addCoin)),
	}
	return this.Send(protodata.StatusCode_OK, response)
}
