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

	actionValue := this.Role.ActionValue()
	if actionValue < 1 {
		return this.Send(lineNum(), fmt.Errorf("体力不足"))
	} else {
		if err := this.Role.SetActionValue(actionValue - 1); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	// 是否使用临时道具
	if len(tempItemList) > 0 {
		tempItemDiamond := tempItemDiamond()
		var diamond int
		desc := "use items: "
		for _, val := range tempItemList {
			diamond += tempItemDiamond[val]
			desc += fmt.Sprintf("%d", val+1)
		}
		if err := this.Role.SubDiamond(diamond, models.FINANCE_DUPLICATE_USE, desc); err != nil {
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
		TempItems:  request.TempItems,
		Role:       roleProto(this.Role)}
	return this.Send(StatusOK, response)
}

func (this *Connect) BattleResult() error {

	request := &protodata.FightEndRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

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

	response := &protodata.FightEndResponse{Role: roleProto(this.Role)}

	if killNum > 0 {
		general := models.General.General(this.Uid, this.Role.GeneralConfigId)
		general.AddKillNum(killNum)
		response.General = generalProto(general, models.BaseGeneral(general.BaseId, nil))
	}

	var generalData *protodata.GeneralData
	if BattleLogModel.Type == 0 && request.GetIsWin() {

		var find bool
		DuplicateModel := models.NewDuplicateModel(this.Uid)
		//duplicateNum := len(DuplicateModel.List())
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

		//var length int = 1
		//if len(DuplicateModel.List()) != duplicateNum {
		//	length = 2
		//}
		//duplicateNum = len(DuplicateModel.List())
		//begin := duplicateNum - length
		//if begin < 0 {
		//	begin = 0
		//}

		list := models.ConfigDuplicateList()

		last := DuplicateModel.List()[len(DuplicateModel.List())-1]
		var base *models.ConfigDuplicate
		find = false
		for _, duplicate := range list {
			if find == true {
				base = duplicate
				break
			}
			if duplicate.Chapter == last.Chapter && duplicate.Section == last.Section {
				base = duplicate
				find = true
			}
		}

		section := new(protodata.SectionData)
		section.SectionId = proto.Int32(int32(base.Section))
		section.SectionName = proto.String(base.SectionName)
		section.IsUnlock = proto.Bool(true)
		chapter := new(protodata.ChapterData)
		chapter.ChapterId = proto.Int32(int32(base.Chapter))
		chapter.ChapterName = proto.String(base.ChapterName)
		chapter.IsUnlock = proto.Bool(true)
		chapter.Sections = append(chapter.Sections, section)
		response.Chapter = append(response.Chapter, chapter)

		// 奖励英雄
		baseId := 0
		for _, val := range list {
			if val.Chapter == BattleLogModel.Chapter && val.Section == BattleLogModel.Section {
				if val.GenId > 0 {
					baseId = val.GenId
				}
				break
			}
		}

		if baseId > 0 {
			if general := models.General.General(this.Uid, baseId); general == nil {
				baseGeneral := models.BaseGeneral(baseId, nil)
				if general := models.General.Insert(this.Uid, baseGeneral); general == nil {
					return this.Send(lineNum(), fmt.Errorf("数据库错误:新建英雄失败"))
				} else {
					generalData = generalProto(general, baseGeneral)
				}
			}
		}
	} else {
		if request.GetIsWin() {
			this.Role.UnlimitedNum += 1
			if this.Role.UnlimitedNum > this.Role.UnlimitedMaxNum {
				this.Role.UnlimitedMaxNum = this.Role.UnlimitedNum
			}
		} else if this.Role.UnlimitedNum != 0 {
			this.Role.UnlimitedNum = 0
		}
	}

	response.Reward = rewardProto(coin, diamond, 0, generalData)

	oldCoin, oldDiamond := this.Role.Coin, this.Role.Diamond

	if killNum > 0 {
		this.Role.KillNum += killNum
	}
	if coin > 0 {
		this.Role.Coin += coin
	}
	if diamond > 0 {
		this.Role.Diamond += diamond
	}

	if err := this.Role.Set(); err != nil {
		this.Role = nil
		return this.Send(lineNum(), err)
	}

	if coin > 0 {
		models.InsertAddCoinFinanceLog(this.Uid, models.FINANCE_DUPLICATE_GET, oldCoin, this.Role.Coin, "")
	}
	if diamond > 0 {
		models.InsertAddDiamondFinanceLog(this.Uid, models.FINANCE_DUPLICATE_GET, oldDiamond, this.Role.Diamond, "")
	}

	return this.Send(StatusOK, response)
}

func rewardProto(coin, diamond, actionValue int, generalData *protodata.GeneralData) *protodata.RewardData {
	return &protodata.RewardData{
		RewardCoin:    proto.Int32(int32(coin)),
		RewardDiamond: proto.Int32(int32(diamond)),
		Stamina:       proto.Int32(int32(actionValue)),
		General:       generalData}
}
