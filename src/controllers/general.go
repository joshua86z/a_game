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

	baseId := int(request.GetGeneralId())
	baseGeneral := models.BaseGeneral(baseId, nil)
	needDiamond := baseGeneral.BuyDiamond

	if models.General.General(this.Uid, baseId) != nil {
		return this.Send(lineNum(), fmt.Errorf("已有这个英雄: %d", baseId))
	}

	if needDiamond > this.Role.Diamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	if err := this.Role.SubDiamond(needDiamond, models.FINANCE_BUY_GENERAL, fmt.Sprintf("genId : %d", baseId)); err != nil {
		return this.Send(lineNum(), err)
	}

	general, err := models.General.Insert(this.Uid, baseGeneral)
	if err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.BuyGeneralResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, baseGeneral),
	}
	return this.Send(StatusOK, response)
}

func (this *Connect) GeneralLevelUp() error {

	request := new(protodata.GeneralLevelUpRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}
	baseId := int(request.GetGeneralId())

	general := models.General.General(this.Uid, baseId)
	if general == nil {
		return this.Send(lineNum(), fmt.Errorf("英雄数据出错"))
	}

	baseGeneral := models.BaseGeneral(general.BaseId, nil)

	if general.Level >= len(baseGeneral.LevelUpCoin)-1 {
		return this.Send(lineNum(), fmt.Errorf("英雄已是最高等级"))
	}

	coin := baseGeneral.LevelUpCoin[general.Level]
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	if err := this.Role.SubCoin(coin, models.FINANCE_GENERAL_LEVELUP, baseGeneral.Name); err != nil {
		return this.Send(lineNum(), err)
	}

	if err := general.LevelUp(baseGeneral); err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.GeneralLevelUpResponse{
		Role:    roleProto(this.Role),
		General: generalProto(general, baseGeneral),
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
		if general.BaseId == generalId {
			find = true
			break
		}
	}

	if !find {
		return this.Send(lineNum(), fmt.Errorf("武将没有解锁"))
	}
	if err := this.Role.SetGeneralBaseId(generalId); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.SetLeaderResponse{LeaderId: request.LeaderId})
}

func generalProtoList(generalList []*models.GeneralData, configs map[int]*models.Base_General) []*protodata.GeneralData {

	var result []*protodata.GeneralData
	for _, baseGeneral := range configs {

		var find bool
		for _, general := range generalList {
			if general.BaseId == baseGeneral.BaseId {
				result = append(result, generalProto(general, baseGeneral))
				find = true
				break
			}
		}
		if !find {
			result = append(result, generalProto(new(models.GeneralData), baseGeneral))
		}
	}
	return result
}

func generalProto(general *models.GeneralData, baseGeneral *models.Base_General) *protodata.GeneralData {

	var unlock bool
	if general.BaseId == 0 {
		general.Atk = baseGeneral.Atk
		general.Def = baseGeneral.Def
		general.Hp = baseGeneral.Hp
		general.Speed = baseGeneral.Speed
		general.Dex = baseGeneral.Dex
		general.Range = baseGeneral.Range

	} else {
		unlock = true
	}

	return &protodata.GeneralData{
		GeneralId:   proto.Int32(int32(baseGeneral.BaseId)),
		GeneralName: proto.String(baseGeneral.Name),
		GeneralDesc: proto.String(baseGeneral.Desc),
		Level:       proto.Int32(int32(general.Level)),
		Atk:         proto.Int32(int32(general.Atk)),
		Def:         proto.Int32(int32(general.Def)),
		Hp:          proto.Int32(int32(general.Hp)),
		Speed:       proto.Int32(int32(general.Speed)),
		Dex:         proto.Int32(int32(general.Dex)),
		TriggerR:    proto.Int32(int32(general.Range)),
		AtkR:        proto.Int32(int32(baseGeneral.AtkRange)),
		SkillAtk:    proto.Int32(int32(baseGeneral.SkillAtk)),
		GeneralType: proto.Int32(int32(baseGeneral.Type)),
		LevelUpCoin: proto.Int32(int32(baseGeneral.LevelUpCoin[general.Level])),
		BuyDiamond:  proto.Int32(int32(baseGeneral.BuyDiamond)),
		KillNum:     proto.Int32(int32(general.KillNum)),
		IsUnlock:    proto.Bool(unlock)}
}
