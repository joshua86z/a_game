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

	signDay := this.Role.SignTimes % 7
	if signDay == 0 {
		signDay = 7
	}

	var isReceive bool
	if this.Role.SignDate == time.Now().Format("20060102") {
		isReceive = true
	}
	if !isReceive {
		if err := this.Role.Sign(); err != nil {
			return this.Send(lineNum(), err)
		}
		//		SignModel.GetReward()
	}

	var rewardList []*protodata.RewardData
	for i := 1; i <= 7; i++ {
		var temp protodata.RewardData
		temp.RewardCoin = proto.Int32(5)
		temp.RewardDiamond = proto.Int32(5)
		rewardList = append(rewardList, &temp)
	}

	signProto := &protodata.SignRewardData{
		Reward:    rewardList,
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
			Price:          proto.Int32(int32(val.Rmb)),
		})
	}

	GeneralModel := models.NewGeneralModel(this.Uid)
	if len(GeneralModel.List()) == 0 {
		Lua, _ := lua.NewLua("conf/new_role.lua")
		s := Lua.GetString("init_generals")
		Lua.Close()
		fmt.Println(s)
		array := strings.Split(s, ",")
		configs := models.ConfigGeneralMap()
		GeneralModel.Insert(configs[models.Atoi(array[0])])
		GeneralModel.Insert(configs[models.Atoi(array[1])])
		GeneralModel.Insert(configs[models.Atoi(array[2])])
	}

	response := &protodata.UserDataResponse{
		Role:             roleProto(this.Role),
		Items:            itemProtoList(models.NewItemModel(this.Uid).List()),
		Generals:         generalProtoList(GeneralModel.List()),
		SignReward:       signProto,
		Chapters:         duplicateProtoList(models.NewDuplicateModel(this.Uid).List()),
		TempItemDiamonds: []int32{5, 5, 5, 5},
		CoinProducts:     coinProductProtoList,
		DiamondProducts:  productProtoList}

	return this.Send(StatusOK, response)
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
	if name == "" {
		return this.Send(lineNum(), fmt.Errorf("名字不能为空"))
	}

	rune := []rune(name)
	if len(rune) > 7 {
		rune = rune[:7]
		name = string(rune)
	}

	// 判断是否存在此用户名
	if n := models.NumberByRoleName(name); n > 0 {
		return this.Send(lineNum(), fmt.Errorf("这个名字已被使用"))
	}

	if err := this.Role.SetName(name); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.SetUpNameResponse{Role: roleProto(this.Role)})
}

func (this *Connect) BuyStaminaRequest() error {

	if this.Role.ActionValue() >= models.MaxActionValue {
		return this.Send(lineNum(), fmt.Errorf("体力已满"))
	}

	needDiamond := actionValueDiamond()
	if this.Role.Diamond < needDiamond {
		return this.Send(lineNum(), fmt.Errorf("钻石不足"))
	}

	err := this.Role.BuyActionValue(needDiamond, models.MaxActionValue)
	if err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.BuyStaminaResponse{})
}

func roleProto(RoleModel *models.RoleModel) *protodata.RoleData {

	var roleData protodata.RoleData
	roleData.RoleId = proto.Int64(RoleModel.Uid)
	roleData.RoleName = proto.String(RoleModel.Name)
	roleData.Stamina = proto.Int32(int32(RoleModel.ActionValue()))
	roleData.MaxStamina = proto.Int32(int32(models.MaxActionValue))
	roleData.Coin = proto.Int32(int32(RoleModel.Coin))
	roleData.Diamond = proto.Int32(int32(RoleModel.Diamond))
	roleData.SuppleStaminaTime = proto.Int32(int32(RoleModel.ActionRecoverTime()))
	roleData.SuppleStaDiamond = proto.Int32(int32(actionValueDiamond()))
	roleData.KillNum = proto.Int32(int32(0))
	return &roleData
}

func actionValueDiamond() int {
	return 5
}
