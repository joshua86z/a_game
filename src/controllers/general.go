package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) BuyGeneral() error {

	request := &protodata.BuyGeneralRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	configId := int(request.GetGeneralId())
	config := models.ConfigGeneralById(configId)
	needDiamond := config.BuyDiamond

	GeneralModel := models.NewGeneralModel(this.Role.Uid)
	if GeneralModel.General(configId) != nil {
		return this.Send(lineNum(), fmt.Errorf("已有这个英雄: %d", configId))
	}

	if needDiamond > this.Role.Diamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	if err := this.Role.SubDiamond(needDiamond, models.BUY_GENERAL, fmt.Sprintf("genId : %d", configId)); err != nil {
		return this.Send(lineNum(), err)
	}

	general := GeneralModel.Insert(config)
	if general == nil {
		return this.Send(lineNum(), fmt.Errorf("失败:数据库错误"))
	}

	response := &protodata.BuyGeneralResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, config),
	}
	return this.Send(protodata.StatusCode_OK, response)
}

func (this *Connect) GeneralLevelUp() error {

	request := &protodata.GeneralLevelUpRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}
	configId := int(request.GetGeneralId())

	general := models.NewGeneralModel(this.Role.Uid).General(configId)

	config := models.ConfigGeneralById(general.ConfigId)

	coin := levelUpCoin(general.Level)
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足,无法升级"))
	}

	if err := this.Role.SubCoin(coin, models.GENERAL_LEVELUP, config.Name); err != nil {
		return this.Send(lineNum(), err)
	}

	if err := general.LevelUp(config); err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.GeneralLevelUpResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, config),
	}

	return this.Send(protodata.StatusCode_OK, response)
}

func generalProtoList(generalList []*models.GeneralData) []*protodata.GeneralData {

	var result []*protodata.GeneralData

	configs := models.ConfigGeneralMap()
	for _, config := range configs {

		var generalData protodata.GeneralData
		generalData.GeneralId = proto.Int32(int32(config.ConfigId))
		generalData.GeneralName = proto.String(config.Name)
		generalData.GeneralDesc = proto.String(config.Desc)
		generalData.AtkR = proto.Int32(int32(config.AtkRange))
		generalData.GeneralType = proto.Int32(int32(config.Type))
		generalData.BuyDiamond = proto.Int32(int32(config.BuyDiamond))

		var find bool
		for _, general := range generalList {
			if general.ConfigId == config.ConfigId {
				generalData.Level = proto.Int32(int32(general.Level))
				generalData.Atk = proto.Int32(int32(general.Atk))
				generalData.Def = proto.Int32(int32(general.Def))
				generalData.Hp = proto.Int32(int32(general.Hp))
				generalData.Speed = proto.Int32(int32(general.Speed))
				generalData.Dex = proto.Int32(int32(general.Dex))
				generalData.TriggerR = proto.Int32(int32(general.Range))
				generalData.LevelUpCoin = proto.Int32(int32(generallevelUpCoin(general.Level)))
				generalData.IsUnlock = proto.Bool(true)
				generalData.KillNum = proto.Int32(int32(general.KillNum))

				find = true
			}
		}

		if !find {
			generalData.Level = proto.Int32(1)
			generalData.Atk = proto.Int32(int32(config.Atk))
			generalData.Def = proto.Int32(int32(config.Def))
			generalData.Hp = proto.Int32(int32(config.Hp))
			generalData.Speed = proto.Int32(int32(config.Speed))
			generalData.Dex = proto.Int32(int32(config.Dex))
			generalData.TriggerR = proto.Int32(int32(config.Range))
			generalData.LevelUpCoin = proto.Int32(int32(generallevelUpCoin(1)))
			generalData.IsUnlock = proto.Bool(false)
			generalData.KillNum = proto.Int32(0)
		}

		result = append(result, &generalData)
	}
	return result
}

func generallevelUpCoin(level int) int {
	return level
}

func generalProto(general *models.GeneralData, config *models.ConfigGeneral) *protodata.GeneralData {

	var generalData protodata.GeneralData
	generalData.GeneralId = proto.Int32(int32(config.ConfigId))
	generalData.GeneralName = proto.String(config.Name)
	generalData.GeneralDesc = proto.String(config.Desc)
	generalData.Level = proto.Int32(int32(general.Level))
	generalData.Atk = proto.Int32(int32(general.Atk))
	generalData.Def = proto.Int32(int32(general.Def))
	generalData.Hp = proto.Int32(int32(general.Hp))
	generalData.Speed = proto.Int32(int32(general.Speed))
	generalData.TriggerR = proto.Int32(int32(general.Range))
	generalData.AtkR = proto.Int32(int32(config.AtkRange))
	generalData.GeneralType = proto.Int32(int32(config.Type))
	generalData.LevelUpCoin = proto.Int32(int32(general.Level))
	generalData.BuyDiamond = proto.Int32(int32(config.BuyDiamond))
	generalData.KillNum = proto.Int32(int32(general.KillNum))
	generalData.IsUnlock = proto.Bool(true)
	return &generalData
}
