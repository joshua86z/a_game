package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Pay struct {
}

func (this *Pay) BuyCoinRequest(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	request := &protodata.BuyCoinRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return lineNum(), nil, err
	}

	index := int(request.GetProductIndex())
	configList := models.ConfigCoinDiamondList()
	needDiamond := configList[index-1].Diamond
	addCoin := configList[index-1].Coin

	if needDiamond > RoleModel.Diamond {
		return lineNum(), nil, fmt.Errorf("钻石不足")
	}

	oldCoin := RoleModel.Coin
	RoleModel.Coin += addCoin

	err := RoleModel.SubDiamond(needDiamond, models.BUY_COIN, fmt.Sprintf("index: %d , coin:%d -> %d , diamond:%d -> %d", index, oldCoin, RoleModel.Coin, RoleModel.Diamond, RoleModel.Diamond-needDiamond))
	if err != nil {
		return lineNum(), nil, err
	}

	response := &protodata.BuyCoinResponse{
		Role: roleProto(RoleModel),
		Coin: proto.Int32(int32(addCoin)),
	}
	return protodata.StatusCode_OK, response, nil
}
