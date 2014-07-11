package models

var general_config_map map[int]*ConfigGeneral

func init() {
	var temp []*ConfigGeneral
	if _, err := DB().Select(&temp, "SELECT * FROM config_general"); err != nil {
		panic(err)
	}
	general_config_map = make(map[int]*ConfigGeneral)
	for _, general := range temp {
		general_config_map[general.ConfigId] = general
	}
}

// config_general
type ConfigGeneral struct {
	ConfigId   int    `db:"general_config_id"`
	Name       string `db:"general_name"`
	Type       int    `db:"general_type"`
	Atk        int    `db:"general_atk"`
	Def        int    `db:"general_def"`
	Hp         int    `db:"general_hp"`
	Speed      int    `db:"general_speed"`
	Range      int    `db:"general_range"`
	AtkGroup   int    `db:"general_atk_group"`
	DefGroup   int    `db:"general_def_group"`
	HpGroup    int    `db:"general_hp_group"`
	SpeedGroup int    `db:"general_speed_group"`
	RangeGroup int    `db:"general_range_group"`
	Desc       string `db:"general_desc"`
}

func ConfigGeneralMap() map[int]*ConfigGeneral {
	return general_config_map
}
