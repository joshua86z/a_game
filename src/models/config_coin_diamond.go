package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

type ConfigCoinDiamond struct {
	Index   int
	Name    string
	Coin    int
	Diamond int
	Desc    string
}

func init() {
}

func ConfigCoinDiamondList() []*ConfigCoinDiamond {

	var result []*ConfigCoinDiamond

	Lua, _ := lua.NewLua("conf/coin_diamond.lua")
	var i int
	for {
		i++
		itemStr := Lua.GetString(fmt.Sprintf("coin_diamond_%d", i))
		if itemStr == "" {
			break
		}

		array := strings.Split(itemStr, "\\,")

		result = append(result, &ConfigCoinDiamond{
			Index:   i,
			Name:    array[0],
			Coin:    Atoi(array[1]),
			Diamond: Atoi(array[2]),
			Desc:    array[3]})
	}

	Lua.Close()
	return result
}
