package controllers

import (
//	"code.google.com/p/goprotobuf/proto"
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

	var c, s int

	if chapterId != 1 || sectionId != 1 {
		configs := models.ConfigDuplicateList()
		for index, val := range configs {
			if val.Chapter > chapterId {
				break
			}
			if val.Chapter == chapterId && val.Section == sectionId {
				if index > 0 {
					index -= 1

				}
				c, s = configs[index].Chapter, configs[index].Section
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

	}

	var BattleLogModel models.BattleLogModel
	BattleLogModel.Uid = this.Uid
	BattleLogModel.Chapter = chapterId
	BattleLogModel.Section = sectionId
	BattleLogModel.Type = models.BattleType(request.GetFightMode())
	if err := models.InsertBattleLog(&BattleLogModel); err != nil {
		return this.Send(lineNum(), err)
	}

	return nil
}

func (this *Connect) BattleResult() error {

	request := &protodata.FightInitResponse{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	return nil
}
