package models

var item_config_map map[int]*ConfigItem

func init() {
	var temp []*ConfigItem
	if _, err := DB().Select(&temp, "SELECT * FROM config_item"); err != nil {
		panic(err)
	}
	item_config_map = make(map[int]*ConfigItem)
	for _, item := range temp {
		item_config_map[item.ConfigId] = item
	}
}

// config_item
type ConfigItem struct {
	ConfigId   int    `db:"item_config_id"`
	Name       string `db:"item_name"`
	Desc       string `db:"item_desc"`
}

func ConfigItemMap() map[int]*ConfigItem {
	return item_config_map
}
