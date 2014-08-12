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

	var config *models.ConfigItem
	for _, c := range models.ConfigItemList() {
		if c.ConfigId == configId {
			config = c
			break
		}
	}
	if config == nil {
		return this.Send(lineNum(), fmt.Errorf("参数错误:没有这个道具Id:%d", configId))
	}

	var coin, level int
	item := models.Item.Item(this.Uid, configId)
	if item == nil {
		level = 0
	} else {
		level = item.Level
	}

	if level >= len(config.LevelUpCoin)-1 {
		return this.Send(lineNum(), fmt.Errorf("道具已经最大等级"))
	}

	coin = config.LevelUpCoin[level]
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	err := this.Role.SubCoin(coin, models.FINANCE_ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", config.Name, level, level+1))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	if item == nil {
		if item = models.Item.Insert(this.Uid, config); item == nil {
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

		var find bool
		for _, item := range itemList {
			if item.ConfigId == config.ConfigId {
				find = true
				result = append(result, itemProto(item, config))
				break
			}
		}
		if !find {
			result = append(result, itemProto(new(models.ItemData), config))
		}
	}
	return result
}

func itemProto(item *models.ItemData, config *models.ConfigItem) *protodata.ItemData {

	return &protodata.ItemData{
		ItemId:      proto.Int32(int32(config.ConfigId)),
		ItemName:    proto.String(config.Name),
		ItemDesc:    proto.String(config.Desc),
		ItemValue:   proto.Int32(int32(config.Value + item.Level*config.Group)),
		ItemPro:     proto.Int32(int32(config.Probability)),
		Level:       proto.Int32(int32(item.Level)),
		LevelUpCoin: proto.Int32(int32(config.LevelUpCoin[item.Level]))}
}

//func levelUpCoinMap() map[int]int {
//
//	Lua, _ := lua.NewLua("conf/item.lua")
//	s := Lua.GetString("levelUp")
//	array := strings.Split(s, ",")
//
//	result := make(map[int]int)
//	for index, val := range array {
//		result[index] = models.Atoi(val)
//	}
//
//	return result
//}

func tempItemCoin() [4]int {
	Lua, _ := lua.NewLua("conf/item.lua")
	s := Lua.GetString("temp_item_coin")
	Lua.Close()
	tempItemCoin := strings.Split(s, ",")
	return [4]int{models.Atoi(tempItemCoin[0]), models.Atoi(tempItemCoin[1]), models.Atoi(tempItemCoin[2]), models.Atoi(tempItemCoin[3])}
}
