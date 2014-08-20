package models

import (
	"database/sql"
	"fmt"
	"libs/lua"
	"strconv"
	"time"
)

var Role RoleModel

type RoleModel struct {
	MaxActionValue int
	ActionWaitTime int
}

func init() {
	Role.MaxActionValue = 5
	Role.ActionWaitTime = 1800
	var RoleData RoleData
	DB().AddTableWithName(RoleData, "role").SetKeys(false, "Uid")
}

func (this RoleModel) Role(uid int64) *RoleData {

	RoleData := new(RoleData)
	if err := DB().SelectOne(RoleData, "SELECT * FROM role WHERE uid = ?", uid); err != nil {
		if err == sql.ErrNoRows {
			return nil
			//			RoleData.Uid = uid
			//			RoleData.Coin = 10000
			//			RoleData.Diamond = 10000
			//			RoleData.ActionTime = time.Now().Unix()
			//			RoleData.UnixTime = RoleData.ActionTime
			//			if err = DB().Insert(&RoleData); err != nil {
			//				DBError(err)
			//			} else {
			//				return NewRoleData(uid)
			//			}
		} else {
			DBError(err)
		}
	}

	return RoleData
}

func (this RoleModel) Insert(uid int64) *RoleData {

	Lua, _ := lua.NewLua("conf/role.lua")

	RoleData := new(RoleData)
	RoleData.Uid = uid
	RoleData.Coin = Lua.GetInt("new_coin")
	RoleData.Diamond = Lua.GetInt("new_diamond")
	RoleData.GeneralBaseId = Lua.GetInt("new_leader")
	RoleData.UnixTime = time.Now().Unix()

	Lua.Close()

	if err := DB().Insert(RoleData); err != nil {
		DBError(err)
	}

	return RoleData
}

func (this RoleModel) NumberByRoleName(name string) int64 {
	n, err := DB().SelectInt("SELECT COUNT(*) FROM role WHERE role_name = ?", name)
	if err != nil {
		DBError(err)
	}
	return n
}

func (this RoleModel) FriendList(uids []int64) []*RoleData {

	var s string
	for _, uid := range uids {
		s += strconv.Itoa(int(uid)) + ","
	}

	var result []*RoleData
	_, err := DB().Select(&result, "SELECT * FROM `role` WHERE `uid` IN("+s[:len(s)-1]+") ORDER BY `role_kill_num` DESC")
	if err != nil {
		DBError(err)
	}

	return result
}

// role
type RoleData struct {
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
	SignNum         int    `db:"role_sign_num"`
	BuyActionDate   string `db:"role_buy_action_date"` // 购买体力日期
	BuyActionNum    int    `db:"role_buy_action_num"`  // 购买体力次数
	GeneralBaseId   int    `db:"general_base_id"`
	UnixTime        int64  `db:"role_time"`
}

func (this *RoleData) SetName(name string) error {
	this.Name = name
	this.UnixTime = time.Now().Unix()
	_, err := DB().Update(this)
	return err
}

func (this *RoleData) ActionValue() int {

	now := time.Now()
	n := (int(now.Unix() - this.ActionTime)) / Role.ActionWaitTime
	if n > Role.MaxActionValue {
		n = Role.MaxActionValue
	}

	return int(n) + this.OtherAction
}

func (this *RoleData) SetActionValue(n int) error {

	RoleData := *this

	if n > Role.MaxActionValue {
		this.OtherAction = n - Role.MaxActionValue
		n = Role.MaxActionValue
	} else {
		this.OtherAction = 0
	}

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % Role.ActionWaitTime
	this.ActionTime = nowUnix - int64(remainder) - int64(Role.ActionWaitTime*n)
	this.UnixTime = nowUnix

	_, err := DB().Update(this)
	if err != nil {
		this = &RoleData
	}
	return err
}

