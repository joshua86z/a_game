package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) BuyGeneral() error {

	request := new(protodata.BuyGeneralRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	configId := int(request.GetGeneralId())
	config := models.ConfigGeneralMap()[configId]
	needDiamond := config.BuyDiamond

	if models.General.General(this.Uid, configId) != nil {
		return this.Send(lineNum(), fmt.Errorf("已有这个英雄: %d", configId))
	}

	if needDiamond > this.Role.Diamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	if err := this.Role.SubDiamond(needDiamond, models.FINANCE_BUY_GENERAL, fmt.Sprintf("genId : %d", configId)); err != nil {
		return this.Send(lineNum(), err)
	}

	general := models.General.Insert(this.Uid, config)
	if general == nil {
		return this.Send(lineNum(), fmt.Errorf("失败:数据库错误"))
	}

	response := &protodata.BuyGeneralResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, config),
	}
	return this.Send(StatusOK, response)
}

func (this *Connect) GeneralLevelUp() error {

	request := new(protodata.GeneralLevelUpRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}
	configId := int(request.GetGeneralId())

	general := models.General.General(this.Uid, configId)
	if general == nil {
		return this.Send(lineNum(), fmt.Errorf("英雄数据出错"))
	}

	config := models.ConfigGeneralMap()[general.ConfigId]

	if general.Level >= len(config.LevelUpCoin)-1 {
		return this.Send(lineNum(), fmt.Errorf("英雄已是最高等级"))
	}

	coin := config.LevelUpCoin[general.Level]
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	if err := this.Role.SubCoin(coin, models.FINANCE_GENERAL_LEVELUP, config.Name); err != nil {
		return this.Send(lineNum(), err)
	}

	if err := general.LevelUp(config); err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.GeneralLevelUpResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, config),
	}

	return this.Send(StatusOK, response)
}

func (this *Connect) SetLeader() error {

	request := new(protodata.SetLeaderRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	generalId := int(request.GetLeaderId())

	var find bool
	for _, general := range models.General.List(this.Uid) {
		if general.ConfigId == generalId {
			find = true
			break
		}
	}

	if !find {
		return this.Send(lineNum(), fmt.Errorf("武将没有解锁"))
	}
	if err := this.Role.SetGeneralConfigId(generalId); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.SetLeaderResponse{LeaderId: request.LeaderId})
}

func generalProtoList(generalList []*models.GeneralData, configs map[int]*models.ConfigGeneral) []*protodata.GeneralData {

	var result []*protodata.GeneralData
	for _, config := range configs {

		var find bool
		for _, general := range generalList {
			if general.ConfigId == config.ConfigId {
				result = append(result, generalProto(general, config))
				find = true
				break
			}
		}
		if !find {
			result = append(result, generalProto(new(models.GeneralData), config))
		}
	}
	return result
}

func generalProto(general *models.GeneralData, config *models.ConfigGeneral) *protodata.GeneralData {

	var use bool
	if general.Id == 0 {
		general.Atk = config.Atk
		general.Def = config.Def
		general.Hp = config.Hp
		general.Speed = config.Speed
		general.Dex = config.Dex
		general.Range = config.Range

	} else {
		use = true
	}

	return &protodata.GeneralData{
		GeneralId:   proto.Int32(int32(config.ConfigId)),
		GeneralName: proto.String(config.Name),
		GeneralDesc: proto.String(config.Desc),
		Level:       proto.Int32(int32(general.Level)),
		Atk:         proto.Int32(int32(general.Atk)),
		Def:         proto.Int32(int32(general.Def)),
		Hp:          proto.Int32(int32(general.Hp)),
		Speed:       proto.Int32(int32(general.Speed)),
		TriggerR:    proto.Int32(int32(general.Range)),
		AtkR:        proto.Int32(int32(config.AtkRange)),
		GeneralType: proto.Int32(int32(config.Type)),
		LevelUpCoin: proto.Int32(int32(config.LevelUpCoin[general.Level])),
		BuyDiamond:  proto.Int32(int32(config.BuyDiamond)),
		KillNum:     proto.Int32(int32(general.KillNum)),
		IsUnlock:    proto.Bool(use)}
}
