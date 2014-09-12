package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"models"
	"protodata"
	"time"
)

func (this *Connect) Sign() error {

	response := new(protodata.SignRewardResponse)

	if this.Role.SignDate == time.Now().Format("20060102") {

		signDay := this.Role.SignNum % 7
		if signDay == 0 {
			signDay = 7
		}

		response.SignReward = &protodata.SignRewardData{
			SignDay: proto.Int32(int32(signDay)),
			Reward:  signProto(this.Role.SignNum)}
		response.Role = roleProto(this.Role)
		response.SignReward.IsReceive = proto.Bool(true)
		return this.Send(StatusOK, response)
	}

	if err := this.Role.Sign(); err != nil {
		return this.Send(lineNum(), err)
	}

	signDay := this.Role.SignNum % 7
	if signDay == 0 {
		signDay = 7
	}

	response.SignReward = &protodata.SignRewardData{
		SignDay: proto.Int32(int32(signDay)),
		Reward:  signProto(this.Role.SignNum)}

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
		baseGeneral := models.BaseGeneral(generalId, nil)
		if find {
			this.Role.AddDiamond(baseGeneral.BuyDiamond, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
		} else {
			general, err := models.General.Insert(this.Uid, baseGeneral)
			if err != nil {
				return this.Send(lineNum(), response)
			}
			response.General = generalProto(general, baseGeneral)
		}
	}

	response.Role = roleProto(this.Role)
	return this.Send(StatusOK, response)
}

func signProto(day int) []*protodata.RewardData {

	var result []*protodata.RewardData
	Lua, _ := lua.NewLua("conf/sign_reward.lua")

	begin := day/7 - 1
	if begin < 0 {
		begin = 0
	}
	begin = begin*7 + 1
	for i := begin; i <= begin+7; i++ {

		//Lua.L.GetGlobal("signReward")
		Lua.L.DoString(fmt.Sprintf("coin, diamond, action, generalId = signReward(%d)", i))
		coin := Lua.GetInt("coin")
		diamond := Lua.GetInt("diamond")
		action := Lua.GetInt("action")
		generalId := Lua.GetInt("generalId")

		temp := new(protodata.RewardData)
		temp.RewardCoin = proto.Int32(int32(coin))
		temp.RewardDiamond = proto.Int32(int32(diamond))
		temp.Stamina = proto.Int32(int32(action))
		if generalId > 0 {
			config := models.BaseGeneral(generalId, nil)
			temp.General = generalProto(new(models.GeneralData), config)
		}

		result = append(result, temp)
	}

	Lua.Close()
	return result
}
