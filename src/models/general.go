package models

import (
	"database/sql"
	"time"
)

// role_generals
type GeneralData struct {
	Id       int    `db:"general_id"`
	Uid      int64  `db:"uid"`
	ConfigId int    `db:"general_config_id"`
	Name     string `db:"general_name"`
	Level    int    `db:"general_level"`
	Atk      int    `db:"general_atk"`
	Def      int    `db:"general_def"`
	Hp       int    `db:"general_hp"`
	Speed    int    `db:"general_speed"`
	Dex      int    `db:"general_dex"`
	Range    int    `db:"general_range"`
	KillNum  int    `db:"general_kill_num"`
	UnixTime int64  `db:"general_time"`
}

func init() {
	DB().AddTableWithName(GeneralData{}, "role_generals").SetKeys(true, "Id")
}

var General GeneralModel

type GeneralModel struct {
}

func (this *GeneralModel) List(uid int64) []*GeneralData {
	var result []*GeneralData
	_, err := DB().Select(&result, "SELECT * FROM role_generals WHERE uid = ? ", uid)
	if err != nil {
		DBError(err)
	}
	return result
}

func (this *GeneralModel) General(uid int64, configId int) *GeneralData {
	GeneralData := new(GeneralData)
	err := DB().SelectOne(GeneralData, "SELECT * FROM role_generals WHERE uid = ? AND general_config_id = ? LIMIT 1", uid, configId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		DBError(err)
	}
	return GeneralData
}

func (this *GeneralModel) Insert(uid int64, config *ConfigGeneral) *GeneralData {

	general := new(GeneralData)
	general.Uid = uid
	general.ConfigId = config.ConfigId
	general.Name = config.Name
	general.Level = 0
	general.Atk = config.Atk
	general.Def = config.Def
	general.Hp = config.Hp
	general.Speed = config.Speed
	general.Dex = config.Dex
	general.Range = config.Range
	general.UnixTime = time.Now().Unix()

	if err := DB().Insert(general); err != nil {
		return nil
	}
	return general
}

func (this *GeneralData) LevelUp(config *ConfigGeneral) error {

	this.Level += 1
	this.Atk += config.AtkGroup
	this.Def += config.DefGroup
	this.Hp += config.HpGroup
	this.Speed += config.SpeedGroup
	this.Dex += config.DexGroup
	this.Range += config.RangeGroup
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}
