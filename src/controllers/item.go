package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"models"
	"protodata"
	"strings"
)

func (this *Connect) ItemLevelUp() error {

	request := new(protodata.ItemLevelUpRequest)
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

	if level >= 5 {
		return this.Send(lineNum(), fmt.Errorf("道具已经最大等级"))
	}

	coin = levelUpCoinMap()[level]
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	err := this.Role.SubCoin(coin, models.FINANCE_ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", config.Name, level, level+1))
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

	m := levelUpCoinMap()

	var result []*protodata.ItemData
	for _, config := range models.ConfigItemList() {

		itemData := new(protodata.ItemData)
		itemData.ItemId = proto.Int32(int32(config.ConfigId))
		itemData.ItemName = proto.String(config.Name)
		itemData.ItemDesc = proto.String(config.Desc)
		itemData.LevelUpCoin = proto.Int32(int32(m[0]))
		itemData.ItemValue = proto.Int32(int32(config.Value))
		itemData.ItemPro = proto.Int32(int32(config.Probability))

		for _, item := range itemList {
			if item.ConfigId == config.ConfigId {
				itemData.Level = proto.Int32(int32(item.Level))
				itemData.LevelUpCoin = proto.Int32(int32(m[item.Level]))
				itemData.ItemValue = proto.Int32(int32(config.Value + item.Level*config.Group))
				break
			}
		}

		result = append(result, itemData)
	}
	return result
}

func itemProto(item *models.ItemData, config *models.ConfigItem) *protodata.ItemData {

	m := levelUpCoinMap()

	itemData := new(protodata.ItemData)
	itemData.ItemId = proto.Int32(int32(config.ConfigId))
	itemData.ItemName = proto.String(config.Name)
	itemData.ItemDesc = proto.String(config.Desc)
	itemData.ItemValue = proto.Int32(int32(config.Value + item.Level*config.Group))
	itemData.ItemPro = proto.Int32(int32(config.Probability))
	itemData.Level = proto.Int32(int32(item.Level))
	itemData.LevelUpCoin = proto.Int32(int32(m[item.Level]))
	return itemData
}

func levelUpCoinMap() map[int]int {

	Lua, _ := lua.NewLua("conf/item.lua")
	s := Lua.GetString("levelUp")
	array := strings.Split(s, ",")

	result := make(map[int]int)
	for index, val := range array {
		result[index] = models.Atoi(val)
	}

	return result
}
