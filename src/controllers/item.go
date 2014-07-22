package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) ItemLevelUp() error {

	request := &protodata.ItemLevelUpRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	configId := int(request.GetItemId())

	config := models.ConfigItemList()[configId-1]
	ItemModel := models.NewItemModel(this.Role.Uid)

	var coin, level int
	item := ItemModel.Item(configId)
	if item == nil {
		level = 0
	} else {
		level = item.Level
	}
	coin = levelUpCoin(level)

	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	err := this.Role.SubCoin(coin, models.ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", config.Name, level, level+1))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	if item == nil {
		if item = ItemModel.Insert(config); item == nil {
			return this.Send(lineNum(), err)
		}
	} else {
		if err = item.LevelUp(); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	response := &protodata.ItemLevelUpResponse{
		Role: roleProto(this.Role),
		Item: itemProto(item, config),
	}
	return this.Send(StatusOK, response)
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
