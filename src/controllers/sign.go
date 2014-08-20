package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
	"time"
)

func (this *Connect) Sign() error {

	response := new(protodata.SignRewardResponse)

	signDay := this.Role.SignNum % 7
	if signDay == 0 {
		signDay = 7
	}

	response.SignReward = &protodata.SignRewardData{SignDay: proto.Int32(int32(signDay))}

	if this.Role.SignDate == time.Now().Format("20060102") {
		response.Role = roleProto(this.Role)
		response.SignReward.IsReceive = proto.Bool(true)
		return this.Send(StatusOK, response)
	}

	if err := this.Role.Sign(); err != nil {
		return this.Send(lineNum(), err)
	}

	configs := models.BaseGeneralMap()
	var rewardList []*protodata.RewardData
	for i := this.Role.SignNum; i < this.Role.SignNum+7; i++ {

		c, d, s, g := signReward(i)

		temp := new(protodata.RewardData)
		temp.RewardCoin = proto.Int32(int32(c))
		temp.RewardDiamond = proto.Int32(int32(d))
		temp.Stamina = proto.Int32(int32(s))
		if g > 0 {
			temp.General = generalProto(new(models.GeneralData), configs[g])
		}

		rewardList = append(rewardList, temp)
	}

	coin, diamond, action, generalId := signReward(this.Role.SignNum)
	if coin > 0 {
		this.Role.AddCoin(coin, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
	} else if diamond > 0 {
		this.Role.AddDiamond(diamond, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
	} else if action > 0 {
		this.Role.SetActionValue(this.Role.ActionValue() + action)
	} else if generalId > 0 {
		var find bool
		for _, val := range models.General.List(this.Uid) {
			if generalId == val.BaseId {
				find = true
				break
			}
		}
		config := configs[generalId]
		if find {
			this.Role.AddDiamond(config.BuyDiamond, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
		} else {
			response.General = generalProto(models.General.Insert(this.Uid, config), config)
		}
	}

	response.Role = roleProto(this.Role)
	return this.Send(StatusOK, response)
}
