package models

import (
	"fmt"
	"libs/lua"
	"strconv"
	"strings"
)

//var duplicate_config_list []*ConfigDuplicate
//
//func init() {
//	if _, err := DB().Select(&duplicate_config_list, "SELECT * FROM config_duplicate ORDER BY duplicate_config_id ASC "); err != nil {
//		panic(err)
//	}
//}

// config_duplicate
type ConfigDuplicate struct {
	ConfigId    int    `db:"duplicate_config_id"`
	Chapter     int    `db:"duplicate_chapter"`
	Section     int    `db:"duplicate_section"`
	ChapterName string `db:"duplicate_chapter_name"`
	SectionName string `db:"duplicate_section_name"`
	NPC         string `db:"duplicate_npc"`
	Items       string `db:"duplicate_items"`
	Generals    string `db:"duplicate_generals"`
	ChapterDesc string `db:"duplicate_chapter_desc"`
	SectionDesc string `db:"duplicate_section_desc"`
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
		chapter, _ := strconv.Atoi(array[0])
		section, _ := strconv.Atoi(array[1])
		result = append(result, &ConfigDuplicate{
			ConfigId:    i,
			Chapter:     chapter,
			Section:     section,
			ChapterName: array[2],
			SectionName: array[3],
			NPC:         array[4],
			Items:       array[5],
			Generals:    array[6],
			ChapterDesc: array[7],
			SectionDesc: array[8],
		})
	}

	Lua.Close()

	return result
}
