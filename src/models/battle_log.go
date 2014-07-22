package models

import (
	"time"
)

type BattleType int

const (
	_ BattleType = iota
	BATTLE_TYPE_A
	BATTLE_TYPE_B
	BATTLE_TYPE_C
)

func init() {
	DB().AddTableWithName(BattleLogModel{}, "battle_log").SetKeys(true, "Id")
}

// battle_log
type BattleLogModel struct {
	Id        int        `db:"battle_id"`
	Uid       int64      `db:"uid"`
	Chapter   int        `db:"duplicate_chapter"`
	Section   int        `db:"duplicate_section"`
	GeneralId int        `db:"general_id"`
	Type      BattleType `db:"battle_type"`
	KillNum   int        `db:"battle_kill_num"`
	Result    bool       `db:"battle_result"`
	Time      int64      `db:"battle_time"`
}

func InsertBattleLog(battleLog *BattleLogModel) error {
	battleLog.Time = time.Now().Unix()
	return DB().Insert(battleLog)
}

func NewBattleLogModel(battle_id int) *BattleLogModel {
	var battleLog BattleLogModel
	err := DB().SelectOne(&battleLog, "SELECT * FROM battle_log WHERE battle_id = ?", battle_id)
	if err != nil {
		return nil
	} else {
		return &battleLog
	}
}

func (this *BattleLogModel) Win(killNum int) error {

	this.Result = true
	this.KillNum = killNum
	this.Time = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}
