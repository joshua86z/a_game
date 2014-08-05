package models

import (
	"libs/log"
	"time"
)

func init() {
	DB().AddTableWithName(finance{}, "role_finance_logs").SetKeys(true, "Id")
}

type FinanceType int

const (
	_ FinanceType = iota
	FINANCE_ADMIN
	FINANCE_BUY_DIAMOND
	FINANCE_BUY_ACTION
	FINANCE_BUY_GENERAL
	FINANCE_BUY_COIN
	FINANCE_ITEM_LEVELUP
	FINANCE_GENERAL_LEVELUP
	FINANCE_MAIL_GET
	FINANCE_DUPLICATE_USE
	FINANCE_DUPLICATE_GET
	FINANCE_SIGN_GET
)

type finance struct {
	Id          int         `db:"rfl_id"`
	Uid         int64       `db:"uid"`
	OldMoney    int         `db:"rfl_old_money"`
	NewMoney    int         `db:"rfl_new_money"`
	IsAdd       bool        `db:"rfl_type"`
	MoneyType   int         `db:"rfl_money_type"`
	Desc        string      `db:"rfl_desc"`
	UnixTime    int64       `db:"rfl_time"`
	FinanceType FinanceType `db:"rfl_static_type"`
}

var (
	financeStrMap map[FinanceType]string
	financeChan   chan *finance
)

func init() {
	financeChan = make(chan *finance, 1000)
	go checkfinanceChan()

	financeStrMap = make(map[FinanceType]string)
	financeStrMap[FINANCE_ADMIN] = "内部添加"
	financeStrMap[FINANCE_BUY_DIAMOND] = "充值"
	financeStrMap[FINANCE_BUY_ACTION] = "买体力"
	financeStrMap[FINANCE_BUY_GENERAL] = "买英雄"
	financeStrMap[FINANCE_BUY_COIN] = "兑换金币"
	financeStrMap[FINANCE_ITEM_LEVELUP] = "道具升级"
	financeStrMap[FINANCE_GENERAL_LEVELUP] = "英雄升级"
	financeStrMap[FINANCE_MAIL_GET] = "邮件领取"
	financeStrMap[FINANCE_DUPLICATE_USE] = "临时道具"
	financeStrMap[FINANCE_DUPLICATE_GET] = "副本获得"
	financeStrMap[FINANCE_SIGN_GET] = "签到"
}

func checkfinanceChan() {

	defer func() {
		if err := recover(); err != nil {
			log.Critical("financeChan panic : %v", err)
			checkfinanceChan()
		}
	}()

	for finance := range financeChan {
		DB().Insert(finance)
		//sql := "INSERT INTO g_role_finance_log SET roles_unique = ? , rfl_old_money = ? , rfl_new_money = ? , rfl_type = ? , rfl_mtype = ? , rfl_desc = ? , rfl_time = UNIX_TIMESTAMP() , rfl_static_type = ? "
		//DB().Exec(sql, finance.uid, finance.oldMoney, finance.newMoney, isAdd, finance.moneyType, finance.desc, finance.financeType)
	}
}

//rfl_type  1 增加 0 减少
//rfl_mtype 1 coin 2 diamond
func InsertAddDiamondFinanceLog(uid int64, finance_type FinanceType, old_diamond int, new_diamond int, desc string) {

	desc = financeLogDesc(finance_type, desc)

	financeChan <- &finance{

		Uid:         uid,
		OldMoney:    old_diamond,
		NewMoney:    new_diamond,
		IsAdd:       true,
		MoneyType:   2,
		Desc:        desc,
		UnixTime:    time.Now().Unix(),
		FinanceType: finance_type,
	}
}

func InsertSubDiamondFinanceLog(uid int64, finance_type FinanceType, old_diamond int, new_diamond int, desc string) {

	desc = financeLogDesc(finance_type, desc)

	financeChan <- &finance{
		Uid:         uid,
		OldMoney:    old_diamond,
		NewMoney:    new_diamond,
		IsAdd:       false,
		MoneyType:   2,
		Desc:        desc,
		UnixTime:    time.Now().Unix(),
		FinanceType: finance_type,
	}
}

func InsertAddCoinFinanceLog(uid int64, finance_type FinanceType, old_coin int, new_coin int, desc string) {

	desc = financeLogDesc(finance_type, desc)

	financeChan <- &finance{
		Uid:         uid,
		OldMoney:    old_coin,
		NewMoney:    new_coin,
		IsAdd:       true,
		MoneyType:   1,
		Desc:        desc,
		UnixTime:    time.Now().Unix(),
		FinanceType: finance_type,
	}
}

func InsertSubCoinFinanceLog(uid int64, finance_type FinanceType, old_coin int, new_coin int, desc string) {

	desc = financeLogDesc(finance_type, desc)

	financeChan <- &finance{
		Uid:         uid,
		OldMoney:    old_coin,
		NewMoney:    new_coin,
		IsAdd:       false,
		MoneyType:   1,
		Desc:        desc,
		UnixTime:    time.Now().Unix(),
		FinanceType: finance_type,
	}
}

func financeLogDesc(finance_type FinanceType, desc string) string {
	if desc != "" {
		return financeStrMap[finance_type] + " " + desc
	} else {
		return financeStrMap[finance_type]
	}
}
