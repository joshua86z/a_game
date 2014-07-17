package models

import (
	"fmt"
	"libs/lua"
	"strconv"
	"strings"
)

//var general_config_map map[int]*ConfigGeneral
//
//func init() {
//	var temp []*ConfigGeneral
//	if _, err := DB().Select(&temp, "SELECT * FROM config_general"); err != nil {
//		panic(err)
//	}
//	general_config_map = make(map[int]*ConfigGeneral)
//	for _, general := range temp {
//		general_config_map[general.ConfigId] = general
//	}
//}

// config_general
type ConfigGeneral struct {
	ConfigId   int    `db:"general_config_id"`
	Name       string `db:"general_name"`
	Type       int    `db:"general_type"`
	Atk        int    `db:"general_atk"`
	Def        int    `db:"general_def"`
	Hp         int    `db:"general_hp"`
	Speed      int    `db:"general_speed"`
	Dex        int    `db:"general_dex"`
	Range      int    `db:"general_range"`
	AtkRange   int    `db:"general_atk_range"`
	AtkGroup   int    `db:"general_atk_group"`
	DefGroup   int    `db:"general_def_group"`
	HpGroup    int    `db:"general_hp_group"`
	SpeedGroup int    `db:"general_speed_group"`
	DexGroup   int    `db:"general_dex_group"`
	RangeGroup int    `db:"general_range_group"`
	BuyDiamond int    `db:"general_buy_diamond"`
	Desc       string `db:"general_desc"`
}

func ConfigGeneralMap() map[int]*ConfigGeneral {

	result := make(map[int]*ConfigGeneral)

	Lua, _ := lua.NewLua("conf/general.lua")

	var configId int
	for {
		configId++
		itemStr := Lua.GetString(fmt.Sprintf("general_%d", configId))
		if itemStr == "" {
			break
		}
		array := strings.Split(itemStr, "\\,")
		result[configId] = genByStr(configId, array)
	}

	Lua.Close()
	return result
}

func ConfigGeneralById(configId int) *ConfigGeneral {

	Lua, _ := lua.NewLua("conf/general.lua")

	itemStr := Lua.GetString(fmt.Sprintf("general_%d", configId))
	array := strings.Split(itemStr, "\\,")

	Lua.Close()

	return genByStr(configId, array)
}

func genByStr(configId int, array []string) *ConfigGeneral {
	return &ConfigGeneral{
		ConfigId:   configId,
		Name:       array[0],
		Type:       Atoi(array[1]),
		Atk:        Atoi(array[2]),
		Def:        Atoi(array[3]),
		Hp:         Atoi(array[4]),
		Speed:      Atoi(array[5]),
		Dex:        Atoi(array[6]),
		Range:      Atoi(array[7]),
		AtkRange:   Atoi(array[8]),
		AtkGroup:   Atoi(array[9]),
		DefGroup:   Atoi(array[10]),
		HpGroup:    Atoi(array[11]),
		SpeedGroup: Atoi(array[12]),
		DexGroup:   Atoi(array[13]),
		RangeGroup: Atoi(array[14]),
		BuyDiamond: Atoi(array[15]),
		Desc:       array[16]}
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
