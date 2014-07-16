package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"models"
	"protodata"
)

func getMailProto(mail *models.MailData) protodata.MailData {

	var mailData protodata.MailData
	mailData.MailId = proto.Int32(int32(mail.Id))
	mailData.MailTitle = proto.String(mail.Title)
	mailData.MailContent = proto.String(mail.Content)
	mailData.Reward = nil
	mailData.IsReceive = proto.Bool(mail.IsReceive)
	return mailData
}

//type Item struct {
//	mailMaxLevel map[int]int
//}
//
//func NewItem() *Item {
//
//	mailMaxLevel := make(map[int]int)
//	mailMaxLevel[5] = 90
//	mailMaxLevel[4] = 75
//	mailMaxLevel[3] = 60
//	mailMaxLevel[2] = 45
//	mailMaxLevel[1] = 30
//
//	return &Item{mailMaxLevel: mailMaxLevel}
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
//	mailList := models.NewItemModel(uniqueId).ItemList
//
//	for i := 0; i < len(mailList); i++ {
//		for j := len(mailList) - 1; j > i; j-- {
//			if mailList[j].GeneralId > 0 && mailList[j-1].GeneralId == 0 {
//				temp := mailList[j-1]
//				mailList[j-1] = mailList[j]
//				mailList[j] = temp
//			} else {
//				if configs.ConfigItemById(mailList[j].ConfigId).Quality > configs.ConfigItemById(mailList[j-1].ConfigId).Quality {
//					temp := mailList[j-1]
//					mailList[j-1] = mailList[j]
//					mailList[j] = temp
//				}
//			}
//		}
//	}
//
//	mailDataList := []*pb.RoleProps{}
//	rbs := []*pb.RoleBook{}
//
//	for _, mail := range mailList {
//
//		if request_type == int32(mail.Type) || request_type == 0 {
//			if mailData, err := getItemData(mail, configs.ConfigItemById(mail.ConfigId)); err != nil {
//				return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//			} else {
//				mailDataList = append(mailDataList, mailData)
//			}
//		}
//	}
//
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.GetPropsResponse{
//		ECode:      proto.Int32(1),
//		PropsCount: proto.Int32(int32(len(mailDataList))),
//		Props:      mailDataList,
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
//	mailId := int(request.GetId())
//
//	// 判断道具是否属于角色的
//	ItemModel := models.NewItemModel(uniqueId)
//	mail := ItemModel.GetItem(mailId)
//
//	if mail.Id == 0 || !InSlice(int(mail.Type), []interface{}{1, 2, 3, 4}) {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SalePropertyResponse{
//			ECode: proto.Int32(2),
//			EStr:  proto.String(ESTR_cant_sale_prop),
//			Coins: proto.Int32(0),
//		}), nil
//	}
//
//	//判断道具是否被穿戴
//	if mail.GeneralId != 0 {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.SalePropertyResponse{
//			ECode: proto.Int32(3),
//			EStr:  proto.String(ESTR_mail_prop_used),
//			Coins: proto.Int32(0),
//		}), nil
//	}
//
//	//计算道具卖出价格 升级所需的总金币 * 0.7
//	addCoin := mail.GetSaleCoin()
//
//	// 出售道具
//
//	if err := ItemModel.DeleteItem(mailId); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	money := models.GetMoney(uniqueId)
//	money.AddCoin(addCoin, models.CO_SALE_PR, fmt.Sprintf("mail : %v , level : %d", configs.ConfigItemById(mail.ConfigId), mail.Level))
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
//	mailId := int(request.GetId())
//	otherEquipIds := request.GetMixEquipmentId()
//
//	ItemModel := models.NewItemModel(uniqueId)
//	mail := ItemModel.GetItem(mailId)
//
//	configItem := configs.ConfigItemById(mail.ConfigId)
//	maxLevel := this.mailMaxLevel[configItem.Quality]
//
//	if mail.Level >= maxLevel {
//
//		mailData, _ := getItemData(mail, configItem)
//
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//			ECode:          proto.Int32(3),
//			EStr:           proto.String(ESTR_prop_up_top),
//			Id:             proto.Int32(int32(mail.Id)),
//			NewLevel:       proto.Int32(int32(mail.Level)),
//			MixEquipmentId: otherEquipIds,
//			Props:          mailData,
//		}), fmt.Errorf("Equip is max Level , mailId : %d", mail.Id)
//	}
//
//	if len(otherEquipIds) > 0 {
//
//		if mail.Type != configs.EQUIP_TYPE_BOOK {
//
//			if mail.Level%10 != 0 {
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
//			if temp.ConfigId != mail.ConfigId {
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
//		if mail.Type == configs.EQUIP_TYPE_BOOK {
//			return pd.ReturnStr(cq, pb.StatusCode_DATA_ERROR, ""), fmt.Errorf("兵书不能升级")
//		}
//	}
//
//	// 计算升级金币
//	needCoin := models.GetItemLevelUpCoin(mail.Level)
//
//	money := models.GetMoney(uniqueId)
//
//	if money.Coin < needCoin {
//		return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//			ECode:          proto.Int32(2),
//			EStr:           proto.String(ESTR_not_enough_coin),
//			Id:             proto.Int32(int32(mail.Id)),
//			NewLevel:       proto.Int32(int32(mail.Level)),
//			MixEquipmentId: otherEquipIds,
//		}), nil
//	}
//
//	if err := money.SubCoin(needCoin, models.CO_UP_ITEM, fmt.Sprintf("mail : %s , level %d -> %d", configs.ConfigItemById(mail.ConfigId).Name, mail.Level, mail.Level+1)); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	//	newLevel := mail.Level + 1
//	if err := mail.LevelUp(1); err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	mailData, err := getItemData(mail, configs.ConfigItemById(mail.ConfigId))
//	if err != nil {
//		return pd.ReturnStr(cq, pb.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	// 每日任务
//	models.NewMissionModel(uniqueId).GetMission(5).Add(1)
//
//	if mail.Type == configs.EQUIP_TYPE_BOOK {
//		// 兵书成就
//		models.Achievement(uniqueId, 14)
//	}
//
//	var generalData *pb.GeneralData
//	if mail.GeneralId > 0 {
//		generalData, _ = getGeneralData(uniqueId, models.GetGeneralModel(uniqueId).GetGeneral(mail.GeneralId))
//	}
//
//	// 升级成功
//	return pd.ReturnStr(cq, pb.StatusCode_OK, &pb.UpgradeEquipmentResponse{
//		ECode:          proto.Int32(1),
//		EStr:           proto.String(""),
//		Id:             proto.Int32(int32(mail.Id)),
//		NewLevel:       proto.Int32(int32(mail.Level)),
//		MixEquipmentId: otherEquipIds,
//		Props:          mailData,
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
//		if err := models.GetMoney(uniqueId).AddGold(int(coins), models.CO_VIP_BAG, "mail : "+configGift.Name); err != nil {
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
//		mailConfigId := int(content[0])
//		mailLevel := int(content[1])
//
//		configItem := configs.ConfigItemById(mailConfigId)
//
//		if configItem.Type == configs.ITEM_TYPE_PIECE {
//			models.AddGnerealPiece(uniqueId, configItem.GenId, mailLevel)
//			peiceGenId = configItem.GenId
//			peiceNum = mailLevel
//			continue
//		}
//
//		mail, err := ItemModel.InsertItem(configItem, mailLevel)
//		if err != nil {
//			continue
//		}
//
//		propIdList = append(propIdList, int32(mail.Id))
//		baseIdList = append(baseIdList, int32(mailConfigId))
//		qualityList = append(qualityList, int32(configItem.Quality))
//	}
//
//	// Delete mail prop
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
//func getItemData(mail *models.ItemData, mailConfig *configs.Config_Item) (*pb.RoleProps, error) {
//
//	return wrapper.RolePropWrapper{mail, mailConfig}.ToProto(), nil
//}
