package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"math/rand"
	"models"
	"protodata"
	"strings"
	"time"
)

func (this *Connect) UserDataRequest() error {

	var isReceive bool
	if this.Role.SignDate == time.Now().Format("20060102") {
		isReceive = true
	}
	if !isReceive {
		if err := this.Role.Sign(); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	configs := models.BaseGeneralMap()

	//var rewardList []*protodata.RewardData
	//for i := this.Role.SignNum; i < this.Role.SignNum+7; i++ {

	//	c, d, s, g := signReward(i)

	//	temp := new(protodata.RewardData)
	//	temp.RewardCoin = proto.Int32(int32(c))
	//	temp.RewardDiamond = proto.Int32(int32(d))
	//	temp.Stamina = proto.Int32(int32(s))
	//	if g > 0 {
	//		temp.General = generalProto(new(models.GeneralData), configs[g])
	//	}

	//	rewardList = append(rewardList, temp)
	//}

	signDay := this.Role.SignNum % 7
	if signDay == 0 {
		signDay = 7
	}

	signProto := &protodata.SignRewardData{
		Reward:    signProto(this.Role.SignNum),
		IsReceive: proto.Bool(isReceive),
		SignDay:   proto.Int32(int32(signDay)),
	}

	var coinProductProtoList []*protodata.CoinProductData
	coinDiamond := models.ConfigCoinDiamondList()
	for _, val := range coinDiamond {
		coinProductProtoList = append(coinProductProtoList, &protodata.CoinProductData{
			ProductIndex: proto.Int32(int32(val.Index)),
			ProductName:  proto.String(val.Name),
			ProductDesc:  proto.String(val.Desc),
			ProductCoin:  proto.Int32(int32(val.Coin)),
			PriceDiamond: proto.Int32(int32(val.Diamond)),
		})
	}

	var productProtoList []*protodata.DiamondProductData
	productList := models.ConfigPayCenterList()
	for _, val := range productList {
		productProtoList = append(productProtoList, &protodata.DiamondProductData{
			ProductIndex:   proto.Int32(int32(val.Id)),
			ProductName:    proto.String(val.Name),
			ProductDesc:    proto.String(val.Desc),
			ProductDiamond: proto.Int32(int32(val.Diamond)),
			Price:          proto.Int32(int32(val.Money)),
		})
	}

	generalList := models.General.List(this.Uid)
	//if len(generalList) == 0 {
	//	Lua, _ := lua.NewLua("conf/general.lua")
	//	s := Lua.GetString("init_generals")
	//	Lua.Close()
	//	array := strings.Split(s, ",")
	//	generalList = append(generalList, models.General.Insert(this.Uid, configs[models.Atoi(array[0])]))
	//	generalList = append(generalList, models.General.Insert(this.Uid, configs[models.Atoi(array[1])]))
	//	generalList = append(generalList, models.General.Insert(this.Uid, configs[models.Atoi(array[2])]))
	//}

	if !isReceive {
		coin, diamond, action, generalId := signReward(this.Role.SignNum)
		if coin > 0 {
			this.Role.AddCoin(coin, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
		} else if diamond > 0 {
			this.Role.AddDiamond(diamond, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
		} else if action > 0 {
			this.Role.SetActionValue(this.Role.ActionValue() + action)
		} else if generalId > 0 {
			var find bool
			for _, val := range generalList {
				if generalId == val.BaseId {
					find = true
					break
				}
			}
			baseGeneral := configs[generalId]
			if find {
				this.Role.AddDiamond(baseGeneral.BuyDiamond, models.FINANCE_SIGN_GET, fmt.Sprintf("signDay : %d", signDay))
			} else {
				newGeneral, err := models.General.Insert(this.Uid, baseGeneral)
				if err != nil {
					return this.Send(lineNum(), err)
				}
				generalList = append(generalList, newGeneral)
			}
		}
	}

	tempItemDiamond := tempItemDiamond()

	return this.Send(StatusOK, &protodata.UserDataResponse{
		Role:             roleProto(this.Role),
		Items:            itemProtoList(models.Item.List(this.Uid)),
		Generals:         generalProtoList(generalList, configs),
		SignReward:       signProto,
		Chapters:         duplicateProtoList(models.NewDuplicateModel(this.Uid).List(), models.ConfigDuplicateList()),
		TempItemDiamonds: []int32{int32(tempItemDiamond[0]), int32(tempItemDiamond[1]), int32(tempItemDiamond[2]), int32(tempItemDiamond[3])},
		CoinProducts:     coinProductProtoList,
		DiamondProducts:  productProtoList,
		LeaderId:         proto.Int32(int32(this.Role.GeneralBaseId)),
		StaminaTimeMax:   proto.Int32(int32(models.Role.ActionWaitTime))})
}

// 随机生成角色名字
func (this *Connect) RandomName() error {

	L, err := lua.NewLua("conf/random_name.lua")
	if err != nil {
		return this.Send(lineNum(), err)
	}

	firstNameStr := L.GetString("first_name")
	SecondNameStr := L.GetString("second_name")

	L.Close()

	firstNameArray := strings.Split(firstNameStr, ",")
	SecondNameArray := strings.Split(SecondNameStr, ",")

	firstName := firstNameArray[rand.Intn(len(firstNameArray))]
	secondName := SecondNameArray[rand.Intn(len(SecondNameArray))]

	response := &protodata.RandomNameResponse{
		Name: proto.String(firstName + secondName),
	}

	return this.Send(StatusOK, response)
}

func (this *Connect) SetRoleName() error {

	request := &protodata.SetUpNameRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	name := request.GetName()

	L, _ := lua.NewLua("conf/random_name.lua")
	sensitiveWord := L.GetString("sensitive_word")
	L.Close()

	wordArray := strings.Split(sensitiveWord, ",")

	// 过滤敏感字
	for _, word := range wordArray {
		name = strings.Replace(name, word, "", len(name))
	}

	if name == "" {
		return this.Send(lineNum(), fmt.Errorf("名字不能为空"))
	}

	rune := []rune(name)
	if len(rune) > 7 {
		rune = rune[:7]
		name = string(rune)
	}

	// 判断是否存在此用户名

	if n, _ := models.Role.NumberByRoleName(name); n > 0 {
		return this.Send(lineNum(), fmt.Errorf("这个名字已被使用"))
	}

	if err := this.Role.SetName(name); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.SetUpNameResponse{Role: roleProto(this.Role)})
}

func (this *Connect) BuyStaminaRequest() error {

	if this.Role.ActionValue() >= models.Role.MaxActionValue {
		return this.Send(lineNum(), fmt.Errorf("体力已满"))
	}

	needDiamond := buyActionDiamond(this.Role.BuyActionNum)
	if this.Role.Diamond < needDiamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	this.Role.BuyActionNum += 1
	err := this.Role.BuyActionValue(needDiamond, models.Role.MaxActionValue)
	if err != nil {
		this.Role.BuyActionNum -= 1
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.BuyStaminaResponse{
		Role:    roleProto(this.Role),
		Stamina: proto.Int32((5))})
}

func roleProto(RoleModel *models.RoleData) *protodata.RoleData {

	return &protodata.RoleData{
		RoleId:            proto.Int64(RoleModel.Uid),
		RoleName:          proto.String(RoleModel.Name),
		Stamina:           proto.Int32(int32(RoleModel.ActionValue())),
		MaxStamina:        proto.Int32(int32(models.Role.MaxActionValue)),
		Coin:              proto.Int32(int32(RoleModel.Coin)),
		Diamond:           proto.Int32(int32(RoleModel.Diamond)),
		SuppleStaminaTime: proto.Int32(int32(RoleModel.ActionRecoverTime())),
		SuppleStaDiamond:  proto.Int32(int32(buyActionDiamond(RoleModel.BuyActionNum))),
		KillNum:           proto.Int32(int32(RoleModel.KillNum))}
}

func buyActionDiamond(n int) int {
	Lua, _ := lua.NewLua("conf/role.lua")
	Lua.L.GetGlobal("buyActionDiamond")
	Lua.L.DoString(fmt.Sprintf("diamond = buyActionDiamond(%d)", n))
	diamond := Lua.GetInt("diamond")
	Lua.Close()
	return diamond
}

func signReward(times int) (coin, diamond, action, generalId int) {
	Lua, _ := lua.NewLua("conf/sign_reward.lua")
	Lua.L.GetGlobal("signReward")
	Lua.L.DoString(fmt.Sprintf("coin, diamond, action, generalId = signReward(%d)", times))
	coin = Lua.GetInt("coin")
	diamond = Lua.GetInt("diamond")
	action = Lua.GetInt("action")
	generalId = Lua.GetInt("generalId")
	Lua.Close()
	return
}
