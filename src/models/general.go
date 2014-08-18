package models

import (
	"database/sql"
	"time"
)

// role_generals
type GeneralData struct {
	Id       int    `db:"general_id"`
	Uid      int64  `db:"uid"`
	BaseId   int    `db:"general_base_id"`
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

func (this *GeneralModel) General(uid int64, baseId int) *GeneralData {
	general := new(GeneralData)
	err := DB().SelectOne(general, "SELECT * FROM role_generals WHERE uid = ? AND general_base_id = ? LIMIT 1", uid, baseId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		DBError(err)
	}
	return general
}

func (this *GeneralModel) Insert(uid int64, base *Base_General) *GeneralData {

	general := new(GeneralData)
	general.Uid = uid
	general.BaseId = base.BaseId
	general.Name = base.Name
	general.Level = 0
	general.Atk = base.Atk
	general.Def = base.Def
	general.Hp = base.Hp
	general.Speed = base.Speed
	general.Dex = base.Dex
	general.Range = base.Range
	general.UnixTime = time.Now().Unix()

	if err := DB().Insert(general); err != nil {
		return nil
	}
	return general
}

func (this *GeneralData) LevelUp(base *Base_General) error {

	this.Level += 1
	this.Atk += base.AtkGroup
	this.Def += base.DefGroup
	this.Hp += base.HpGroup
	this.Speed += base.SpeedGroup
	this.Dex += base.DexGroup
	this.Range += base.RangeGroup
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func (this *GeneralData) AddKillNum(num int) error {

	this.KillNum += num
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}
