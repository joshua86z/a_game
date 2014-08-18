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
	DB().AddTableWithName(BattleLogModel{}, "battle_logs").SetKeys(true, "Id")
}

// battle_logs
type BattleLogModel struct {
	Id        int        `db:"battle_id"`
	Uid       int64      `db:"uid"`
	Chapter   int        `db:"duplicate_chapter"`
	Section   int        `db:"duplicate_section"`
	GeneralId int        `db:"general_id"`
	Type      BattleType `db:"battle_type"`
	KillNum   int        `db:"battle_kill_num"`
	Result    int        `db:"battle_result"`
	Time      int64      `db:"battle_time"`
}

func InsertBattleLog(battleLog *BattleLogModel) error {
	battleLog.Time = time.Now().Unix()
	return DB().Insert(battleLog)
}

func NewBattleLogModel(battle_id int) *BattleLogModel {
	battleLog := new(BattleLogModel)
	err := DB().SelectOne(battleLog, "SELECT * FROM battle_logs WHERE battle_id = ?", battle_id)
	if err != nil {
		return nil
	} else {
		return battleLog
	}
}

func (this *BattleLogModel) SetResult(isWin bool, killNum int) error {

	this.KillNum = killNum
	this.Time = time.Now().Unix()
	if isWin {
		this.Result = 2
	} else {
		this.Result = 1
	}
	_, err := DB().Update(this)
	return err
}

func LastBattleLog(uid int64) *BattleLogModel {
	battleLog := new(BattleLogModel)
	DB().SelectOne(battleLog, "SELECT * FROM battle_logs WHERE uid = ? ORDER BY battle_id DESC LIMIT 1", uid)
	return battleLog
}
