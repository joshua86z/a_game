package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Item struct {
}

func (this *Item) LevelUp(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	request := &protodata.ItemLevelUpRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return lineNum(), nil, err
	}

	configId := int(request.GetItemId())

	config := models.ConfigItemList()[configId-1]
	ItemModel := models.NewItemModel(RoleModel.Uid)

	var coin, level int
	item := ItemModel.Item(configId)
	if item == nil {
		level = 0
	} else {
		level = item.Level
	}
	coin = levelUpCoin(level)

	if coin > RoleModel.Coin {
		return lineNum(), nil, fmt.Errorf("金币不足")
	}

	err := RoleModel.SubCoin(coin, models.ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", config.Name, level, level+1))
	if err != nil {
		return lineNum(), nil, err
	}

	if item == nil {
		if item = ItemModel.Insert(config); item == nil {
			return lineNum(), nil, err
		}
	} else {
		if err = item.LevelUp(); err != nil {
			return lineNum(), nil, err
		}
	}

	response := &protodata.ItemLevelUpResponse{
		Role: roleProto(RoleModel),
		Item: itemProto(item, config),
	}
	return protodata.StatusCode_OK, response, nil
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
