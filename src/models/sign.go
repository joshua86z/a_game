package models

import (
	"fmt"
	"time"
)

// role_sign
type SignData struct {
	Id     int    `db:"sign_id"`
	Uid    int64  `db:"uid"`
	Date   string `db:"sign_date"`
	Reward bool   `db:"sign_reward"`
}

func init() {
	DB().AddTableWithName(SignData{}, "role_sign").SetKeys(true, "Id")
}

type SignModel struct {
	Uid      int64
	SignList []*SignData
}

func GetSignModel(uid int64) *SignModel {

	var Sign SignModel

	var temp []*SignData
	_, err := DB().Select(&temp, "SELECT * FROM role_sign WHERE uid = ? ", uid)
	if err != nil {
		DBError(err)
	}

	Sign.Uid = uid
	Sign.SignList = temp

	return &Sign
}

func (this *SignModel) List() []*SignData {
	return this.SignList
}

func (this *SignModel) GetSign(signId int) *SignData {
	for _, sign := range this.SignList {
		if sign.Id == signId {
			return sign
		}
	}
	panic(fmt.Sprintf("没有这条数据 %d", signId))
}

func InsertSign(uid int64) *SignData {

	sign := &SignData{}

	sign.Uid = uid
	sign.Date = time.Now().Format("20060102")

	if err := DB().Insert(sign); err != nil {
		DBError(err)
	}

	return sign
}
