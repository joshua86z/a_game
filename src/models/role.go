package models

import (
	"database/sql"
	"time"
)

func init() {
	var RoleModel RoleModel
	DB().AddTableWithName(RoleModel, RoleModel.tableName()).SetKeys(false, "Uid")
}

// role
type RoleModel struct {
	Uid         int64  `db:"uid"`
	Name        string `db:"role_name"`
	Coin        int    `db:"role_coin"`
	Diamond     int    `db:"role_diamond"`
	OtherAction int    `db:"role_other_action"`
	ActionTime  int64  `db:"role_action_time"` // 上次体力更新时间
	UnixTime    int64  `db:"role_time"`
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
	return DB().Insert(RoleModel)
}

func NumberByRoleName(name string) int64 {
	n, err := DB().SelectInt("SELECT COUNT(*) FROM role WHERE role_name = ?", name)
	if err != nil {
		DBError(err)
	}
	return n
}

var (
	MaxActionValue int64 = 5
	ActionWaitTime int64 = 1800
)

func (this *RoleModel) ActionValue() int {

	now := time.Now()
	n := (now.Unix() - this.ActionTime) / ActionWaitTime
	if n > MaxActionValue {
		n = MaxActionValue
	}

	return int(n) + this.OtherAction
}

func (this *RoleModel) SetActionValue(n int) error {

	if n > int(MaxActionValue) {
		this.OtherAction = n - int(MaxActionValue)
		n = int(MaxActionValue)
	}

	nowUnix := time.Now().Unix()
	remainder := (nowUnix - this.ActionTime) % ActionWaitTime
	this.ActionTime = nowUnix - remainder - ActionWaitTime*int64(n)
	this.UnixTime = nowUnix

	_, err := DB().Update(this)
	return err
}

func (this *RoleModel) ActionRecoverTime() int {

	nowUnix := time.Now().Unix()
	remainder := (nowUnix - this.ActionTime) % ActionWaitTime

	return int(ActionWaitTime - remainder)
}

func (this *RoleModel) AddCoin(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Coin

	this.Coin += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Coin, desc)
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
	}

	return err
}

func (this *RoleModel) AddDiamond(n int, finance_type FinanceType, desc string) error {

	oldMoney := this.Diamond

	this.Diamond += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, finance_type, oldMoney, this.Diamond, desc)
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
	}

	return err
}
