package models

import (
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

type GeneralModel struct {
	Uid         int64
	GeneralList []*GeneralData
}

func NewGeneralModel(uid int64) *GeneralModel {

	var General GeneralModel

	var temp []*GeneralData
	_, err := DB().Select(&temp, "SELECT * FROM role_generals WHERE uid = ?  ", uid)
	if err != nil {
		DBError(err)
	}

	General.Uid = uid
	General.GeneralList = temp

	return &General
}

func (this *GeneralModel) List() []*GeneralData {
	return this.GeneralList
}

func (this *GeneralModel) General(configId int) *GeneralData {
	for _, general := range this.GeneralList {
		if general.ConfigId == configId {
			return general
		}
	}
	return nil
}

func (this *GeneralModel) Insert(config *ConfigGeneral) *GeneralData {

	general := &GeneralData{}
	general.Uid = this.Uid
	general.ConfigId = config.ConfigId
	general.Name = config.Name
	general.Level = 1
	general.Atk = config.Atk
	general.Def = config.Def
	general.Hp = config.Hp
	general.Speed = config.Speed
	general.Dex = config.Dex
	general.Range = config.Range
	general.UnixTime = time.Now().Unix()

	if err := DB().Insert(general); err != nil {
		return nil
	} else {
		this.GeneralList = append(this.GeneralList, general)
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
