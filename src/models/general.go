package models

import (
	"fmt"
	"time"
)

// g_role_generals
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
	Range    int    `db:"general_range"`
	UnixTime int64  `db:"general_time"`
}

func init() {
	DB().AddTableWithName(GeneralData{}, "role_generals").SetKeys(true, "Id")
}

type GeneralModel struct {
	Uid         int64
	GeneralList []*GeneralData
}

func GetGeneralModel(uid int64) *GeneralModel {

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

func (this *GeneralModel) GetGeneral(generalId int) *GeneralData {
	for _, general := range this.GeneralList {
		if general.Id == generalId {
			return general
		}
	}
	panic(fmt.Sprintf("武将数据错误: 没找到这个武将 %d", generalId))
}

func (this *GeneralData) LevelUp() error {

	config := ConfigGeneralMap()[this.ConfigId]

	this.Level += 1
	this.Atk += config.AtkGroup
	this.Def += config.DefGroup
	this.Hp += config.HpGroup
	this.Speed += config.SpeedGroup
	this.Range += config.RangeGroup
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func InsertGeneral(uid int64, configId int) *GeneralData {

	c := ConfigGeneralMap()[configId]

	general := &GeneralData{}

	general.Uid = uid
	general.ConfigId = configId
	general.Name = c.Name
	general.Level = 1
	general.Atk += c.Atk
	general.Def += c.Def
	general.Hp += c.Hp
	general.Speed += c.Speed
	general.Range += c.Range
	general.UnixTime = time.Now().Unix()

	if err := DB().Insert(general); err != nil {
		DBError(err)
	}

	return general
}

func DeleteGeneral(generalId int) error {

	_, err := DB().Exec("DELETE FROM role_generals WHERE general_id = ? ", generalId)
	return err
}
