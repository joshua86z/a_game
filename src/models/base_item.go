package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

func init() {
	BaseItemList()
}

type Base_Item struct {
	BaseId      int
	Name        string
	Desc        string
	Value       int
	Group       int
	Probability int
	LevelUpCoin []int
}

func BaseItemList() []*Base_Item {

	var result []*Base_Item

	Lua, _ := lua.NewLua("conf/item.lua")
	//Lua.L.GetGlobal("item")
	//Lua.L.GetGlobal("levelUpCoin")
	maxLevel := Lua.GetInt("max_level")

	var coinList []int
	for i := 0; i <= maxLevel; i++ {
		Lua.L.DoString(fmt.Sprintf("c = levelUpCoin(%d)", i))
		coinList = append(coinList, Lua.GetInt("c"))
	}

	indexStr := Lua.GetString("index")
	indexArr := strings.Split(indexStr, ",")
	for _, index := range indexArr {
		baseId := Atoi(index)
		item := BaseItem(baseId, Lua)
		item.LevelUpCoin = coinList
		result = append(result, item)
	}

	Lua.Close()
	return result
}

func BaseItem(baseId int, Lua *lua.Lua) *Base_Item {

	var sign bool
	if Lua == nil {
		sign = true
		Lua, _ = lua.NewLua("conf/item.lua")
		defer Lua.Close()
	}

	Lua.L.DoString(fmt.Sprintf("name, desc, value, group, probability = item(%d)", baseId))
	name, desc, value, group, probability := Lua.GetString("name"), Lua.GetString("desc"), Lua.GetInt("value"), Lua.GetInt("group"), Lua.GetInt("probability")
	if name == "" {
		return nil
	}

	var coinList []int
	if sign {
		maxLevel := Lua.GetInt("max_level")

		for i := 0; i <= maxLevel; i++ {
			Lua.L.DoString(fmt.Sprintf("c = levelUpCoin(%d)", i))
			coinList = append(coinList, Lua.GetInt("c"))
		}
	}

	return &Base_Item{
		BaseId:      baseId,
		Name:        name,
		Desc:        desc,
		Value:       value,
		Group:       group,
		Probability: probability,
		LevelUpCoin: coinList}
}
