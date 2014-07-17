package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

type Config_Pay_Center struct {
	Id      int    `db:"pay_config_id"`
	Name    string `db:"pay_name"`
	Rmb     int    `db:"pay_rmb"`
	Diamond int    `db:"pay_diamond"`
	Desc    string `db:"pay_desc"`
}

var config_pay_center []*Config_Pay_Center

func init() {
	if _, err := DB().Select(&config_pay_center, "SELECT * FROM `config_pay_center` ORDER BY `pay_config_id` ASC "); err != nil {
		panic(err)
	}
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
			Rmb:     Atoi(array[1]),
			Diamond: Atoi(array[2]),
			Desc:    array[3]})
	}

	Lua.Close()
	return result
}
