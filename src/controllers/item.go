package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Item struct {
}

func (this *Item) LevelUp(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	request := &protodata.ItemLevelUpRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return ReturnStr(commandRequest, 17, ""), err
	}

	configId := int(request.GetItemId())

	config := models.ConfigItemList()[configId-1]
	ItemModel := models.NewItemModel(uid)
	RoleModel := models.NewRoleModel(uid)

	var coin, level int
	item := ItemModel.Item(configId)
	if item == nil {
		level = 0
	} else {
		level = item.Level
	}
	coin = levelUpCoin(level)

	if coin > RoleModel.Coin {
		return ReturnStr(commandRequest, 36, "金币不足"), fmt.Errorf("金币不足")
	}

	err := RoleModel.SubCoin(coin, models.ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", config.Name, level, level+1))
	if err != nil {
		return ReturnStr(commandRequest, 41, "失败,数据库错误"), err
	}

	if item == nil {
		if item = ItemModel.Insert(config); item == nil {
			return ReturnStr(commandRequest, 46, "失败,数据库错误"), err
		}
	} else {
		if err = item.LevelUp(); err != nil {
			return ReturnStr(commandRequest, 50, "失败,数据库错误"), err
		}
	}

	response := &protodata.ItemLevelUpResponse{
		Role: roleProto(RoleModel),
		Item: itemProto(item, config),
	}
	return ReturnStr(commandRequest, protodata.StatusCode_OK, response), nil
}

func itemProtoList(itemList []*models.ItemData) []*protodata.ItemData {

	var result []*protodata.ItemData
	for _, config := range models.ConfigItemList() {

		var itemData protodata.ItemData
		itemData.ItemId = proto.Int32(int32(config.ConfigId))
		itemData.ItemName = proto.String(config.Name)
		itemData.ItemDesc = proto.String(config.Desc)
		itemData.LevelUpCoin = proto.Int32(int32(levelUpCoin(1)))
		itemData.Level = proto.Int32(1)

		for _, item := range itemList {
			if item.ConfigId == config.ConfigId {
				itemData.Level = proto.Int32(int32(item.Level))
				itemData.LevelUpCoin = proto.Int32(int32(levelUpCoin(item.Level)))
				break
			}
		}

		result = append(result, &itemData)
	}
	return result
}

func levelUpCoin(level int) int {
	return level
}

func itemProto(item *models.ItemData, config *models.ConfigItem) *protodata.ItemData {

	var itemData protodata.ItemData
	itemData.ItemId = proto.Int32(int32(config.ConfigId))
	itemData.ItemName = proto.String(config.Name)
	itemData.ItemDesc = proto.String(config.Desc)
	itemData.Level = proto.Int32(int32(item.Level))
	itemData.LevelUpCoin = proto.Int32(int32(levelUpCoin(item.Level)))
	return &itemData
}
