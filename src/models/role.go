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
	Uid        int64 `db:"uid"`
	Coin       int   `db:"role_coin"`
	Diamond    int   `db:"role_diamond"`
	ActionTime int64 `db:"role_action_time"` // 上次体力更新时间
	UnixTime   int64 `db:"role_time"`
}

func (this RoleModel) tableName() string {
	return "role"
}

func NewRoleModel(uid int64) *RoleModel {

	var RoleModel RoleModel
	if err := DB().SelectOne(&RoleModel, "SELECT * FROM "+RoleModel.tableName()+" WHERE uid = ?", uid); err != nil {
		if err == sql.ErrNoRows {
			RoleModel.Uid = uid
			RoleModel.Coin = 10000
			RoleModel.Diamond = 10000
			RoleModel.ActionTime = time.Now().Unix()
			RoleModel.UnixTime = RoleModel.ActionTime
			if err = DB().Insert(&RoleModel); err != nil {
				DBError(err)
			} else {
				return NewRoleModel(uid)
			}
		} else {
			DBError(err)
		}
	}

	return &RoleModel
}

var (
	MaxActionValue int64 = 20
	ActionWaitTime int64 = 1800
)

func (this *RoleModel) ActionValue() int {

	now := time.Now()
	n := (now.Unix() - this.ActionTime) / ActionWaitTime
	if n > MaxActionValue {
		n = MaxActionValue
	}

	return int(n)
}

func (this *RoleModel) SetActionValue(n int) error {

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

func (this *RoleModel) AddCoin(n int) error {

	this.Coin += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func (this *RoleModel) SubCoin(n int) error {

	this.Coin -= n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func (this *RoleModel) AddDiamond(n int) error {

	this.Diamond += n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func (this *RoleModel) SubDiamond(n int) error {

	this.Diamond -= n
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}
