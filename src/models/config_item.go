package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

var (
	item_config_map  map[int]*ConfigItem
	item_config_list []*ConfigItem
)

func init() {
	var temp []*ConfigItem
	if _, err := DB().Select(&temp, "SELECT * FROM config_item"); err != nil {
		panic(err)
	}
	item_config_map = make(map[int]*ConfigItem)
	for _, item := range temp {
		item_config_map[item.ConfigId] = item
	}
}

// config_item
type ConfigItem struct {
	ConfigId int    `db:"item_config_id"`
	Name     string `db:"item_name"`
	Desc     string `db:"item_desc"`
}

func ConfigItemMap() map[int]*ConfigItem {
	return item_config_map
}

func ConfigItemList() []*ConfigItem {

	var result []*ConfigItem

	Lua, _ := lua.NewLua("conf/item.lua")
	i := 1
	for {
		itemStr := Lua.GetString(fmt.Sprintf("item_%d", i))
		if itemStr == "" {
			break
		} else {
			i++
		}
		array := strings.Split(itemStr, "\\,")
		result = append(result, &ConfigItem{i, array[0], array[1]})
	}

	Lua.Close()

	return result
}
