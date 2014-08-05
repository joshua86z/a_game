package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) BattleRequest() error {

	request := &protodata.FightInitRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	chapterId := int(request.GetChapterId())
	sectionId := int(request.GetSectionId())
	tempItemList := request.GetTempItems()

	var (
		c, s   int
		config *models.ConfigDuplicate
	)

	configs := models.ConfigDuplicateList()
	if chapterId != 1 || sectionId != 1 {
		for index, val := range configs {
			if val.Chapter > chapterId {
				break
			}
			if val.Chapter == chapterId && val.Section == sectionId {
				if index > 0 {
					index -= 1

				}
				c, s = configs[index].Chapter, configs[index].Section
				config = val
				break
			}
		}
		if c == 0 {
			return this.Send(lineNum(), fmt.Errorf("没有这个副本"))
		}

		dplicateList := models.NewDuplicateModel(this.Uid).List()
		var find bool
		for _, val := range dplicateList {
			if val.Chapter == c && val.Section == s {
				find = true
				break
			}
		}
		if !find {
			return this.Send(lineNum(), fmt.Errorf("你还没有解锁这个关卡"))
		}
	} else {
		c, s = configs[0].Chapter, configs[0].Section
		config = configs[0]
	}

	// 是否使用临时道具
	if len(tempItemList) > 0 {
		tempItemCoin := tempItemCoin()
		var coin int
		desc := "use items: "
		for index := range tempItemList {
			if tempItemList[index] > 0 {
				coin += tempItemCoin[index]
				desc += fmt.Sprintf("%d , ", index+1)
			}
		}
		if err := this.Role.SubCoin(coin, models.FINANCE_DUPLICATE_USE, desc); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	BattleLogModel := new(models.BattleLogModel)
	BattleLogModel.Uid = this.Uid
	BattleLogModel.Chapter = chapterId
	BattleLogModel.Section = sectionId
	BattleLogModel.Type = models.BattleType(request.GetFightMode())
	if err := models.InsertBattleLog(BattleLogModel); err != nil {
		return this.Send(lineNum(), err)
	}

	response := &protodata.FightInitResponse{
		BattleData: proto.String(config.Value),
		FightMode:  request.FightMode,
		Role:       roleProto(this.Role)}
	return this.Send(StatusOK, response)
}

func (this *Connect) BattleResult() error {

	request := &protodata.FightEndRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	generalCId := int(request.GetGeneralId())
	coin := int(request.GetCoinNum())
	diamond := int(request.GetDiamondNum())
	killNum := int(request.GetKillNum())

	BattleLogModel := models.LastBattleLog(this.Uid)
	if BattleLogModel == nil || BattleLogModel.Result != 0 {
		return this.Send(lineNum(), fmt.Errorf("战斗数据非法:没有战斗初始化"))
	}

	if err := BattleLogModel.SetResult(request.GetIsWin(), killNum); err != nil {
		return this.Send(lineNum(), err)
	}

	var generalData *protodata.GeneralData
	GeneralModel := models.NewGeneralModel(this.Role.Uid)
	if GeneralModel.General(generalCId) != nil {
		//
	} else {
		config := models.ConfigGeneralMap()[generalCId]
		if general := GeneralModel.Insert(config); general != nil {
			return this.Send(lineNum(), fmt.Errorf("数据库错误:新建英雄失败"))
		} else {
			generalData = generalProto(general, config)
		}
	}

	err := this.Role.AddKillNum(killNum, coin, diamond, fmt.Sprintf("Chapter: %d , Section: %d ", BattleLogModel.Chapter, BattleLogModel.Section))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	var find bool
	DuplicateModel := models.NewDuplicateModel(this.Uid)
	duplicateNum := len(DuplicateModel.List())
	for _, val := range DuplicateModel.List() {
		if val.Chapter == BattleLogModel.Chapter && val.Section == BattleLogModel.Section {
			find = true
			break
		}
	}

	if !find {
		if nil == DuplicateModel.Insert(BattleLogModel.Chapter, BattleLogModel.Section) {
			return this.Send(lineNum(), fmt.Errorf("数据库错误:新增数据失败"))
		}
	}

	var length int = 1
	if len(DuplicateModel.List()) != duplicateNum {
		length = 2
	}
	duplicateNum = len(DuplicateModel.List())

	response := &protodata.FightEndResponse{
		Role:    roleProto(this.Role),
		Reward:  rewardProto(coin, diamond, 0, generalData),
		Chapter: duplicateProtoList(DuplicateModel.List()[duplicateNum-length:])}
	return this.Send(StatusOK, response)
}

func rewardProto(coin, diamond, actionValue int, generalData *protodata.GeneralData) *protodata.RewardData {
	return &protodata.RewardData{
		RewardCoin:    proto.Int32(int32(coin)),
		RewardDiamond: proto.Int32(int32(diamond)),
		Stamina:       proto.Int32(int32(actionValue)),
		General:       generalData}
}