func (this *RoleData) BuyActionValue(diamond int, n int) error {

	RoleData := *this

	if n > Role.MaxActionValue {
		this.OtherAction = n - Role.MaxActionValue
		n = Role.MaxActionValue
	}

	oldDiamond := this.Diamond
	oldAction := this.ActionValue()

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % Role.ActionWaitTime
	this.ActionTime = nowUnix - int64(remainder) - int64(Role.ActionWaitTime*n)
	this.UnixTime = nowUnix
	this.Diamond -= diamond

	_, err := DB().Update(this)
	if err == nil {
		InsertSubDiamondFinanceLog(this.Uid, FINANCE_BUY_ACTION, oldDiamond, this.Diamond, fmt.Sprintf("%d -> %d", oldAction, n))
	} else {
		this = &RoleData
	}
	return err
}

func (this *RoleData) ActionRecoverTime() int {

	nowUnix := time.Now().Unix()
	remainder := int(nowUnix-this.ActionTime) % Role.ActionWaitTime

	return Role.ActionWaitTime - remainder
}

func (this *RoleData) AddCoin(n int, finance_type FinanceType, desc string) error {

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

func (this *RoleData) SubCoin(n int, finance_type FinanceType, desc string) error {

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

func (this *RoleData) AddDiamond(n int, finance_type FinanceType, desc string) error {

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

func (this *RoleData) SubDiamond(n int, finance_type FinanceType, desc string) error {

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

func (this *RoleData) DiamondIntoCoin(diamond int, coin int, desc string) error {

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

func (this *RoleData) SetGeneralConfigId(baseId int) error {

	temp := this.GeneralBaseId
	this.GeneralBaseId = baseId
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	if err != nil {
		this.GeneralBaseId = temp
	}
	return err
}

//func (this *RoleData) AddKillNum(num int, coin int, diamond int, desc string) error {

//	killNum, oldCoin, oldDiamond := this.KillNum, this.Coin, this.Diamond

//	this.KillNum += num
//	this.Coin += coin
//	this.Diamond += diamond
//	this.UnixTime = time.Now().Unix()

//	_, err := DB().Update(this)
//	if err != nil {
//		this.KillNum, this.Coin, this.Diamond = killNum, oldCoin, oldDiamond
//	} else {
//		if coin > 0 {
//			InsertAddCoinFinanceLog(this.Uid, FINANCE_DUPLICATE_GET, oldCoin, this.Coin, desc)
//		}
//		if diamond > 0 {
//			InsertAddDiamondFinanceLog(this.Uid, FINANCE_DUPLICATE_GET, oldDiamond, this.Diamond, desc)
//		}
//	}
//	return err
//}

//func (this *RoleData) SetUnlimitedNum(num int) error {

//	temp1, temp2 := this.UnlimitedNum, this.UnlimitedMaxNum

//	this.UnixTime = time.Now().Unix()
//	this.UnlimitedNum = num
//	if num > this.UnlimitedMaxNum {
//		this.UnlimitedMaxNum = num
//	}

//	_, err := DB().Update(this)
//	if err != nil {
//		this.UnlimitedNum, this.UnlimitedMaxNum = temp1, temp2
//	}
//	return err
//}

func (this *RoleData) Sign() error {

	temp1 := this.SignDate
	temp2 := this.SignNum
	now := time.Now()
	this.UnixTime = now.Unix()
	if this.SignDate == now.Format("20060102") {
		return nil
	}

	if this.SignDate == now.AddDate(0, 0, -1).Format("20060102") {
		this.SignNum++
	} else {
		this.SignNum = 1
	}

	this.SignDate = now.Format("20060102")

	_, err := DB().Update(this)
	if err != nil {
		this.SignDate = temp1
		this.SignNum = temp2
	}

	return err
}

func (this *RoleData) UpdateDate() error {

	now := time.Now()
	if this.BuyActionDate == now.Format("20060102") {
		return nil
	}

	temp1 := this.BuyActionDate
	temp2 := this.BuyActionNum

	this.UnixTime = now.Unix()
	this.BuyActionDate = now.Format("20060102")
	this.BuyActionNum = 0

	_, err := DB().Update(this)
	if err != nil {
		this.BuyActionDate = temp1
		this.BuyActionNum = temp2
	}

	return err
}

func (this *RoleData) Set() error {
	this.UnixTime = time.Now().Unix()
	_, err := DB().Update(this)
	return err
}
