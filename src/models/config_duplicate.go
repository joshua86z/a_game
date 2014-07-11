package models

var duplicate_config_list []*ConfigDuplicate

func init() {
	if _, err := DB().Select(&duplicate_config_list, "SELECT * FROM config_duplicate ORDER BY duplicate_config_id ASC "); err != nil {
		panic(err)
	}
}

// config_duplicate
type ConfigDuplicate struct {
	ConfigId    int    `db:"duplicate_config_id"`
	ChapterName string `db:"duplicate_chapter_name"`
	SectionName string `db:"duplicate_section_name"`
	Chapter     int    `db:"duplicate_chapter"`
	Section     int    `db:"duplicate_section"`
	NPC         string `db:"duplicate_npc"`
	Items       string `db:"duplicate_items"`
	Generals    string `db:"duplicate_generals"`
	ChapterDesc string `db:"duplicate_chapter_desc"`
	SectionDesc string `db:"duplicate_section_desc"`
}

func ConfigDuplicateList() []*ConfigDuplicate {
	return duplicate_config_list
}
