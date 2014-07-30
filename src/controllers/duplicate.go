package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"models"
	"protodata"
)

func duplicateProtoList(duplicates []*models.DuplicateData) []*protodata.ChapterData {

	list := models.ConfigDuplicateList()

	var result []*protodata.ChapterData
	result = append(result, &protodata.ChapterData{
		ChapterId:   proto.Int32(int32(list[0].Chapter)),
		ChapterName: proto.String(list[0].ChapterName),
		ChapterDesc: proto.String(list[0].ChapterDesc),
		IsUnlock:    proto.Bool(true),
	})

	for index, section := range list {

		var sectionProto protodata.SectionData
		sectionProto.SectionId = proto.Int32(int32(section.Section))
		sectionProto.SectionName = proto.String(section.SectionName)
		sectionProto.SectionDesc = proto.String(section.SectionDesc)
		sectionProto.IsUnlock = proto.Bool(true)

		var find bool
		if index > 0 {
			for _, d := range duplicates {
				if d.Chapter == list[index-1].Chapter && d.Section == list[index-1].Section {
					find = true
					break
				} else {
					find = false
				}
			}
			if !find {
				sectionProto.IsUnlock = proto.Bool(false)
			}
		}

		if section.Chapter != int(*result[len(result)-1].ChapterId) {

			result = append(result, &protodata.ChapterData{
				ChapterId:   proto.Int32(int32(section.Chapter)),
				ChapterName: proto.String(section.ChapterName),
				ChapterDesc: proto.String(section.ChapterDesc),
				IsUnlock:    proto.Bool(find),
			})
		}
		result[len(result)-1].Sections = append(result[len(result)-1].Sections, &sectionProto)
	}

	return result
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
//	duplicate := &Role{
//
//		firstName:  []string{},
//		secondName: []string{},
//		sensitive:  []string{},
//
//		//		random: rand.New(rand.NewSource(time.Now().UnixNano())),
//	}
//
//	duplicate.updateNameList()
//
//	return duplicate
//}
//
//// 随机生成角色名字
//func (duplicate *Role) RoleRandomName(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	// random name
//	//	random := duplicate.random
//
//	firstName := duplicate.firstName[random.Intn(len(duplicate.firstName))]
//	length := len(duplicate.secondName)
//	secondName := duplicate.secondName[random.Intn(length)]
//
//	name := firstName + secondName
//
//	response := &pb.RoleRandomNameResponse{
//		ECode: proto.Int32(1),
//		EStr:  proto.String(""),
//		Name:  proto.String(name),
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 设定角色名字 (初始化用户)
//func (duplicate *Role) SetRoleName(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	request := &pb.SetRoleNameRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	// 创建新角色
//	name := request.GetName()
//	if !duplicate.checkNameValid(name) {
//
//		// 用户名非法，你懂的
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SetRoleNameResponse{
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
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SetRoleNameResponse{
//			ECode: proto.Int32(3),
//			EStr:  proto.String(ESTR_duplicate_name_exist),
//		}), nil
//	}
//
//	// 初始化角色 从unique截取uid
//	uid, err := strconv.ParseInt(Substr(fmt.Sprintf("%d", uniqueId), 3, -1), 10, 32)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_SERVER_INTERNAL_ERROR, ""), err
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
//		return pd.ReturnStr(cq, pb.StatusCode_CACHE_ERROR, ""), err
//	} else {
//
//		genIds := [4]int{1106, 1416, 1435, 1446}
//		// 初始化角色的武将等其他数据 四个武将国籍不同 初始所带兵种不同
//		for _, genId := range genIds {
//			_, err = models.InsertGeneral(uniqueId, genId, true)
//			if err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//			}
//		}
//	}
//
//	L, _ := lua.NewLua("conf/new_duplicate.lua")
//	newCoin := L.GetInt("newCoin")
//	newIngot := L.GetInt("newIngot")
//	newPoint := L.GetInt("newPoint")
//	L.Close()
//
//	if err = models.InsertAccount(uniqueId, newCoin, newIngot, newPoint); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_CACHE_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SetRoleNameResponse{
//		ECode: proto.Int32(1),
//	}), nil
//
//}
//
//// 获取玩家角色信息 (如果没有角色，返回0值)
//func (this *Role) RoleBaseInfo(player *Player, cq *pb.CommandRequest) (string, error) {
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
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 设置新手进度
//func (duplicate *Role) SetGuideStatus(player *Player, cq *pb.CommandRequest) (string, error) {
//	request := &pb.SetGuideStatusRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	newProc := request.GetStatus()
//
//	err := models.UpdateNewProc(player.UniqueId, int(newProc))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SetGuideStatusResponse{
//		ECode: proto.Int32(1),
//	}), nil
//}
//
//// 玩家角色成就信息
//func (duplicate *Role) AchieveInfo(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.AchieveInfoRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	achievementList, err := models.GetRoleAchievements(uniqueId)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
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
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 10194 领取成就奖励
//func (duplicate *Role) GetAchieveRewards(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	var err error
//	uniqueId := player.UniqueId
//
//	request := &pb.GetAchieveRewardsRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	achId := int(request.GetAchieveId())
//
//	if achievement, err := models.GetAchievementByAchId(uniqueId, achId); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	} else if achievement.Status != 1 {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("这个成就已领取过")
//	}
//
//	configAchievement := configs.ConfigAchievementGetAll()[achId]
//
//	if err = models.GetAchieveIngot(uniqueId, achId); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if err = models.GetMoney(uniqueId).AddGold(configAchievement.Ingot, models.IG_ACHIVES, "achievement : "+configAchievement.Name); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	response := &pb.GetAchieveRewardsResponse{
//		ECode: proto.Int32(1),
//		Ingot: proto.Int32(int32(configAchievement.Ingot)),
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 获取角色Buff状态信息
//func (duplicate *Role) RoleBuffInfo(player *Player, cq *pb.CommandRequest) (string, error) {
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
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 角色购买体力检查
//func (this *Role) RoleBuyStaminaCheck(player *Player, cq *pb.CommandRequest) (string, error) {
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
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 角色购买体力
//func (this *Role) RoleBuyStamina(player *Player, cq *pb.CommandRequest) (string, error) {
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
//					return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
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
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 获取姓名列表供随机组合
//func (duplicate *Role) updateNameList() {
//
//	//duplicate.firstName, err = duplicate.model.GetFirstNameList()
//	duplicate.firstName = configs.GetFirstNameList()
//
//	//	length := 10
//	if len(duplicate.firstName) < 5 {
//		//		length = len(duplicate.firstName)
//	}
//
//	//duplicate.secondName, err = duplicate.model.GetSecondNameList()
//	duplicate.secondName = configs.GetSecondNameList()
//
//	if len(duplicate.secondName) < 5 {
//		//		length = len(duplicate.secondName)
//	}
//
//	// Sensitive Word
//	sens := configs.ConfigSensitiveWordGetAll()
//
//	for _, sen := range sens {
//		duplicate.sensitive = append(duplicate.sensitive, sen.Word)
//	}
//
//}
//
//// 过滤非法字符
//func (duplicate *Role) checkNameValid(name string) bool {
//
//	if len(name) < 1 {
//		return false
//	}
//
//	for _, sens := range duplicate.sensitive {
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
//func getRoleData(duplicate models.Role, money models.RoleMoney) *pb.RolePackData {
//
//	level, exp, displayExp := duplicate.GetLevel()
//
//	return &pb.RolePackData{
//		RoleId:     proto.Int64(int64(duplicate.Unique)),
//		Exp:        proto.Int32(int32(exp)),
//		MaxExp:     proto.Int32(int32(displayExp)),
//		Level:      proto.Int32(int32(level)),
//		MaxLevel:   proto.Int32(int32(duplicate.GetMaxLevel())),
//		Coin:       proto.Int32(int32(money.Coin)),
//		Ingot:      proto.Int32(int32(money.Gold)),
//		Point:      proto.Int32(int32(money.Point)),
//		Stamina:    proto.Int32(int32(models.ActionPoint(duplicate.Unique).Num())),
//		MaxStamina: proto.Int32(int32(20)),
//		VipLevel:   proto.Int32(int32(duplicate.Vip)),
//	}
//}
