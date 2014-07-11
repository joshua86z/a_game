package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"models"
	"protodata"
)

func getItemProto(item *models.ItemData) protodata.ItemData {

	config := models.ConfigItemMap()[item.Id]

	var itemData protodata.ItemData
	itemData.ItemId = proto.Int32(int32(item.Id))
	itemData.ItemName = proto.String(config.Name)
	itemData.ItemDesc = proto.String(config.Desc)
	itemData.Level = proto.Int32(int32(item.Level))
	itemData.LevelUpCoin = proto.Int32(int32(item.LevelUpCoin()))
	return itemData
}

//type Item struct {
//	itemMaxLevel map[int]int
//}
//
//func NewItem() *Item {
//
//	itemMaxLevel := make(map[int]int)
//	itemMaxLevel[5] = 90
//	itemMaxLevel[4] = 75
//	itemMaxLevel[3] = 60
//	itemMaxLevel[2] = 45
//	itemMaxLevel[1] = 30
//
//	return &Item{itemMaxLevel: itemMaxLevel}
//}
//
//// 10151 获取玩家装备列表
//func (p *Item) GetItemList(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.GetPropsRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	request_type := request.GetPropType()
//
//	itemList := models.NewItemModel(uniqueId).ItemList
//
//	for i := 0; i < len(itemList); i++ {
//		for j := len(itemList) - 1; j > i; j-- {
//			if itemList[j].GeneralId > 0 && itemList[j-1].GeneralId == 0 {
//				temp := itemList[j-1]
//				itemList[j-1] = itemList[j]
//				itemList[j] = temp
//			} else {
//				if configs.ConfigItemById(itemList[j].ConfigId).Quality > configs.ConfigItemById(itemList[j-1].ConfigId).Quality {
//					temp := itemList[j-1]
//					itemList[j-1] = itemList[j]
//					itemList[j] = temp
//				}
//			}
//		}
//	}
//
//	itemDataList := []*pb.RoleProps{}
//	rbs := []*pb.RoleBook{}
//
//	for _, item := range itemList {
//
//		if request_type == int32(item.Type) || request_type == 0 {
//			if itemData, err := getItemData(item, configs.ConfigItemById(item.ConfigId)); err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//			} else {
//				itemDataList = append(itemDataList, itemData)
//			}
//		}
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.GetPropsResponse{
//		ECode:      proto.Int32(1),
//		PropsCount: proto.Int32(int32(len(itemDataList))),
//		Props:      itemDataList,
//		BookCount:  proto.Int32(int32(len(rbs))),
//		Books:      rbs,
//	}), nil
//}
//
//// 10153 卖道具
//func (p *Item) SaleItem(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.SalePropertyRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	itemId := int(request.GetId())
//
//	// 判断道具是否属于角色的
//	ItemModel := models.NewItemModel(uniqueId)
//	item := ItemModel.GetItem(itemId)
//
//	if item.Id == 0 || !InSlice(int(item.Type), []interface{}{1, 2, 3, 4}) {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SalePropertyResponse{
//			ECode: proto.Int32(2),
//			EStr:  proto.String(ESTR_cant_sale_prop),
//			Coins: proto.Int32(0),
//		}), nil
//	}
//
//	//判断道具是否被穿戴
//	if item.GeneralId != 0 {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SalePropertyResponse{
//			ECode: proto.Int32(3),
//			EStr:  proto.String(ESTR_item_prop_used),
//			Coins: proto.Int32(0),
//		}), nil
//	}
//
//	//计算道具卖出价格 升级所需的总金币 * 0.7
//	addCoin := item.GetSaleCoin()
//
//	// 出售道具
//
//	if err := ItemModel.DeleteItem(itemId); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	money := models.GetMoney(uniqueId)
//	money.AddCoin(addCoin, models.CO_SALE_PR, fmt.Sprintf("item : %v , level : %d", configs.ConfigItemById(item.ConfigId), item.Level))
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SalePropertyResponse{
//		ECode: proto.Int32(1),
//		Coins: proto.Int32(int32(addCoin)),
//	}), nil
//}
//
//// 10155 升级装备
//func (this *Item) EquipmentLevelUp(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//	request := &pb.UpgradeEquipmentRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//
//	itemId := int(request.GetId())
//	otherEquipIds := request.GetMixEquipmentId()
//
//	ItemModel := models.NewItemModel(uniqueId)
//	item := ItemModel.GetItem(itemId)
//
//	configItem := configs.ConfigItemById(item.ConfigId)
//	maxLevel := this.itemMaxLevel[configItem.Quality]
//
//	if item.Level >= maxLevel {
//
//		itemData, _ := getItemData(item, configItem)
//
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//			ECode:          proto.Int32(3),
//			EStr:           proto.String(ESTR_prop_up_top),
//			Id:             proto.Int32(int32(item.Id)),
//			NewLevel:       proto.Int32(int32(item.Level)),
//			MixEquipmentId: otherEquipIds,
//			Props:          itemData,
//		}), fmt.Errorf("Equip is max Level , itemId : %d", item.Id)
//	}
//
//	if len(otherEquipIds) > 0 {
//
//		if item.Type != configs.EQUIP_TYPE_BOOK {
//
//			if item.Level%10 != 0 {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("装备等级不是10的倍数不能进阶")
//			}
//
//			if len(otherEquipIds) != 2 {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("合成需要三个相同的武器")
//			}
//		} else {
//			if len(otherEquipIds) != 3 {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("合成需要四个相同的兵书")
//			}
//		}
//
//		for _, tempItemId := range otherEquipIds {
//			temp := ItemModel.GetItem(int(tempItemId))
//
//			if temp.ConfigId != item.ConfigId {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("不是相同的类型不能合成")
//			}
//		}
//		for _, tempItemId := range otherEquipIds {
//
//			if err := ItemModel.DeleteItem(int(tempItemId)); err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//			}
//		}
//
//	} else {
//		if item.Type == configs.EQUIP_TYPE_BOOK {
//			return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("兵书不能升级")
//		}
//	}
//
//	// 计算升级金币
//	needCoin := models.GetItemLevelUpCoin(item.Level)
//
//	money := models.GetMoney(uniqueId)
//
//	if money.Coin < needCoin {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//			ECode:          proto.Int32(2),
//			EStr:           proto.String(ESTR_not_enough_coin),
//			Id:             proto.Int32(int32(item.Id)),
//			NewLevel:       proto.Int32(int32(item.Level)),
//			MixEquipmentId: otherEquipIds,
//		}), nil
//	}
//
//	if err := money.SubCoin(needCoin, models.CO_UP_ITEM, fmt.Sprintf("item : %s , level %d -> %d", configs.ConfigItemById(item.ConfigId).Name, item.Level, item.Level+1)); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	//	newLevel := item.Level + 1
//	if err := item.LevelUp(1); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	itemData, err := getItemData(item, configs.ConfigItemById(item.ConfigId))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	// 每日任务
//	models.NewMissionModel(uniqueId).GetMission(5).Add(1)
//
//	if item.Type == configs.EQUIP_TYPE_BOOK {
//		// 兵书成就
//		models.Achievement(uniqueId, 14)
//	}
//
//	var generalData *pb.GeneralData
//	if item.GeneralId > 0 {
//		generalData, _ = getGeneralData(uniqueId, models.GetGeneralModel(uniqueId).GetGeneral(item.GeneralId))
//	}
//
//	// 升级成功
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//		ECode:          proto.Int32(1),
//		EStr:           proto.String(""),
//		Id:             proto.Int32(int32(item.Id)),
//		NewLevel:       proto.Int32(int32(item.Level)),
//		MixEquipmentId: otherEquipIds,
//		Props:          itemData,
//		Coins:          proto.Int32(int32(money.Coin)),
//		General:        generalData,
//	}), nil
//}
//
//// 10154 使用礼包
//func (p *Item) UseGiftBag(player *Player, cq *pb.CommandRequest) (string, error) {
//
//	uniqueId := player.UniqueId
//
//	request := &pb.UseGiftBagRequest{}
//	if err := pd.Unmarshal(cq.GetSerializedString(), request); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), err
//	}
//	giftId := int(request.GetId())
//
//	ItemModel := models.NewItemModel(uniqueId)
//	gift := ItemModel.GetItem(giftId)
//
//	configGift := configs.ConfigItemById(gift.ConfigId)
//
//	if gift.Id == 0 || gift.Type != 5 {
//		return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, "道具 : "+configGift.Name+", 不是礼包类型"), fmt.Errorf("道具 : " + configGift.Name + ", 不是礼包类型")
//	}
//
//	coins, randContents := (wrapper.GiftWrapper{configGift}).GetContent()
//
//	// 获得钱
//	if coins > 0 {
//		if err := models.GetMoney(uniqueId).AddGold(int(coins), models.CO_VIP_BAG, "item : "+configGift.Name); err != nil {
//			coins = 0
//		}
//	}
//
//	// 获得礼包内道具
//	var propIdList []int32
//	var baseIdList []int32
//	var qualityList []int32
//	var peiceGenId int
//	var peiceNum int
//	for _, content := range randContents {
//
//		itemConfigId := int(content[0])
//		itemLevel := int(content[1])
//
//		configItem := configs.ConfigItemById(itemConfigId)
//
//		if configItem.Type == configs.ITEM_TYPE_PIECE {
//			models.AddGnerealPiece(uniqueId, configItem.GenId, itemLevel)
//			peiceGenId = configItem.GenId
//			peiceNum = itemLevel
//			continue
//		}
//
//		item, err := ItemModel.InsertItem(configItem, itemLevel)
//		if err != nil {
//			continue
//		}
//
//		propIdList = append(propIdList, int32(item.Id))
//		baseIdList = append(baseIdList, int32(itemConfigId))
//		qualityList = append(qualityList, int32(configItem.Quality))
//	}
//
//	// Delete item prop
//	if err := ItemModel.DeleteItem(giftId); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	//更新道具类的成就
//	models.Achievement(uniqueId, 12, 13, 24)
//
//	response := &pb.UseGiftBagResponse{
//		ECode:          proto.Int32(1),
//		PropsCount:     proto.Int32(int32(len(propIdList))),
//		Baseid:         baseIdList,
//		Id:             propIdList,
//		PropQuality:    qualityList,
//		GeneralsChipId: proto.Int32(int32(peiceGenId)),
//		ChipNum:        proto.Int32(int32(peiceNum)),
//	}
//	return pd.ReturnStr(cq, pb.StatusCode_OK, response), nil
//}
//
//func getItemData(item *models.ItemData, itemConfig *configs.Config_Item) (*pb.RoleProps, error) {
//
//	return wrapper.RolePropWrapper{item, itemConfig}.ToProto(), nil
//}
