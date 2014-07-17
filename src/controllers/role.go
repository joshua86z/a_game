package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"math/rand"
	"models"
	"protodata"
	"strings"
)

type Role struct {
}

func (this *Role) UserDataRequest(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	if RoleModel == nil {
		RoleModel = &models.RoleModel{}
		RoleModel.Coin = 0
		RoleModel.Diamond = 0
		if err := models.InsertRole(RoleModel); err != nil {
			return lineNum(), nil, err
		}
	}

	SignModel := models.NewSignModel(RoleModel.Uid)
	signDay := SignModel.Times % 7
	if signDay == 0 {
		signDay = 7
	}
	signProto := &protodata.SignRewardData{
		Reward:    nil,
		IsReceive: proto.Bool(true),
		//		IsReceive: proto.Bool(SignModel.Reward),
		SignDay: proto.Int32(int32(signDay)),
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

	response := &protodata.UserDataResponse{
		Role:             roleProto(RoleModel),
		Items:            itemProtoList(models.NewItemModel(RoleModel.Uid).List()),
		Generals:         generalProtoList(models.NewGeneralModel(RoleModel.Uid).List()),
		SignReward:       signProto,
		Chapters:         getDuplicateProto(models.NewDuplicateModel(RoleModel.Uid)),
		TempItemDiamonds: []int32{5, 5, 5, 5},
		CoinProducts:     coinProductProtoList,
		DiamondProducts:  productProtoList}

	return protodata.StatusCode_OK, response, nil
}

// 随机生成角色名字
func (role *Role) RandomName(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	L, err := lua.NewLua("conf/random_name.lua")
	if err != nil {
		return lineNum(), nil, err
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

	return protodata.StatusCode_OK, response, nil
}

func (role *Role) SetRoleName(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	request := &protodata.SetUpNameRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return lineNum(), nil, err
	}

	name := request.GetName()
	if name == "" {
		return lineNum(), nil, fmt.Errorf("名字不能为空")
	}

	rune := []rune(name)
	if len(rune) > 7 {
		rune = rune[:7]
		name = string(rune)
	}

	// 判断是否存在此用户名
	if n := models.NumberByRoleName(name); n > 0 {
		return lineNum(), nil, fmt.Errorf("这个名字已被使用")
	}

	RoleModel.SetName(name)

	return protodata.StatusCode_OK, &protodata.SetUpNameResponse{}, nil
}

func (this *Role) BuyStaminaRequest(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	if RoleModel.ActionValue() >= models.MaxActionValue {
		return lineNum(), nil, fmt.Errorf("体力已满")
	}

	needDiamond := actionValueDiamond()
	if RoleModel.Diamond < needDiamond {
		return lineNum(), nil, fmt.Errorf("钻石不足")
	}

	oldDiamond := RoleModel.Diamond
	oldAction := RoleModel.ActionValue()
	RoleModel.Diamond -= needDiamond
	err := RoleModel.SetActionValue(models.MaxActionValue)
	if err != nil {
		return lineNum(), nil, err
	} else {
		models.InsertSubDiamondFinanceLog(RoleModel.Uid, models.BUY_ACTION, oldDiamond, RoleModel.Diamond, fmt.Sprintf("%d -> %d", oldAction, models.MaxActionValue))
	}

	return protodata.StatusCode_OK, &protodata.BuyStaminaResponse{}, nil
}

func roleProto(RoleModel *models.RoleModel) *protodata.RoleData {

	var roleData protodata.RoleData
	roleData.RoleId = proto.Int64(RoleModel.Uid)
	roleData.RoleName = proto.String("")
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

//// 玩家角色接口struct
//type Role struct {
//
//	// name list for random
//	firstName  []string
//	secondName []string
//	sensitive  []string
//}
//
//// 初始化Role接口
//func NewRole() *Role {
//
//	role := &Role{
//
//		firstName:  []string{},
//		secondName: []string{},
//		sensitive:  []string{},
//
//		//		random: rand.New(rand.NewSource(time.Now().UnixNano())),
//	}
//
//	role.updateNameList()
//
//	return role
//}
//
//// 随机生成角色名字
//func (role *Role) RoleRandomName(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	// random name
//	//	random := role.random
//
//	firstName := role.firstName[random.Intn(len(role.firstName))]
//	length := len(role.secondName)
//	secondName := role.secondName[random.Intn(length)]
//
//	name := firstName + secondName
//
//	response := &pb.RoleRandomNameResponse{
//		ECode: proto.Int32(1),
//		EStr:  proto.String(""),
//		Name:  proto.String(name),
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 设定角色名字 (初始化用户)
//func (role *Role) SetRoleName(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	request := &pb.SetRoleNameRequest{}
//	if err := pd.Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	// 创建新角色
//	name := request.GetName()
//	if !role.checkNameValid(name) {
//
//		// 用户名非法，你懂的
//		return pd.ReturnStr(commandRequest, pb.StatusCode_OK, &pb.SetRoleNameResponse{
//			ECode: proto.Int32(2),
//			EStr:  proto.String(ESTR_character_notvalid),
//		}), nil
//	}
//
//	rune := []rune(name)
//	if len(rune) > 7 {
//		rune = rune[:7]
//		name = string(rune)
//	}
//
//	// 判断是否存在此用户名
//	if r := models.GetRoleByName(name); r.Unique > 0 {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_OK, &pb.SetRoleNameResponse{
//			ECode: proto.Int32(3),
//			EStr:  proto.String(ESTR_role_name_exist),
//		}), nil
//	}
//
//	// 初始化角色 从unique截取uid
//	uid, err := strconv.ParseInt(Substr(fmt.Sprintf("%d", uniqueId), 3, -1), 10, 32)
//	if err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_SERVER_INTERNAL_ERROR, ""), err
//	}
//
//	// 角色初始化数据
//	userRole := models.Role{
//		Unique:     uniqueId,
//		Uid:        int32(uid),
//		Name:       name,
//		Gender:     int(request.GetGender()),
//		HeadId:     int(request.GetHeadId()),
//		Money:      0,
//		Vip:        1,
//		Exp:        0,
//		NewProcess: 0,
//		Status:     1,
//		RegTime:    time.Now().Unix(),
//	}
//	if err = userRole.Insert(); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_CACHE_ERROR, ""), err
//	} else {
//
//		genIds := [4]int{1106, 1416, 1435, 1446}
//		// 初始化角色的武将等其他数据 四个武将国籍不同 初始所带兵种不同
//		for _, genId := range genIds {
//			_, err = models.InsertGeneral(uniqueId, genId, true)
//			if err != nil {
//				return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//			}
//		}
//	}
//
//	L, _ := lua.NewLua("conf/new_role.lua")
//	newCoin := L.GetInt("newCoin")
//	newIngot := L.GetInt("newIngot")
//	newPoint := L.GetInt("newPoint")
//	L.Close()
//
//	if err = models.InsertAccount(uniqueId, newCoin, newIngot, newPoint); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_CACHE_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, &pb.SetRoleNameResponse{
//		ECode: proto.Int32(1),
//	}), nil
//
//}
//
//// 获取玩家角色信息 (如果没有角色，返回0值)
//func (this *Role) RoleBaseInfo(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	userRole := models.GetRole(uniqueId)
//
//	money := models.GetMoney(uniqueId)
//
//	level, exp, maxExp := userRole.GetLevel()
//
//	response := &pb.RoleBaseInfoResponse{
//		ECode:             proto.Int32(1),
//		Name:              proto.String(userRole.Name),
//		Coin:              proto.Int32(int32(money.Coin)),
//		Ingot:             proto.Int32(int32(money.Gold)),
//		Point:             proto.Int32(int32(money.Point)),
//		Stamina:           proto.Int32(int32(models.ActionPoint(uniqueId).Num())),
//		MaxStamina:        proto.Int32(int32(FullStaminaNum)),
//		VipLevel:          proto.Int32(userRole.Vip),
//		StoryStatus:       proto.Int32(int32(models.GetPlotCityId(uniqueId))),
//		GuideStatus:       proto.Int32(userRole.NewProcess),
//		RoleId:            proto.Int64(player.UniqueId),
//		HeadId:            proto.Int32(int32(userRole.HeadId)),
//		RechargeIngot:     proto.Int32(userRole.Money / 10),
//		Exp:               proto.Int32(int32(exp)),
//		MaxExp:            proto.Int32(int32(maxExp)),
//		Level:             proto.Int32(int32(level)),
//		MaxLevel:          proto.Int32(int32(userRole.GetMaxLevel())),
//		BossOpenTime:      proto.Int32(int32(models.NextBossTime())),
//		AutoBattleTimes:   proto.Int32(int32(getTodayAutoBattleTimes(uniqueId))),
//		WipeoutTimes:      proto.Int32(int32(getTodayWipeOutTimes(uniqueId))),
//		WipeoutTotalTimes: proto.Int32(int32(todayWipeOutMaxTimes(uniqueId))),
//		RolePackData:      getRoleData(userRole, money),
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 设置新手进度
//func (role *Role) SetGuideStatus(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//	request := &pb.SetGuideStatusRequest{}
//	if err := pd.Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	newProc := request.GetStatus()
//
//	err := models.UpdateNewProc(player.UniqueId, int(newProc))
//	if err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, &pb.SetGuideStatusResponse{
//		ECode: proto.Int32(1),
//	}), nil
//}
//
//// 玩家角色成就信息
//func (role *Role) AchieveInfo(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.AchieveInfoRequest{}
//	if err := pd.Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	achievementList, err := models.GetRoleAchievements(uniqueId)
//	if err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	var result []struct {
//		Id     int32
//		Name   string
//		Desc   string
//		Ingot  int32
//		Status int32
//		Type   int32
//	}
//
//	configAchievesList := configs.ConfigAchievementGetAll()
//	for _, val := range configAchievesList {
//
//		if val.DisplayType == int(request.GetType()) || request.GetType() == 0 {
//
//			var temp struct {
//				Id     int32
//				Name   string
//				Desc   string
//				Ingot  int32
//				Status int32
//				Type   int32
//			}
//
//			temp.Id = int32(val.Id)
//			temp.Name = string(val.Name)
//			temp.Desc = string(val.Desc)
//			temp.Ingot = int32(val.Ingot)
//			temp.Status = 0
//			temp.Type = int32(val.DisplayType)
//
//			for _, v := range achievementList {
//				if v.Id == val.Id {
//					temp.Status = int32(v.Status)
//				}
//			}
//
//			result = append(result, temp)
//		}
//	}
//
//	//排序 Status大的在前面
//	for i := 0; i < len(result); i++ {
//		for j := len(result) - 1; j > i; j-- {
//			if result[j].Status > result[j-1].Status {
//				tmp := result[j-1]
//				result[j-1] = result[j]
//				result[j] = tmp
//			} else if result[j].Id < result[j-1].Id && result[j].Status == result[j-1].Status {
//				tmp := result[j-1]
//				result[j-1] = result[j]
//				result[j] = tmp
//			}
//		}
//	}
//
//	response := &pb.AchieveInfoResponse{
//		ECode:         proto.Int32(1),
//		AchieveCount:  proto.Int32(int32(len(result))),
//		AchieveId:     []int32{},
//		AchieveName:   []string{},
//		AchieveDesc:   []string{},
//		AchieveStatus: []int32{},
//		Ingot:         []int32{},
//		Type:          []int32{},
//	}
//
//	for _, val := range result {
//		response.AchieveId = append(response.AchieveId, val.Id)
//		response.AchieveName = append(response.AchieveName, val.Name)
//		response.AchieveDesc = append(response.AchieveDesc, val.Desc)
//		response.AchieveStatus = append(response.AchieveStatus, val.Status)
//		response.Ingot = append(response.Ingot, val.Ingot)
//		response.Type = append(response.Type, val.Type)
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 10194 领取成就奖励
//func (role *Role) GetAchieveRewards(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	var err error
//	uniqueId := player.UniqueId
//
//	request := &pb.GetAchieveRewardsRequest{}
//	if err := pd.Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	achId := int(request.GetAchieveId())
//
//	if achievement, err := models.GetAchievementByAchId(uniqueId, achId); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), err
//	} else if achievement.Status != 1 {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("这个成就已领取过")
//	}
//
//	configAchievement := configs.ConfigAchievementGetAll()[achId]
//
//	if err = models.GetAchieveIngot(uniqueId, achId); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if err = models.GetMoney(uniqueId).AddGold(configAchievement.Ingot, models.IG_ACHIVES, "achievement : "+configAchievement.Name); err != nil {
//		return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	response := &pb.GetAchieveRewardsResponse{
//		ECode: proto.Int32(1),
//		Ingot: proto.Int32(int32(configAchievement.Ingot)),
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 获取角色Buff状态信息
//func (role *Role) RoleBuffInfo(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	response := &pb.RoleBuffInfoResponse{
//		ECode:              proto.Int32(1),
//		ExpMultiple:        proto.Int32(1),
//		ExpBuffTime:        proto.Int32(1),
//		PropsMultiple:      proto.Int32(1),
//		PropsBuffTime:      proto.Int32(1),
//		CoinMultiple:       proto.Int32(1),
//		CoinBuffTime:       proto.Int32(1),
//		SoldierExpMultiple: proto.Int32(1),
//		SoldierExpBuffTime: proto.Int32(1),
//		GeneralExpMultiple: proto.Int32(1),
//		GeneralExpBuffTime: proto.Int32(1),
//	}
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 角色购买体力检查
//func (this *Role) RoleBuyStaminaCheck(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	// 判断购买体力状态
//	uniqueId := player.UniqueId
//
//	userRole := models.GetRole(uniqueId)
//
//	times := this.todayBuyActionPointNum(uniqueId)
//
//	// 购买限制次数
//	var limit int
//	for _, vip := range configs.ConfigVipGetAll() {
//		if vip.Level == int(userRole.Vip) {
//			limit = vip.BuyStamina
//			break
//		}
//	}
//
//	needIngot := 0
//	if times < 2 {
//		needIngot = 50
//	} else if times < 6 {
//		needIngot = 100
//	} else {
//		needIngot = 200
//	}
//
//	// 返回
//	response := &pb.RoleBuyStaminaCheckResponse{
//		ECode:           proto.Int32(1),
//		BuyStatus:       proto.Bool(times < limit),
//		BuyStaminaTimes: proto.Int32(int32(times)),
//		BuyCost:         proto.Int32(int32(needIngot)),
//	}
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 角色购买体力
//func (this *Role) RoleBuyStamina(player *Player, commandRequest *pb.CommandRequest) (string, error) {
//
//	// 判断购买体力状态
//	uniqueId := player.UniqueId
//
//	userRole := models.GetRole(uniqueId)
//
//	times := this.todayBuyActionPointNum(uniqueId)
//
//	// 增加体力逻辑
//	var response *pb.RoleBuyStaminaResponse
//
//	// 购买限制次数
//	var limit int = 2
//	for _, vip := range configs.ConfigVipGetAll() {
//		if vip.Level == int(userRole.Vip) {
//			limit = vip.BuyStamina
//			break
//		}
//	}
//
//	// 可购买
//	if times < limit {
//
//		actionPoint := models.ActionPoint(uniqueId)
//		if actionPoint.Num() < FullStaminaNum {
//
//			// 购买所需元宝
//			needIngot := 0
//			if times < 2 {
//				needIngot = 50
//			} else if times < 6 {
//				needIngot = 100
//			} else {
//				needIngot = 200
//			}
//
//			money := models.GetMoney(uniqueId)
//
//			if money.Gold < needIngot {
//
//				// 钱不够
//				response = &pb.RoleBuyStaminaResponse{
//					ECode:      proto.Int32(3),
//					EStr:       proto.String(ESTR_not_enough_ingot),
//					Stamina:    proto.Int32(int32(actionPoint.Num())),
//					MaxStamina: proto.Int32(int32(FullStaminaNum)),
//				}
//
//			} else {
//
//				if err := money.SubGold(needIngot, models.IG_BUY_STA, ""); err != nil {
//					return pd.ReturnStr(commandRequest, pb.StatusCode_DATABASE_ERROR, ""), err
//				}
//
//				// 增加体力
//				actionPoint.Add(20 - actionPoint.Num())
//
//				response = &pb.RoleBuyStaminaResponse{
//					ECode:      proto.Int32(1),
//					Stamina:    proto.Int32(int32(actionPoint.Num())),
//					MaxStamina: proto.Int32(int32(FullStaminaNum)),
//				}
//			}
//		} else {
//
//			// 满体状态
//			response = &pb.RoleBuyStaminaResponse{
//				ECode:      proto.Int32(1),
//				Stamina:    proto.Int32(int32(FullStaminaNum)),
//				MaxStamina: proto.Int32(int32(FullStaminaNum)),
//			}
//		}
//	} else {
//
//		// 超过购买次数限制，无法购买
//		response = &pb.RoleBuyStaminaResponse{
//			ECode:      proto.Int32(2),
//			EStr:       proto.String(ESTR_cant_buy_stamina),
//			Stamina:    proto.Int32(int32(models.ActionPoint(uniqueId).Num())),
//			MaxStamina: proto.Int32(int32(FullStaminaNum)),
//		}
//	}
//
//	this.addBuyActionPointNum(uniqueId)
//
//	return pd.ReturnStr(commandRequest, pb.StatusCode_OK, response), nil
//}
//
//// 获取姓名列表供随机组合
//func (role *Role) updateNameList() {
//
//	//role.firstName, err = role.model.GetFirstNameList()
//	role.firstName = configs.GetFirstNameList()
//
//	//	length := 10
//	if len(role.firstName) < 5 {
//		//		length = len(role.firstName)
//	}
//
//	//role.secondName, err = role.model.GetSecondNameList()
//	role.secondName = configs.GetSecondNameList()
//
//	if len(role.secondName) < 5 {
//		//		length = len(role.secondName)
//	}
//
//	// Sensitive Word
//	sens := configs.ConfigSensitiveWordGetAll()
//
//	for _, sen := range sens {
//		role.sensitive = append(role.sensitive, sen.Word)
//	}
//
//}
//
//// 过滤非法字符
//func (role *Role) checkNameValid(name string) bool {
//
//	if len(name) < 1 {
//		return false
//	}
//
//	for _, sens := range role.sensitive {
//		if pos := strings.Index(name, sens); pos != -1 {
//			return false
//		}
//	}
//
//	return true
//}
//
//// 检查是否可购买体力 返回当日已购买次数
//func (this *Role) todayBuyActionPointNum(uniqueId int64) int {
//
//	key := fmt.Sprintf("%d_BUYACTIONPOINTNUM_%d", ServerId, uniqueId)
//	str, _ := ssdb.SSDB().Get(key)
//	if str != "" {
//		date := time.Now().Format("0102")
//		array := strings.Split(str, ",")
//		if date == array[1] {
//			ret, _ := strconv.Atoi(array[0])
//			return ret
//		}
//	}
//	return 0
//}
//
//func (this *Role) addBuyActionPointNum(uniqueId int64) {
//	num := 0
//	date := time.Now().Format("0102")
//	key := fmt.Sprintf("%d_BUYACTIONPOINTNUM_%d", ServerId, uniqueId)
//	str, _ := ssdb.SSDB().Get(key)
//	if str != "" {
//		array := strings.Split(str, ",")
//		if date == array[1] {
//			num, _ = strconv.Atoi(array[0])
//		}
//	}
//	num += 1
//	ssdb.SSDB().Set(key, fmt.Sprintf("%d,%s", num, date))
//}
//
//func getRoleData(role models.Role, money models.RoleMoney) *pb.RolePackData {
//
//	level, exp, displayExp := role.GetLevel()
//
//	return &pb.RolePackData{
//		RoleId:     proto.Int64(int64(role.Unique)),
//		Exp:        proto.Int32(int32(exp)),
//		MaxExp:     proto.Int32(int32(displayExp)),
//		Level:      proto.Int32(int32(level)),
//		MaxLevel:   proto.Int32(int32(role.GetMaxLevel())),
//		Coin:       proto.Int32(int32(money.Coin)),
//		Ingot:      proto.Int32(int32(money.Gold)),
//		Point:      proto.Int32(int32(money.Point)),
//		Stamina:    proto.Int32(int32(models.ActionPoint(role.Unique).Num())),
//		MaxStamina: proto.Int32(int32(20)),
//		VipLevel:   proto.Int32(int32(role.Vip)),
//	}
//}
