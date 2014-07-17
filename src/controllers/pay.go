package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Pay struct {
}

func (this *Pay) BuyCoinRequest(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	request := &protodata.BuyCoinRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return ReturnStr(commandRequest, 19, ""), err
	}

	index := int(request.GetProductIndex())
	configList := models.ConfigCoinDiamondList()
	needDiamond := configList[index-1].Diamond
	addCoin := configList[index-1].Coin

	RoleModel := models.NewRoleModel(uid)
	if needDiamond > RoleModel.Diamond {
		return ReturnStr(commandRequest, 27, "钻石不足"), fmt.Errorf("钻石不足")
	}

	oldCoin := RoleModel.Coin
	RoleModel.Coin += addCoin

	err := RoleModel.SubDiamond(needDiamond, models.BUY_COIN, fmt.Sprintf("index: %d , coin:%d -> %d , diamond:%d -> %d", index, oldCoin, RoleModel.Coin, RoleModel.Diamond, RoleModel.Diamond-needDiamond))
	if err != nil {
		return ReturnStr(commandRequest, 35, "失败:数据库出错"), err
	}

	response := &protodata.BuyCoinResponse{
		Role: roleProto(RoleModel),
		Coin: proto.Int32(int32(addCoin)),
	}
	return ReturnStr(commandRequest, protodata.StatusCode_OK, response), nil
}
