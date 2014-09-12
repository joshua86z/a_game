package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

type Config_Pay_Center struct {
	Id      int    `db:"pay_config_id"`
	Name    string `db:"pay_name"`
	Money   int    `db:"pay_money"`
	Diamond int    `db:"pay_diamond"`
	Desc    string `db:"pay_desc"`
}

// 充值商店
func ConfigPayCenterList() []*Config_Pay_Center {

	var result []*Config_Pay_Center

	Lua, _ := lua.NewLua("conf/pay_center.lua")
	var i int
	for {
		i++
		itemStr := Lua.GetString(fmt.Sprintf("pay_center_%d", i))
		if itemStr == "" {
			break
		}

		array := strings.Split(itemStr, "\\,")

		result = append(result, &Config_Pay_Center{
			Id:      i,
			Name:    array[0],
			Money:   Atoi(array[1]),
			Diamond: Atoi(array[2]),
			Desc:    array[3]})
	}

	Lua.Close()
	return result
}
