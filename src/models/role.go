package models

import (
	"database/sql"
	"fmt"
	"time"
)

func init() {
	var RoleModel RoleModel
	DB().AddTableWithName(RoleModel, RoleModel.tableName()).SetKeys(false, "Uid")
}

// role
type RoleModel struct {
	Uid             int64  `db:"uid"`
	Name            string `db:"role_name"`
	Coin            int    `db:"role_coin"`
	Diamond         int    `db:"role_diamond"`
	OtherAction     int    `db:"role_other_action"`
	ActionTime      int64  `db:"role_action_time"` // 上次体力更新时间
	UnlimitedMaxNum int    `db:"role_unlimited_max_num"`
	UnlimitedNum    int    `db:"role_unlimited_num"`
	KillNum         int    `db:"role_kill_num"`
	SignDate        string `db:"role_sign_date"`
	SignTimes       int    `db:"role_sign_times"`
	GeneralConfigId int    `db:"general_config_id"`
	UnixTime        int64  `db:"role_time"`
}

func (this RoleModel) tableName() string {
	return "role"
}

func NewRoleModel(uid int64) *RoleModel {

	RoleModel := &RoleModel{}
	if err := DB().SelectOne(RoleModel, "SELECT * FROM "+RoleModel.tableName()+" WHERE uid = ?", uid); err != nil {
		if err == sql.ErrNoRows {
			return nil
			//			RoleModel.Uid = uid
			//			RoleModel.Coin = 10000
			//			RoleModel.Diamond = 10000
			//			RoleModel.ActionTime = time.Now().Unix()
			//			RoleModel.UnixTime = RoleModel.ActionTime
			//			if err = DB().Insert(&RoleModel); err != nil {
			//				DBError(err)
			//			} else {
			//				return NewRoleModel(uid)
			//			}
		} else {
			DBError(err)
		}
	}

	return RoleModel
}

func InsertRole(RoleModel *RoleModel) error {
	RoleModel.UnixTime = time.Now().Unix()
	return DB().Insert(RoleModel)
}

func NumberByRoleName(name string) int64 {
	n, err := DB().SelectInt("SELECT COUNT(*) FROM role WHERE role_name = ?", name)
	if err != nil {
		DBError(err)
	}
	return n
}

const (
	MaxActionValue int = 5
	ActionWaitTime int = 1800
)

func (this *RoleModel) SetName(name string) error {
	this.Name = name
	this.UnixTime = time.Now().Unix()
	_, err := DB().Update(this)
	return err
}

func (this *RoleModel) ActionValue() int {

	now := time.Now()
	n := (int(now.Unix() - this.ActionTime)) / ActionWaitTime
	if n > MaxActionValue {
		n = MaxActionValue
	}

	return int(n) + this.OtherAction
}

func (this *RoleModel) SetActionValue(n int) error {

	RoleModel := *this

	if n > MaxActionValue {
		this.OtherAction = n - MaxActionValue
		n = MaxActionValue
	}

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % ActionWaitTime
	this.ActionTime = nowUnix - int64(remainder) - int64(ActionWaitTime*n)
	this.UnixTime = nowUnix

	_, err := DB().Update(this)
	if err != nil {
		this = &RoleModel
	}
	return err
}

func (this *RoleModel) BuyActionValue(diamond int, n int) error {

	RoleModel := *this

	if n > MaxActionValue {
		this.OtherAction = n - MaxActionValue
		n = MaxActionValue
	}

	oldDiamond := this.Diamond
	oldAction := this.ActionValue()

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % ActionWaitTime
	this.ActionTime = nowUnix - int64(remainder) - int64(ActionWaitTime*n)
	this.UnixTime = nowUnix
	this.Diamond -= diamond

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, FINANCE_BUY_ACTION, oldDiamond, this.Diamond, fmt.Sprintf("%d -> %d", oldAction, n))
	} else {
		this = &RoleModel
	}
	return err
}

func (this *RoleModel) ActionRecoverTime() int {

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % ActionWaitTime

	return ActionWaitTime - remainder
}

func (this *RoleModel) AddCoin(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Coin

	this.Coin += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertAddDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Coin, desc)
	} else {
		this.Coin -= n
	}

	return err
}

func (this *RoleModel) SubCoin(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Coin

	this.Coin -= n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Coin, desc)
	} else {
		this.Coin += n
	}

	return err
}

func (this *RoleModel) AddDiamond(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Diamond

	this.Diamond += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertAddDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Diamond, desc)
	} else {
		this.Diamond -= n
	}

	return err
}

func (this *RoleModel) SubDiamond(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Diamond

	this.Diamond -= n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Diamond, desc)
	} else {
		this.Diamond += n
	}

	return err
}

func (this *RoleModel) DiamondIntoCoin(diamond int, coin int, desc string) error {

	oldCoin, oldDiamond := this.Coin, this.Diamond

	this.Diamond -= diamond
	this.Coin += coin
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, FINANCE_BUY_COIN, oldDiamond, this.Diamond, desc)
		InsertAddCoinFinanceLog(this.Uid, FINANCE_BUY_COIN, oldCoin, this.Coin, desc)
	} else {
		this.Coin, this.Diamond = oldCoin, oldDiamond
	}

	return err
}

func (this *RoleModel) SetGeneralConfigId(configId int) error {

	temp := this.GeneralConfigId
	this.GeneralConfigId = configId
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err != nil {
		this.GeneralConfigId = temp
	}
	return err
}

func (this *RoleModel) AddKillNum(num int, coin int, diamond int, desc string) error {

	killNum, oldCoin, oldDiamond := this.KillNum, this.Coin, this.Diamond

	this.KillNum += num
	this.Coin += coin
	this.Diamond += diamond
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err != nil {
		this.KillNum, this.Coin, this.Diamond = killNum, oldCoin, oldDiamond
	} else {
		if coin > 0 {
			InsertAddCoinFinanceLog(this.Uid, FINANCE_DUPLICATE_GET, oldCoin, this.Coin, desc)
		}
		if diamond > 0 {
			InsertAddDiamondFinanceLog(this.Uid, FINANCE_DUPLICATE_GET, oldDiamond, this.Diamond, desc)
		}
	}
	return err
}

func (this *RoleModel) SetUnlimitedNum(num int) error {

	temp1, temp2 := this.UnlimitedNum, this.UnlimitedMaxNum

	this.UnixTime = time.Now().Unix()
	this.UnlimitedNum = num
	if num > this.UnlimitedMaxNum {
		this.UnlimitedMaxNum = num
	}

	_, err := DB().Update(this)
	if err != nil {
		this.UnlimitedNum, this.UnlimitedMaxNum = temp1, temp2
	}
	return err
}

func (this *RoleModel) Sign() error {

	temp1 := this.SignDate
	temp2 := this.SignTimes
	now := time.Now()
	this.UnixTime = now.Unix()
	if this.SignDate == now.Format("20060102") {
		return nil
	}

	if this.SignDate == now.AddDate(0, 0, -1).Format("20060102") {
		this.SignTimes++
	} else {
		this.SignTimes = 1
	}

	this.SignDate = now.Format("20060102")

	_, err := DB().Update(this)
	if err != nil {
		this.SignDate = temp1
		this.SignTimes = temp2
	}

	return err
}
