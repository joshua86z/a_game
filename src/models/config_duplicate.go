package models

import (
	"fmt"
	"libs/lua"
	"strings"
)

func init() {
	ConfigDuplicateList()
}

// config_duplicate
type ConfigDuplicate struct {
	ConfigId    int    `db:"duplicate_config_id"`
	Chapter     int    `db:"duplicate_chapter"`
	Section     int    `db:"duplicate_section"`
	ChapterName string `db:"duplicate_chapter_name"`
	SectionName string `db:"duplicate_section_name"`
	Value       string `db:"duplicate_value"`
	GenId       int
}

func ConfigDuplicateList() []*ConfigDuplicate {

	var result []*ConfigDuplicate

	Lua, _ := lua.NewLua("conf/duplicate.lua")
	var i int
	for {
		i++
		duplicateStr := Lua.GetString(fmt.Sprintf("duplicate_%d", i))
		if duplicateStr == "" {
			break
		}
		array := strings.Split(duplicateStr, "\\,")
		result = append(result, &ConfigDuplicate{
			ConfigId:    i,
			Chapter:     Atoi(array[0]),
			Section:     Atoi(array[1]),
			ChapterName: array[2],
			SectionName: array[3],
			Value:       array[4],
			GenId:       Atoi(array[5])})
	}

	Lua.Close()

	return result
}
