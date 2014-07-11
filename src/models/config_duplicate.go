package models

var duplicate_config_map map[int]*ConfigDuplicate

func init() {
	var temp []*ConfigDuplicate
	if _, err := DB().Select(&temp, "SELECT * FROM config_duplicate"); err != nil {
		panic(err)
	}
	duplicate_config_map = make(map[int]*ConfigDuplicate)
	for _, duplicate := range temp {
		duplicate_config_map[duplicate.ConfigId] = duplicate
	}
}

// config_duplicate
type ConfigDuplicate struct {
	ConfigId    int    `db:"duplicate_config_id"`
	ChapterName string `db:"duplicate_chapter_name"`
	SectionName string `db:"duplicate_section_name"`
	Chapter     int    `db:"duplicate_chapter"`
	Section     int    `db:"duplicate_section"`
	NPC         int    `db:"duplicate_npc"`
	Items       int    `db:"duplicate_items"`
	Generals    int    `db:"duplicate_generals"`
}

func ConfigDuplicateMap() map[int]*ConfigDuplicate {
	return duplicate_config_map
}
