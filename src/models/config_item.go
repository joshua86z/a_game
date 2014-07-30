package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

func init() {
	ConfigItemList()
}

type ConfigItem struct {
	ConfigId    int
	Name        string
	Desc        string
	Value       int
	Group       int
	Probability int
}

func ConfigItemList() []*ConfigItem {

	var result []*ConfigItem

	Lua, _ := lua.NewLua("conf/item.lua")
	var i int
	for {
		i++
		itemStr := Lua.GetString(fmt.Sprintf("item_%d", i))
		if itemStr == "" {
			break
		}
		array := strings.Split(itemStr, "\\,")
		result = append(result, &ConfigItem{
			ConfigId:    Atoi(array[0]),
			Name:        array[1],
			Desc:        array[2],
			Value:       Atoi(array[3]),
			Group:       Atoi(array[4]),
			Probability: Atoi(array[5])})
	}

	Lua.Close()

	return result
}
