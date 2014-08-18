package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"models"
	"protodata"
)

func (this *Connect) ItemLevelUp() error {

	request := new(protodata.ItemLevelUpRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	baseId := int(request.GetItemId())

	baseItem := models.BaseItem(baseId, nil)
	if baseItem == nil {
		return this.Send(lineNum(), fmt.Errorf("参数错误:没有这个道具Id:%d", baseId))
	}

	var coin, level int
	item := models.Item.Item(this.Uid, baseId)
	if item == nil {
		level = 0
	} else {
		level = item.Level
	}

	if level >= len(baseItem.LevelUpCoin)-1 {
		return this.Send(lineNum(), fmt.Errorf("道具已经最大等级"))
	}

	coin = baseItem.LevelUpCoin[level]
	if coin > this.Role.Coin {
		return this.Send(lineNum(), fmt.Errorf("金币不足"))
	}

	err := this.Role.SubCoin(coin, models.FINANCE_ITEM_LEVELUP, fmt.Sprintf("item: %s , level: %d -> %d", baseItem.Name, level, level+1))
	if err != nil {
		return this.Send(lineNum(), err)
	}

	if item == nil {
		if item = models.Item.Insert(this.Uid, baseItem); item == nil {
			return this.Send(lineNum(), err)
		}
	} else {
		if err = item.LevelUp(); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	response := &protodata.ItemLevelUpResponse{
		Role: roleProto(this.Role),
		Item: itemProto(item, baseItem),
	}
	return this.Send(StatusOK, response)
}

func itemProtoList(itemList []*models.ItemData) []*protodata.ItemData {

	var result []*protodata.ItemData
	for _, config := range models.BaseItemList() {

		var find bool
		for _, item := range itemList {
			if item.BaseId == config.BaseId {
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

func itemProto(item *models.ItemData, config *models.Base_Item) *protodata.ItemData {

	return &protodata.ItemData{
		ItemId:      proto.Int32(int32(config.BaseId)),
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

func tempItemDiamond() [4]int {
	Lua, _ := lua.NewLua("conf/item.lua")
	Lua.L.GetGlobal("tempItemDiamond")
	Lua.L.DoString("d1,d2,d3,d4 = tempItemDiamond()")
	d1, d2, d3, d4 := Lua.GetInt("d1"), Lua.GetInt("d2"), Lua.GetInt("d3"), Lua.GetInt("d4")
	Lua.Close()
	return [4]int{d1, d2, d3, d4}
}
