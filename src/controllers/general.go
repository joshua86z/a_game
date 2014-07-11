package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"models"
	"protodata"
)

func getGeneralProto(general *models.GeneralData) protodata.GeneralData {

	config := models.ConfigGeneralMap()[general.ConfigId]

	var generalData protodata.GeneralData
	generalData.GeneralId = proto.Int32(int32(general.Id))
	generalData.GeneralName = proto.String(config.Name)
	generalData.GeneralDesc = proto.String(config.Desc)
	generalData.Level = proto.Int32(int32(general.Level))
	generalData.Atk = proto.Int32(int32(general.Atk))
	generalData.Def = proto.Int32(int32(general.Def))
	generalData.Hp = proto.Int32(int32(general.Hp))
	generalData.Speed = proto.Int32(int32(general.Speed))
	generalData.Range = proto.Int32(int32(general.Range))
	generalData.GeneralType = proto.Int32(int32(config.Type))
	generalData.LevelUpCoin = proto.Int32(int32(general.Level))
	generalData.IsUnlock = proto.Bool(false)
	return generalData
}

//

//
//// 武将相关接口
//type General struct{}
//
//// 初始化武将接口
//func NewGeneral() *General {
//	return &General{}
//}
//
//// 10121 获取角色武将列表
//func (g *General) GetGeneral(player *Player, cq *protodata.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	request := &pb.GetGeneralRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	// battleStatus 标识获取类型 1 未上阵，3 上阵，其他全部
//	battleStatus := request.GetGeneralClassType()
//	if battleStatus != int32(1) && battleStatus != int32(3) {
//		battleStatus = 0 // 非1或3，既取所有武将
//	}
//
//	// 武将详细信息
//	generalList := models.GetGeneralModel(uniqueId).GetList()
//
//	var generaldataList []*pb.GeneralData
//
//	for _, general := range generalList {
//		generalData, err := getGeneralData(uniqueId, general)
//		if err != nil {
//			return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//		}
//		generaldataList = append(generaldataList, generalData)
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.GetGeneralResponse{
//		ECode:        proto.Int32(1),
//		GeneralCount: proto.Int32(int32(len(generaldataList))),
//		Generals:     generaldataList,
//	}), nil
//}
//
//// 10136 获取单个角色武将
//func (g *General) GetSingleGeneral(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	request := &pb.GetSingleGeneralRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalData, err := getGeneralData(player.UniqueId, models.GetGeneralModel(player.UniqueId).GetGeneral(int(request.GetGeneralId())))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.GetSingleGeneralResponse{
//		ECode:   proto.Int32(1),
//		General: []*pb.GeneralData{generalData},
//	}), nil
//}
//
//// 10124 替换上阵武将
//func (this *General) SetActiveGeneral(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.SetActiveGeneralRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	// 武将相同，非法数据
//	oldGeneralId := int(request.GetOriginGeneralId())
//	newGeneralId := int(request.GetReplaceGeneralId())
//	if oldGeneralId == newGeneralId {
//		log.Warn("Unique:%d Ori gen equal with Rep gen", uniqueId)
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), errors.New(fmt.Sprintf("Unique:%d Ori gen equal with Rep gen", uniqueId))
//	}
//
//	GeneralModel := models.GetGeneralModel(uniqueId)
//
//	oldGeneral := GeneralModel.GetGeneral(oldGeneralId)
//	newGeneral := GeneralModel.GetGeneral(newGeneralId)
//
//	if !oldGeneral.IsUp() {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), errors.New(fmt.Sprintf("oldGeneral not up genId:%d", oldGeneralId))
//	}
//	if newGeneral.IsUp() {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), errors.New(fmt.Sprintf("newGeneral is up genId:%d", newGeneral))
//	}
//
//	isUpType := GeneralModel.GetIsUpSoldierType()
//
//	netSoldierType := newGeneral.UpSoldierType
//	for _, soldierType := range isUpType {
//		if soldierType == int(newGeneral.UpSoldierType) && soldierType != int(oldGeneral.UpSoldierType) {
//			netSoldierType = oldGeneral.UpSoldierType
//		}
//	}
//
//	err := oldGeneral.UpdateGeneralIsUp(false, oldGeneral.UpSoldierType)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//	err = newGeneral.UpdateGeneralIsUp(true, netSoldierType)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	response := &pb.SetActiveGeneralResponse{
//		ECode: proto.Int32(1),
//	}
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 10125 设置武将兵种
//func (g *General) SetGeneralSoldier(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.SetGeneralSoldierRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalId := int(request.GetGeneralId())
//	soldierId := request.GetReplaceSoldierId()
//
//	SoldierModel := models.NewSoldierModel(uniqueId)
//
//	solider := SoldierModel.GetSoldier(int(soldierId))
//
//	soldierType := solider.Type
//
//	GeneralModel := models.GetGeneralModel(uniqueId)
//
//	general := GeneralModel.GetGeneral(generalId)
//
//	var is_up bool
//	if general.IsUp() {
//		is_up = true
//		isUpType := GeneralModel.GetIsUpSoldierType()
//
//		for _, myType := range isUpType {
//			if int(myType) == int(soldierType) && int(soldierType) != int(general.UpSoldierType) {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("This soldier type is up genId:%d", generalId)
//			}
//		}
//	}
//
//	err := general.UpdateGeneralIsUp(is_up, int(soldierType))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SetGeneralSoldierResponse{
//		ECode:            proto.Int32(1),
//		ReplaceSoldierId: proto.Int32(soldierId),
//	}), nil
//}
//
//// 10126 升级武将兵种品质
//func (g *General) SoldierQualityUp(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.UpgradeSoldierLevelRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	SoldierModel := models.NewSoldierModel(uniqueId)
//
//	// 被升级兵种信息
//	roleSoldier := SoldierModel.GetSoldier(int(request.GetSoldierId()))
//
//	General := models.GetGeneralModel(uniqueId)
//	general := General.GetGeneral(roleSoldier.GeneralId)
//
//	configSoldier := configs.ConfigSoldierByArmsId(roleSoldier.ConfigId)
//
//	generalSoldierConfig, err := configs.ConfigGeneralSoldiersByGenId(general.GenId)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	// 判断等级
//	if configSoldier.Quality >= generalSoldierConfig[roleSoldier.Type].SoldierMaxLevel {
//		return pd.ReturnStr(cq, pb.StatusCode_OK,
//			&pb.UpgradeSoldierLevelResponse{
//				ECode:  proto.Int32(2),
//				EStr:   proto.String(ESTR_soldier_up_top),
//				BaseId: proto.Int32(int32(roleSoldier.ConfigId)),
//			}), nil
//	}
//
//	// 升级金钱
//	needCoin, err := gocalc.CalcFx(
//		models.FuncRuleMap[UPGRADE_SOL_COIN_FX].Fx,
//		map[string]interface{}{
//			"e": int32(configs.SoldierLevelUpNeedExp(int(configSoldier.Quality))) - int32(roleSoldier.Exp),
//		})
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_SERVER_INTERNAL_ERROR, ""), err
//	}
//
//	// 付钱升级
//	money := models.GetMoney(uniqueId)
//
//	if money.Coin < int(needCoin) {
//		return pd.ReturnStr(cq, pb.StatusCode_OK,
//			&pb.UpgradeSoldierLevelResponse{
//				ECode:  proto.Int32(3),
//				EStr:   proto.String(ESTR_not_enough_coin),
//				BaseId: proto.Int32(int32(roleSoldier.ConfigId)),
//			}), nil
//	}
//	err = money.SubCoin(int(needCoin), models.CO_UP_SOLD, fmt.Sprintf("soldier : %s , quality %d -> %d", configSoldier.Name, configSoldier.Quality, configSoldier.Quality+1))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if err = roleSoldier.LevelUp(); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	// 每日任务
//	models.NewMissionModel(uniqueId).GetMission(6).Add(1)
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK,
//		&pb.UpgradeSoldierLevelResponse{
//			ECode:  proto.Int32(1),
//			BaseId: proto.Int32(int32(roleSoldier.ConfigId)),
//		}), nil
//}
//
//// 10128 更换武将装备
//func (g *General) SetGeneralEquipment(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.SetGeneralEquipmentRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalId := int(request.GetGeneralId())
//	itemId := int(request.GetReplaceEquipmentsId())
//	equipType := int(request.GetEquipmentType())
//
//	//	var err error
//	General := models.GetGeneralModel(uniqueId)
//	general := General.GetGeneral(generalId)
//
//	//	props, err := models.GetRoleItemList(uniqueId)
//	//	if err != nil {
//	//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	//	}
//
//	ItemModel := models.NewItemModel(uniqueId)
//
//	for index, item := range ItemModel.ItemList {
//		if item.Type == configs.ItemType(equipType) && item.GeneralId == generalId {
//			if err := ItemModel.ItemList[index].UpDateGeneralEquip(0); err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//			}
//		}
//	}
//
//	if itemId != 0 {
//
//		equip := ItemModel.GetItem(itemId)
//		if equip.GeneralId != 0 {
//			return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("这件武器已被装备")
//		}
//
//		if err := equip.UpDateGeneralEquip(generalId); err != nil {
//			return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//		}
//	}
//
//	generalData, err := getGeneralData(uniqueId, general)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	response := &pb.SetGeneralEquipmentResponse{
//		ECode:   proto.Int32(1),
//		General: generalData,
//	}
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
////
////// 10134 检查技能刷新状态
////func (g *General) CheckRefreshSkill(player *Player, cq *pb.CommandRequest) (string, error) {
////
////	request := &pb.CheckRefreshSkillRequest{}
////	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
////		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
////	}
////
////	rgId := request.GetGeneralId()
////	generals, err := models.GetRoleGeneralList(player.UniqueId)
////	if traceErr("[general]", err) != nil {
////		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
////	}
////
////	for _, roleGen := range generals {
////		if int(rgId) == roleGen.Id {
////
////			remain := g.refreshRemainTime(&roleGen)
////
////			return pd.ReturnStr(cq, pb.StatusCode_OK,
////				&pb.CheckRefreshSkillResponse{
////					ECode:           proto.Int32(1),
////					RefreshLeftTime: proto.Int32(int32(remain)),
////					SeleIngot:       proto.Int32(timeToIngot(int32(remain))),
////				}), nil
////			break
////		}
////	}
////
////	return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), nil
////}
////
////// 10132 刷新武将技能
////func (g *General) RefreshGeneralSkill(player *Player, cq *pb.CommandRequest) (string, error) {
////
////	response := &pb.RefreshGeneralSkillResponse{
////		ECode:               proto.Int32(1),
////		SkillId:             proto.Int32(1),
////		NextRefreshTimeLeft: proto.Int32(1000),
////		Ingot:               proto.Int32(1),
////	}
////
////	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
////}
////
////// 判断技能刷新剩余时间
////// 返回刷新剩余时间，返回 = 0 可以立即刷新
////func (g *General) refreshRemainTime(rg *models.RoleGeneral) int {
////
////	// 获取上次刷新的时间
////	lastTime := int(rg.TechTime)
////
////	// 距离上次刷新已过去的时间
////	passTime := int(time.Now().Unix()) - lastTime
////
////	// 刷新时间已到
////	if passTime >= SkillWaitTime {
////		return 0
////	}
////
////	// 刷新间隔未到
////	return SkillWaitTime - passTime
////}
//
//// 10135 升级武将技能
//func (this *General) GeneralSkillLevelUp(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.SetSelectedSkillRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalId := int(request.GetGeneralId())
//	//	skillId := int(request.GetSkillId())
//
//	GeneralModel := models.GetGeneralModel(uniqueId)
//	general := GeneralModel.GetGeneral(generalId)
//
//	//	skill := general.Skills[skillId-1]
//
//	baseSkill := configs.ConfigSkillBySkillId(general.UsingSkillId)
//	needCoin, err := gocalc.CalcFx(baseSkill.CoinFx,
//		map[string]interface{}{
//			"l": general.SkillLevel,
//		})
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_SERVER_INTERNAL_ERROR, ""), err
//	}
//
//	money := models.GetMoney(uniqueId)
//
//	if money.Coin < int(needCoin) {
//
//		skillProto := &pb.SkillData{
//			Id:      proto.Int32(int32(general.Id)),
//			SkillId: proto.Int32(int32(general.UsingSkillId)),
//			Level:   proto.Int32(int32(general.SkillLevel)),
//			Name:    proto.String(baseSkill.Name),
//			Desc:    proto.String(""), // no desc
//		}
//
//		return pd.ReturnStr(cq, pb.StatusCode_OK,
//			&pb.UpgradeGeneralSkillResponse{
//				ECode: proto.Int32(2),
//				EStr:  proto.String(ESTR_not_enough_coin),
//				Skill: skillProto,
//			}), nil
//	}
//
//	if err = money.SubCoin(int(needCoin), models.CO_UP_SKIL, fmt.Sprintf("general : %s , level %d -> %d", general.Name, general.SkillLevel, general.SkillLevel+1)); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if err = general.SkillLevelUp(); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	skillProto := &pb.SkillData{
//		Id:      proto.Int32(int32(general.Id)),
//		SkillId: proto.Int32(int32(general.UsingSkillId)),
//		Level:   proto.Int32(int32(general.SkillLevel)),
//		Name:    proto.String(configs.ConfigSkillBySkillId(general.UsingSkillId).Name),
//		Desc:    proto.String(""), // no desc
//	}
//
//	// 每日任务
//	models.NewMissionModel(uniqueId).GetMission(7).Add(1)
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK,
//		&pb.UpgradeGeneralSkillResponse{
//			ECode: proto.Int32(1),
//			Skill: skillProto,
//		}), nil
//}
//
//func (this *General) GetPieceList(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	pieceList, err := models.GetRoleGeneralPieceList(player.UniqueId)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	var generalIds []int32
//	var pieceNum []int32
//
//	for _, configGeneral := range configs.ConfigGeneralsGetAll() {
//	SIGN:
//		for _, val := range pieceList {
//			if val.GenId == configGeneral.Id {
//				generalIds = append(generalIds, int32(val.GenId))
//				pieceNum = append(pieceNum, int32(val.Num))
//				continue SIGN
//			}
//		}
//		generalIds = append(generalIds, int32(configGeneral.Id))
//		pieceNum = append(pieceNum, 0)
//	}
//
//	response := &pb.GetChipListResponse{
//		ECode:      proto.Int32(1),
//		EStr:       proto.String(""),
//		GeneralsId: generalIds,
//		ChipNum:    pieceNum,
//	}
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//// 武将升阶
//func (this *General) GeneralClassUp(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.GeneralsAdvancedRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalId := int(request.GetGeneralId())
//
//	GeneralModel := models.GetGeneralModel(uniqueId)
//	general := GeneralModel.GetGeneral(generalId)
//
//	if general.Class >= general.GetMaxClass() {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("已经是最高阶不能再升")
//	}
//
//	//需要的碎片数 = =
//	needPiece := general.GetUpClassNeedPiece()
//
//	num := models.GetNumOfGeneralPiece(uniqueId, general.GenId)
//	if num < needPiece {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("碎片数量不够")
//	}
//
//	if err := models.SubGnerealPiece(uniqueId, general.GenId, needPiece); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	if err := general.UpClass(); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	num = models.GetNumOfGeneralPiece(uniqueId, general.GenId)
//
//	genProto, err := getGeneralData(uniqueId, general)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	money := models.GetMoney(uniqueId)
//
//	//获得武将成就
//	models.Achievement(uniqueId, 15)
//
//	response := &pb.GeneralsAdvanceResponse{
//		ECode:        proto.Int32(1),
//		GeneralId:    proto.Int32(int32(generalId)),
//		IsSuccess:    proto.Bool(true),
//		General:      genProto,
//		NewAdvance:   proto.Int32(int32(num)),
//		NewCoins:     proto.Int32(int32(money.Coin)),
//		ConsumeChips: proto.Int32(int32(needPiece)),
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//func (this *General) EatGenerals(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	request := &pb.EatGeneralsRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalId := int(request.GetGeneralId())
//	generalIds := request.GetBeEatGeneralId()
//
//	GeneralModel := models.GetGeneralModel(uniqueId)
//	general := GeneralModel.GetGeneral(generalId)
//
//	var addExp int
//	for _, otherId := range generalIds {
//		beEatGeneral := GeneralModel.GetGeneral(int(otherId))
//
//		// 武将转化的经验 ...
//		var tempExp int
//		for i := 1; i < beEatGeneral.Level; i++ {
//			tempLevelNeedExp, err := gocalc.CalcFx(
//				models.FuncRuleMap[UPGRADE_GEN_EXP_FX].Fx,
//				map[string]interface{}{
//					"l": i,
//				})
//			if err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//			}
//			tempExp += int(tempLevelNeedExp)
//		}
//
//		tempExp += beEatGeneral.Exp
//		floatExp, err := gocalc.CalcFx(
//			models.FuncRuleMap[22].Fx,
//			map[string]interface{}{
//				"e": tempExp, "g": beEatGeneral.Quality, "c": beEatGeneral.Class,
//			})
//		if err != nil {
//			return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//		}
//
//		addExp += int(floatExp)
//	}
//
//	ItemModel := models.NewItemModel(uniqueId)
//
//	var err error
//	var equipIdList []int
//
//	for _, generalId := range generalIds {
//
//		for _, item := range ItemModel.ItemList {
//			if item.GeneralId == int(generalId) {
//				if err = item.UpDateGeneralEquip(0); err != nil {
//					return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//				}
//				equipIdList = append(equipIdList, item.Id)
//			}
//		}
//
//		if err = models.DisableGeneral(uniqueId, int(generalId)); err != nil {
//			return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//		}
//	}
//
//	if err := general.AddExp(addExp); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	generalData, err := getGeneralData(uniqueId, general)
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	var propDataList []*pb.RoleProps
//
//	for _, equipId := range equipIdList {
//		if propData, err := getItemData(ItemModel.GetItem(equipId), configs.ConfigItemById(ItemModel.GetItem(equipId).ConfigId)); err != nil {
//			return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//		} else {
//			propDataList = append(propDataList, propData)
//		}
//	}
//
//	response := &pb.EatGeneralsResponse{
//		ECode:          proto.Int32(1),
//		EStr:           proto.String(""),
//		BeEatGeneralId: generalIds,
//		RewardExp:      proto.Int32(int32(addExp)),
//		ActiveGeneral:  generalData,
//		Props:          propDataList,
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//func getGeneralData(uniqueId int64, general *models.GeneralData) (*pb.GeneralData, error) {
//
//	//	general, err := models.GetGeneralById(uniqueId, generalId)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//
//	ItemModel := models.NewItemModel(uniqueId)
//
//	genConfig := configs.ConfigGeneralByGenId(int(general.GenId))
//	rgWrapper := wrapper.RoleGeneralWrapper{
//		GeneralData: general,
//		Base:        genConfig,
//	}
//
//	for _, item := range ItemModel.ItemList {
//		if item.GeneralId == int(general.Id) {
//			configItem := configs.ConfigItemById(int(item.ConfigId))
//			rpWrapper := wrapper.RolePropWrapper{
//				item,
//				configItem,
//			}
//
//			rpWrapper.Enhance(&rgWrapper)
//		}
//	}
//	return rgWrapper.ToProto(), nil
//}
//
//func getBattleGeneralData(uniqueId int64) ([]*pb.GeneralData, error) {
//
//	var result []*pb.GeneralData
//
//	generals := models.GetGeneralModel(uniqueId).GetBattleGeneralList()
//
//	for _, general := range generals {
//
//		data, err := getGeneralData(uniqueId, general)
//		if err != nil {
//			return result, err
//		}
//
//		result = append(result, data)
//	}
//
//	return result, nil
//}
