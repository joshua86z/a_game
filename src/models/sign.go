package models

import (
	"database/sql"
//	"fmt"
	"time"
)

// role_sign
type SignModel struct {
	Uid    int64  `db:"uid"`
	Date   string `db:"sign_date"`
	Times  int    `db:"sign_times"`
	Reward bool   `db:"sign_reward"`
}

func init() {
	DB().AddTableWithName(SignModel{}, "role_sign").SetKeys(false, "Uid")
}

func NewSignModel(uid int64) *SignModel {

	SignModel := &SignModel{}
	if err := DB().SelectOne(SignModel, "SELECT * FROM role_sign WHERE uid = ? ", uid); err != nil {
		if err != sql.ErrNoRows {
			DBError(err)
		}
		SignModel.Uid = uid
		SignModel.Date = time.Now().Format("20060102")
		SignModel.Times = 1
		SignModel.Reward = false
		err = DB().Insert(SignModel)
		if err != nil {
			DBError(err)
		}
		return NewSignModel(uid)
	}

	return SignModel
}

//
func (this *SignModel) Sign() error {
	now := time.Now()
	if this.Date == now.Format("20060102") {
		return nil
	}
	if this.Date == now.AddDate(0, 0, -1).Format("20060102") {
		this.Date = now.Format("20060102")
		if this.Reward {
			this.Times = 1
		} else {
			this.Times += 1
		}
	}

	_, err := DB().Update(this)
	return err
}

func (this *SignModel) GetReward() error {
	this.Reward = true
	_, err := DB().Update(this)
	return err
}

//
//func InsertSign(uid int64) *SignData {
//
//	sign := &SignData{}
//
//	sign.Uid = uid
//	sign.Date = time.Now().Format("20060102")
//
//	if err := DB().Insert(sign); err != nil {
//		DBError(err)
//	}
//
//	return sign
//}
