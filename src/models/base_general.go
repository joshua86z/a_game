package models

import (
	"fmt"
	"libs/lua"
	"strconv"
	"strings"
)

func init() {
	BaseGeneralMap()
}

type Base_General struct {
	BaseId      int    `db:"general_base_id"`
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

func BaseGeneralMap() map[int]*Base_General {

	result := make(map[int]*Base_General)

	Lua, _ := lua.NewLua("conf/general.lua")

	var coinList []int //升级需要的金币
	levelUpCoin := Lua.GetString("level_up_coin")
	for _, val := range strings.Split(levelUpCoin, ",") {
		coinList = append(coinList, Atoi(val))
	}

	indexStr := Lua.GetString("index")
	indexArr := strings.Split(indexStr, ",")
	for _, index := range indexArr {
		baseId := Atoi(index)
		result[baseId] = BaseGeneral(baseId, Lua)
		result[baseId].LevelUpCoin = coinList
	}

	Lua.Close()
	return result
}

func BaseGeneral(baseId int, Lua *lua.Lua) *Base_General {

	var sign bool
	if Lua == nil {
		sign = true
		Lua, _ = lua.NewLua("conf/general.lua")
		defer Lua.Close()
	}

	Lua.L.DoString(fmt.Sprintf("Name,Type,Atk,Def,Hp,Speed,Dex,Range,AtkRange,AtkGroup,DefGroup,HpGroup,SpeedGroup,DexGroup,RangeGroup,BuyDiamond,Desc,Skillhurt = baseGeneral(%d)", baseId))

	name := Lua.GetString("Name")
	if name == "" {
		return nil
	}

	result := &Base_General{
		BaseId:     baseId,
		Name:       name,
		Type:       Lua.GetInt("Type"),
		Atk:        Lua.GetInt("Atk"),
		Def:        Lua.GetInt("Def"),
		Hp:         Lua.GetInt("Hp"),
		Speed:      Lua.GetInt("Speed"),
		Dex:        Lua.GetInt("Dex"),
		Range:      Lua.GetInt("Range"),
		AtkRange:   Lua.GetInt("AtkRange"),
		AtkGroup:   Lua.GetInt("AtkGroup"),
		DefGroup:   Lua.GetInt("DefGroup"),
		HpGroup:    Lua.GetInt("HpGroup"),
		SpeedGroup: Lua.GetInt("SpeedGroup"),
		DexGroup:   Lua.GetInt("DexGroup"),
		RangeGroup: Lua.GetInt("RangeGroup"),
		BuyDiamond: Lua.GetInt("BuyDiamond"),
		Desc:       Lua.GetString("Desc"),
		SkillAtk:   Lua.GetInt("SkillAtk")}

	if sign {
		levelUpCoin := Lua.GetString("level_up_coin")
		for _, coin := range strings.Split(levelUpCoin, ",") {
			result.LevelUpCoin = append(result.LevelUpCoin, Atoi(coin))
		}
	}

	return result
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
