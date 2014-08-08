package models

import (
	"fmt"
	"libs/lua"
	"strconv"
	"strings"
)

func init() {
	ConfigGeneralMap()
}

// config_general
type ConfigGeneral struct {
	ConfigId    int    `db:"general_config_id"`
	Name        string `db:"general_name"`
	Type        int    `db:"general_type"`
	Atk         int    `db:"general_atk"`
	Def         int    `db:"general_def"`
	Hp          int    `db:"general_hp"`
	Speed       int    `db:"general_speed"`
	Dex         int    `db:"general_dex"`
	Range       int    `db:"general_range"`
	AtkRange    int    `db:"general_atk_range"`
	AtkGroup    int    `db:"general_atk_group"`
	DefGroup    int    `db:"general_def_group"`
	HpGroup     int    `db:"general_hp_group"`
	SpeedGroup  int    `db:"general_speed_group"`
	DexGroup    int    `db:"general_dex_group"`
	RangeGroup  int    `db:"general_range_group"`
	BuyDiamond  int    `db:"general_buy_diamond"`
	SkillAtk    int    `db:"general_skill_atk"`
	Desc        string `db:"general_desc"`
	LevelUpCoin []int  `db:"-"`
}

func ConfigGeneralMap() map[int]*ConfigGeneral {

	result := make(map[int]*ConfigGeneral)

	Lua, _ := lua.NewLua("conf/general.lua")

	var coinList []int //升级需要的金币
	levelUpCoin := Lua.GetString("level_up_coin")
	for _, val := range strings.Split(levelUpCoin, ",") {
		coinList = append(coinList, Atoi(val))
	}

	var i int
	for {
		i++
		itemStr := Lua.GetString(fmt.Sprintf("general_%d", i))
		if itemStr == "" {
			break
		}
		array := strings.Split(itemStr, "\\,")
		result[Atoi(array[0])] = genByStr(array)
		result[Atoi(array[0])].LevelUpCoin = coinList
	}

	Lua.Close()
	return result
}

func genByStr(array []string) *ConfigGeneral {

	return &ConfigGeneral{
		ConfigId:   Atoi(array[0]),
		Name:       array[1],
		Type:       Atoi(array[2]),
		Atk:        Atoi(array[3]),
		Def:        Atoi(array[4]),
		Hp:         Atoi(array[5]),
		Speed:      Atoi(array[6]),
		Dex:        Atoi(array[7]),
		Range:      Atoi(array[8]),
		AtkRange:   Atoi(array[9]),
		AtkGroup:   Atoi(array[10]),
		DefGroup:   Atoi(array[11]),
		HpGroup:    Atoi(array[12]),
		SpeedGroup: Atoi(array[13]),
		DexGroup:   Atoi(array[14]),
		RangeGroup: Atoi(array[15]),
		BuyDiamond: Atoi(array[16]),
		Desc:       array[17],
		SkillAtk:   Atoi(array[18])}
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
